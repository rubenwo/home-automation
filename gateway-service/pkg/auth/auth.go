package auth

import (
	"github.com/rubenwo/home-automation/gateway-service/pkg/auth/models"
	"net/http"
)

//Authenticator is an interface that describes the functionality for authenticating http functions
type Authenticator interface {
	//Login should return an authorization and refresh token to the client
	Login(w http.ResponseWriter, r *http.Request)
	//Logout should clear/invalidate the authorization and refresh tokens
	Logout(w http.ResponseWriter, r *http.Request)
	//Register should create a new user
	Register(w http.ResponseWriter, r *http.Request)
	//RefreshToken receives the refresh token from the client, validates it, removes it and returns a new authorization and refresh token
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

//TokenExchanger is an interface that describes how an oauth token should be exchanged to our own claims
type TokenExchanger interface {
	ExchangeToken(r *http.Request) (models.Claims, error)
}
