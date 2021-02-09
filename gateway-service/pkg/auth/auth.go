package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type Claims struct {
	Username      string        `json:"username"`
	UserID        string        `json:"user_id"`
	Authorization Authorization `json:"authorization"`
	jwt.StandardClaims
}

type Authorization struct {
	Roles []string `json:"roles"`
}

type Authenticator interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type TokenExchanger interface {
	ExchangeToken(r *http.Request) (Claims, error)
}
