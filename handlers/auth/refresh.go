package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	/*
		проверить, храниться ли refresh токен в БД, чтобы убедиться, что у нас не была выполнена операция logout
	*/
	pairIDFromDB := "заглушка. pairid из бд"
	userAgentFromDB := "заглушка. User-Agent из бд"

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
	jwtkey, errGotEnv := GetJwtFromEnv()
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
		проверить IP пользователя. Если запрос с нового IP, то отправить на webhook сообщение о попытке входа с нового устройства.
		(операция должна проболжиться)
	*/
}
