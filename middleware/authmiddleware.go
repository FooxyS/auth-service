package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	auth "github.com/FooxyS/auth-service/auth/handlers"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
)

func ParseTokenFromHeader(s string) (string, error) {
	substr := strings.Split(s, " ")
	if len(substr) < 2 {
		return "", errors.New("massive is too short. Out of range")
	}

	return substr[1], nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		//проверка токена на валидность и возврат userID
		accessClaims := new(auth.MyCustomClaims)

		//достаю secretString из env
		jwtkey, errGetKey := auth.GetFromEnv("JWT_KEY")
		if errGetKey != nil {
			log.Printf("error with env: %v", errGetKey)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		//парсинг JWT в структуру и проверка валидности токена
		token, errWithParseToken := jwt.ParseWithClaims(authToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtkey), nil
		})
		if errWithParseToken != nil || !token.Valid {
			log.Printf("error with parsing JWT: %v", errWithParseToken)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//получение id из claims, user-agent, IP
		id := accessClaims.UserID
		userAgent := r.Header.Get("User-Agent")

		/*
			логика получения IP пользователя
		*/

		/*
			достаём нужные данные из БД по userID
		*/
		fmt.Println("сравниваем ", userAgent)
		fmt.Println("сравниваем IP")

		/*
			если что-то не сошлось отправляем ошибку авторизации

			удаляем refresh токен из БД
		*/

		//кладём в контекст userID, отправляем в обработчик
		ctx := context.WithValue(r.Context(), consts.USER_ID_KEY, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
