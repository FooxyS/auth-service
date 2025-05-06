package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/FooxyS/auth-service/handlers/auth"
	"github.com/golang-jwt/jwt/v5"
)

func ParseTokenFromHeader(s string) string {
	substr := strings.Split(s, " ")
	return substr[1]
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//получение токена, user-agent, IP
		authBearer := r.Header.Get("Authorization")
		userAgent := r.Header.Get("User-Agent")

		//парсинг токена
		authToken := ParseTokenFromHeader(authBearer)

		//проверка токена на валидность и возврат userID
		accessClaims := auth.MyCustomClaims{}

		jwtkey, errGetKey := auth.GetJwtFromEnv()
		if errGetKey != nil {
			log.Printf("error with env: %v", errGetKey)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		jwt.ParseWithClaims(authToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtkey), nil
		})

		fmt.Printf("сравнивание user-agent: %s", userAgent)
	})
}
