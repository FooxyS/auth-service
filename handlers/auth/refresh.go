package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
Требования к операции `refresh`

1. Операцию `refresh` можно выполнить только той парой токенов, которая была выдана вместе.

2. Необходимо запретить операцию обновления токенов при изменении `User-Agent`. При этом,
после неудачной попытки выполнения операции, необходимо деавторизовать пользователя,
который попытался выполнить обновление токенов.

3. При попытке обновления токенов с нового IP необходимо отправить POST-запрос на заданный `webhook`
с информацией о попытке входа со стороннего IP. Запрещать операцию в данном случае не нужно.
*/

func ParseTokenFromHeader(s string) (string, error) {
	substr := strings.Split(s, " ")
	if len(substr) < 2 {
		return "", errors.New("massive is too short. Out of range")
	}

	return substr[1], nil
}

type WebhookJson struct {
	Message string `json:"message"`
}

func SendWebhook(ip string) error {
	//создание ответа
	text := fmt.Sprintf("Попытка входа с нового IP: %s\n", ip)

	respBody := WebhookJson{
		Message: text,
	}

	jsonRespBody, errParseJson := json.Marshal(respBody)
	if errParseJson != nil {
		return errParseJson
	}

	//достаём URL из env
	webhookurl, errGetEnv := GetFromEnv("WEBHOOK_URL")
	if errGetEnv != nil || webhookurl == "" {
		log.Printf("error with get from env: %v\n", errGetEnv)
		return errGetEnv
	}

	//формируем новый запрос
	req, errResp := http.NewRequest(http.MethodPost, webhookurl, bytes.NewBuffer([]byte(jsonRespBody)))
	if errResp != nil {
		log.Printf("error with creating new request: %v\n", errResp)
		return errResp
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, errDoReq := client.Do(req)
	if errDoReq != nil {
		log.Printf("error with sending request: %v\n", errDoReq)
		return errDoReq
	}

	defer resp.Body.Close()

	log.Printf("Webhook отправлен. Статус ответа: %v\n", resp.Status)
	return nil
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	/*
		проверить, храниться ли refresh токен в БД, чтобы убедиться, что у нас не была выполнена операция logout
	*/
	pairIDFromDB := "заглушка. pairid из бд"
	userAgentFromDB := "заглушка. User-Agent из бд"
	IPAddrFromDB := "заглушка. IP из бд"

	/*
		проверить IP пользователя. Если запрос с нового IP, то отправить на webhook сообщение о попытке входа с нового устройства.
		(операция должна проболжиться)
	*/
	//получение ip пользователя из запроса
	ipFromReq, _, errWithSplitIP := net.SplitHostPort(r.RemoteAddr)
	if errWithSplitIP != nil {
		log.Printf("error with parse remoteAddr to IP: %v\n", errWithSplitIP)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//сравнение IP адресов
	if ipFromReq != IPAddrFromDB {
		//логика отправки post запроса на заданный webhook...
		errSendWebhook := SendWebhook(ipFromReq)
	}

	/*
		достать User-Agent из БД и запроса, сравнить. Если не совпадает, то деавторизовать пользователя.
		(удалить сессию из БД, отправить код unauthorized)
	*/
	//достаём User-Agent из запроса
	agentFromReq := r.Header.Get("User-Agent")
	if userAgentFromDB != agentFromReq {
		//логика удаления сессии пользователя из БД...

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*
		распарсить access, проверить pairid -> были выданы вместе -> доступ разрешён
	*/

	//парсинг access токена, сравнение pairid
	//проверка хэдера авторизации
	authBearer := r.Header.Get("Authorization")
	if authBearer == "" {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return
	}
	authToken, errMassShort := ParseTokenFromHeader(authBearer)
	if errMassShort != nil {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return
	}

	//загрузка секретной строки
	jwtkey, errGotEnv := GetFromEnv("JWT_KEY")
	if errGotEnv != nil {
		log.Printf("error with env: %v", errGotEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	accessClaims := new(MyCustomClaims)

	token, errWithParseToken := jwt.ParseWithClaims(authToken, accessClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})
	if errWithParseToken != nil || !token.Valid {
		log.Printf("error with parsing JWT: %v", errWithParseToken)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//проверка pairid
	if accessClaims.PairID != pairIDFromDB {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*
		создание новой пары токенов
	*/
	//создание нового uuid пары токенов
	newPairID := uuid.New().String()

	//создание нового access токена
	newAccessClaims := MyCustomClaims{
		UserID: accessClaims.UserID,
		PairID: newPairID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	newAccessToken, errWithGenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, newAccessClaims).SignedString(jwtkey)
	if errWithGenAccess != nil {
		log.Printf("error with generating access token: %v\n", errWithGenAccess)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	AccessJson := AccessTokenJson{
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

	//создание нового refresh токена
	b := make([]byte, 32)
	newRefreshToken := base64.URLEncoding.EncodeToString(b)

	newCookie := http.Cookie{
		Name:     "refresh-token",
		Value:    newRefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &newCookie)

	//добавление refresh токена в БД
	hashedNewRefresh, errHashRefresh := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if errHashRefresh != nil {
		log.Printf("error with hashing refresh token: %v\n", errHashRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//логика добавления новых данных в БД...
	fmt.Println("загушка добавления ", hashedNewRefresh)

}
