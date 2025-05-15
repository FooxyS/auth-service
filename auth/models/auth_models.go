package models

import "github.com/golang-jwt/jwt/v5"

type Session struct {
	ID           string `json:"id"`
	IP           string `json:"ip"`
	RefreshToken string `json:"refreshtoken"`
	PairID       string `json:"pairid"`
	UserAgent    string `json:"useragent"`
}

type MyCustomClaims struct {
	UserID string
	PairID string
	jwt.RegisteredClaims
}

type AccessTokenJson struct {
	Access string `json:"access"`
}
