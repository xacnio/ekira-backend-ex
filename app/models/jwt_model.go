package models

type JWTClaims struct {
	Email  string `json:"Email"`
	UserId string `json:"UserId"`
}
