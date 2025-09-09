package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type RouteResponse struct{
	Message string `json:"message"`
	ID string `json:"id,omitempty"`
}


type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Id string `json:"id,omitempty"`
}



type UserResponse struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Token string `json:"token"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}


type Claims struct {
	Username string `json:"username"`
	Id string `json:"id"`
	jwt.RegisteredClaims
}