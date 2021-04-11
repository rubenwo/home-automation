package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	BearerSchema = "Bearer "
	BasicSchema  = "Basic "
)

type DefaultClient struct {
	key             []byte
	adminEnabled    bool
	tokenExpiration time.Duration
}

func NewDefaultClient(key []byte, tokenExpiration time.Duration, adminEnabled bool) *DefaultClient {
	return &DefaultClient{key: key, tokenExpiration: tokenExpiration, adminEnabled: adminEnabled}
}

func (c *DefaultClient) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthorizationMiddleware => checking permissions...")
		for _, header := range r.Header["Upgrade"] {
			if header == "websocket" {
				fmt.Println("websocket connection")
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
			claims := &Claims{}

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
	var lr DefaultLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		http.Error(w, "invalid request body", http.StatusUnprocessableEntity)
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
	}

	expirationTime := time.Now().Add(c.tokenExpiration)

	claims := Claims{
		Username:      "admin",
		UserID:        "1",
		Authorization: Authorization{Roles: []string{"admin"}},
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

	response := LoginResponse{
		Username: claims.Username,
		UserID:   claims.UserID,
		Token:    tokenString,
	}

	w.Header().Set("Authorization", "bearer "+tokenString)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *DefaultClient) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *DefaultClient) Register(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Registering a new user is disabled by the admins", http.StatusForbidden)
}

func (c *DefaultClient) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
