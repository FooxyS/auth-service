package auth

import (
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/auth/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//достаём пул подключений из контекста
	pgpool, ok := r.Context().Value("postgres").(*pgxpool.Pool)
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
	jwtkey, errGotEnv := services.GetFromEnv("JWT_KEY")
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

	_, errWithExec := pgpool.Exec(r.Context(), "DELETE FROM session_table WHERE user_id=$1", accessClaims.UserID)
	if errWithExec != nil {
		log.Printf("errWithExec: %v\n", errWithExec)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сессия пользователя успешно удалена"))
}
