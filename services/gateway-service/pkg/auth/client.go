package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rubenwo/home-automation/services/gateway-service/pkg/auth/models"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	BearerSchema = "Bearer "
	BasicSchema  = "Basic "
)

type DefaultClient struct {
	key                          []byte
	adminEnabled                 bool
	authorizationTokenExpiration time.Duration
	refreshTokenExpiration       time.Duration

	refreshTokenLock sync.Mutex
	refreshTokens    map[string]time.Time

	isAllowedAnonymous map[string]bool
	allowedRolesRoutes map[string][]string
}

//NewDefaultClient
func NewDefaultClient(key []byte, authorizationTokenExpiration time.Duration, refreshTokenExpiration time.Duration, adminEnabled bool, isAllowedAnonymous map[string]bool) *DefaultClient {
	return &DefaultClient{
		key:                          key,
		authorizationTokenExpiration: authorizationTokenExpiration,
		refreshTokenExpiration:       refreshTokenExpiration,
		adminEnabled:                 adminEnabled,

		refreshTokens:      make(map[string]time.Time),
		isAllowedAnonymous: isAllowedAnonymous,
		allowedRolesRoutes: make(map[string][]string),
	}
}

func (c *DefaultClient) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthorizationMiddleware => checking permissions...")

		// Skip authorization if websocket
		// TODO: this should be secured
		for _, header := range r.Header["Upgrade"] {
			if header == "websocket" {
				fmt.Println("websocket connection")
				next.ServeHTTP(w, r)
				return
			}
		}

		s := strings.Split(r.RequestURI, "/")

		for k, value := range c.isAllowedAnonymous {
			if !value {
				continue
			}

			ks := strings.Split(k, "/")
			if len(ks) != len(s) {
				continue
			}

			matches := true
			for i := 0; i < len(ks); i++ {
				if strings.Contains(ks[i], "{") {
					continue
				}
				if ks[i] != s[i] {
					matches = false
					break
				}
			}
			if matches && value {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		if authHeader[:len(BasicSchema)] == BasicSchema {
			log.Println("processing basic auth")
			username, password, authOk := r.BasicAuth()
			if !authOk {
				http.Error(w, "basic auth failed", http.StatusUnauthorized)
				return
			}

			if username == "admin" {
				if c.adminEnabled {
					adminPassword := os.Getenv("ADMIN_PWD")
					if password != adminPassword {
						http.Error(w, "username/password is incorrect", http.StatusUnauthorized)
						return
					}
				} else {
					http.Error(w, "admin account is disabled", http.StatusUnauthorized)
					return
				}
			}
		}

		if authHeader[:len(BearerSchema)] == BearerSchema {
			log.Println("processing bearer token")

			tknStr := authHeader[len(BearerSchema):]

			// Initialize a new instance of `Claims`
			claims := &models.Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
				return c.key, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "invalid jwt signature", http.StatusUnauthorized)
					return
				}
				http.Error(w, fmt.Sprintf("error parsing token: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			if !tkn.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if claims.Username == "admin" && !c.adminEnabled {
				http.Error(w, "admin account is disabled", http.StatusUnauthorized)
				return
			}

		}

		log.Println("AuthorizationMiddleware => All OK!")

		next.ServeHTTP(w, r)
		return
	})
}

func (c *DefaultClient) Login(w http.ResponseWriter, r *http.Request) {
	var lr models.DefaultLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		http.Error(w, "invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if lr.Username == "" {
		http.Error(w, "username cannot be empty", http.StatusUnauthorized)
		return
	}
	if lr.Password == "" {
		http.Error(w, "password cannot be empty", http.StatusUnauthorized)
		return
	}

	if lr.Username == "admin" {
		if c.adminEnabled {
			adminPassword := os.Getenv("ADMIN_PWD")
			if lr.Password != adminPassword {
				http.Error(w, "username/password is incorrect", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "admin account is disabled", http.StatusUnauthorized)
			return
		}
	} else {
		log.Printf("'%s' is not a valid username", lr.Username)
		http.Error(w, fmt.Sprintf("'%s' does not exist", lr.Username), http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(c.authorizationTokenExpiration)

	// TODO: Get claims from the database

	claims := models.Claims{
		Username:      "admin",
		UserID:        "1",
		Authorization: models.Authorization{Roles: []string{"admin"}},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(c.key)
	if err != nil {
		http.Error(w, fmt.Sprintf("error signing token"), http.StatusInternalServerError)
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: time.Now().Add(c.refreshTokenExpiration).Unix()})
	refreshTokenString, err := refreshToken.SignedString(c.key)
	if err != nil {
		http.Error(w, fmt.Sprintf("error signing refresh token"), http.StatusInternalServerError)
		return
	}

	c.refreshTokenLock.Lock()
	c.refreshTokens[refreshTokenString] = time.Now().Add(c.refreshTokenExpiration)
	c.refreshTokenLock.Unlock()

	response := models.LoginResponse{
		Username:           claims.Username,
		UserID:             claims.UserID,
		AuthorizationToken: tokenString,
		RefreshToken:       refreshTokenString,
	}

	w.Header().Set("Authorization", "bearer "+tokenString)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *DefaultClient) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: remove authorization and refresh tokens from the database
	var lr models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		http.Error(w, "invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if lr.AuthorizationToken == "" {
		http.Error(w, "username cannot be empty", http.StatusUnauthorized)
		return
	}
	if lr.RefreshToken == "" {
		http.Error(w, "password cannot be empty", http.StatusUnauthorized)
		return
	}

	c.refreshTokenLock.Lock()
	delete(c.refreshTokens, lr.RefreshToken)
	c.refreshTokenLock.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (c *DefaultClient) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: add the option to create new accounts?
	http.Error(w, "Registering a new user is disabled by the admins", http.StatusForbidden)
}

func (c *DefaultClient) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh-token")
	log.Println(refreshToken)

	// Validate the refresh token

	tkn, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return c.key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "invalid jwt signature", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("error parsing token: %s", err.Error()), http.StatusUnauthorized)
		return
	}

	if !tkn.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	//TODO: Check with database, but for now we just use an in-memory map

	// Get the expiration time
	c.refreshTokenLock.Lock()
	expiration, exists := c.refreshTokens[refreshToken]
	c.refreshTokenLock.Unlock()
	if !exists {
		http.Error(w, "token does not exist", http.StatusUnauthorized)
		return
	}
	if time.Now().After(expiration) {
		http.Error(w, "token is expired", http.StatusUnauthorized)
		return
	}

	c.refreshTokenLock.Lock()
	delete(c.refreshTokens, refreshToken)
	c.refreshTokenLock.Unlock()

	// If token is valid, generate a new authorization and refresh token
	// Delete the token
	expirationTime := time.Now().Add(c.authorizationTokenExpiration)

	// TODO: Get claims from the database
	claims := models.Claims{
		Username:      "admin",
		UserID:        "1",
		Authorization: models.Authorization{Roles: []string{"admin"}},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(c.key)
	if err != nil {
		http.Error(w, fmt.Sprintf("error signing token"), http.StatusInternalServerError)
		return
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: time.Now().Add(c.refreshTokenExpiration).Unix()})
	newRefreshTokenString, err := newRefreshToken.SignedString(c.key)
	if err != nil {
		http.Error(w, fmt.Sprintf("error signing refresh token"), http.StatusInternalServerError)
		return
	}

	c.refreshTokenLock.Lock()
	c.refreshTokens[newRefreshTokenString] = time.Now().Add(c.refreshTokenExpiration)
	c.refreshTokenLock.Unlock()

	response := models.LoginResponse{
		Username:           claims.Username,
		UserID:             claims.UserID,
		AuthorizationToken: tokenString,
		RefreshToken:       newRefreshTokenString,
	}

	w.Header().Set("Authorization", "bearer "+tokenString)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
