package models

import "github.com/golang-jwt/jwt/v5"

//используется для работы с таблицой сессий пользователей в postgres
type Session struct {
	ID           string `json:"id"`
	IP           string `json:"ip"`
	RefreshToken string `json:"refreshtoken"`
	PairID       string `json:"pairid"`
	UserAgent    string `json:"useragent"`
}

//используется для генерации access токена
type MyCustomClaims struct {
	UserID string
	PairID string
	jwt.RegisteredClaims
}

//используется для отправки access токена в формате json пользователю
type AccessTokenJson struct {
	Access string `json:"access"`
}

//используется для отправки пользоватею его GUID из обработчика /auth/me
type UserJsonID struct {
	UserID string `json:"userid"`
}

//используется для отправки сообщения на webhook
type WebhookJson struct {
	Message string `json:"message"`
}
