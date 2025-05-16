package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/auth/services"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {

	//достаём пул подключений из контекста
	pgpool, ok := r.Context().Value(consts.CTX_KEY_DB).(*pgxpool.Pool)
	if !ok {
		log.Println("value not found in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//парсинг access токена, сравнение pairid
	//проверка хэдера авторизации
	authBearer := r.Header.Get("Authorization")
	if authBearer == "" {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return
	}
	authToken, errMassShort := services.ParseTokenFromHeader(authBearer)
	if errMassShort != nil {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return
	}

	//загрузка секретной строки
	jwtkey, errGotEnv := services.GetFromEnv(consts.JWT_KEY)
	if errGotEnv != nil {
		log.Printf("error with env: %v", errGotEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	accessClaims := new(models.MyCustomClaims)

	token, errWithParseToken := jwt.ParseWithClaims(authToken, accessClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})
	if errWithParseToken != nil || !token.Valid {
		log.Printf("error with parsing JWT: %v", errWithParseToken)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//проверяем данные о сессии пользователя по userid
	session := new(models.Session)
	errQueryRow := pgpool.QueryRow(r.Context(), "SELECT * FROM session_table WHERE user_id=$1", accessClaims.UserID).Scan(&session.ID, &session.IP, &session.PairID, &session.RefreshToken, &session.UserAgent)
	if errQueryRow != nil {
		log.Printf("error with QueryRow: %v\n", errQueryRow)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//записываем информацию о куке в структуру
	refreshCookie, errNoCookie := r.Cookie("refreshtoken")
	if errNoCookie != nil {
		log.Printf("errNoCookie: %v\n", errNoCookie)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//проверка pairid. Если false, то выдача новой пары запрещена -> деавторизация
	if accessClaims.PairID != session.PairID {
		_, errWithExec := pgpool.Exec(r.Context(), "DELETE FROM session_table WHERE user_id=$1", accessClaims.UserID)
		if errWithExec != nil {
			log.Printf("errWithExec: %v\n", errWithExec)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("деавторизация из-за несовпадения pairID")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//проверяем expires refresh токена. Если истёк, то удаляем инфу о сессиии пользователя
	if isbefore := refreshCookie.Expires.Compare(time.Now()); isbefore == 1 {
		_, errWithExec := pgpool.Exec(r.Context(), "DELETE FROM session_table WHERE user_id=$1", accessClaims.UserID)
		if errWithExec != nil {
			log.Printf("errWithExec: %v\n", errWithExec)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("деавторизация из-за истечения срока годности у токена refresh")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//получение ip пользователя из запроса
	ipFromReq, _, errWithSplitIP := net.SplitHostPort(r.RemoteAddr)
	if errWithSplitIP != nil {
		log.Printf("error with parse remoteAddr to IP: %v\n", errWithSplitIP)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//сравнение IP адресов
	if ipFromReq != session.IP {
		//логика отправки post запроса на заданный webhook...
		errSendWebhook := services.SendWebhook(ipFromReq)
		if errSendWebhook != nil {
			log.Printf("error with sending webhook: %v\n", errSendWebhook)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	ErrWithCompareResresh := bcrypt.CompareHashAndPassword([]byte(session.RefreshToken), []byte(refreshCookie.Value))
	if ErrWithCompareResresh != nil {
		_, errWithExec := pgpool.Exec(r.Context(), "DELETE FROM session_table WHERE user_id=$1", accessClaims.UserID)
		if errWithExec != nil {
			log.Printf("errWithExec: %v\n", errWithExec)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("деавторизация из-за несовпадения refresh токенов")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	/*
		достать User-Agent из БД и запроса, сравнить. Если не совпадает, то деавторизовать пользователя.
		(удалить сессию из БД, отправить код unauthorized)
	*/
	//достаём User-Agent из запроса
	agentFromReq := r.Header.Get("User-Agent")
	if session.UserAgent != agentFromReq {
		_, errWithExec := pgpool.Exec(r.Context(), "DELETE FROM session_table WHERE user_id=$1", accessClaims.UserID)
		if errWithExec != nil {
			log.Printf("errWithExec: %v\n", errWithExec)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("деавторизация из-за несовпадения user-agent")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//создание нового uuid пары токенов
	newPairID := uuid.New().String()

	//создание нового refresh токена
	b := make([]byte, 32)
	_, errWithRand := rand.Read(b)
	if errWithRand != nil {
		log.Printf("error with filling slice: %v\n", errWithRand)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newRefreshToken := base64.URLEncoding.EncodeToString(b)

	newCookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    newRefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(w, &newCookie)

	//создание нового access токена
	newAccessClaims := models.MyCustomClaims{
		UserID: accessClaims.UserID,
		PairID: newPairID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	newAccessToken, errWithGenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, newAccessClaims).SignedString([]byte(jwtkey))
	if errWithGenAccess != nil {
		log.Printf("error with generating access token: %v\n", errWithGenAccess)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	AccessJson := models.AccessTokenJson{
		Access: newAccessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	errParseJson := json.NewEncoder(w).Encode(AccessJson)
	if errParseJson != nil {
		log.Printf("error with parsing json response: %v\n", errParseJson)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//хэшируем refresh токен для записи в БД
	hashedNewRefresh, errHashRefresh := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if errHashRefresh != nil {
		log.Printf("error with hashing refresh token: %v\n", errHashRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//обновляем данные о сессии
	_, errWithExec := pgpool.Exec(r.Context(), "UPDATE session_table SET ip_address=$1, refresh_token=$2, pair_id=$3 WHERE user_id=$4", ipFromReq, string(hashedNewRefresh), newPairID, accessClaims.UserID)
	if errWithExec != nil {
		log.Printf("ErrWithExec: %v\n", errWithExec)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
