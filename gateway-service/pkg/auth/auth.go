package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

type Authenticator interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type FederatedAuthenticator interface {
	ExchangeToken(r *http.Request) (Claims, error)
}