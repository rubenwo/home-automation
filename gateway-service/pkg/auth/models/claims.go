package models

import "github.com/dgrijalva/jwt-go"

//Claims contains the values for the claims part of the JWT
type Claims struct {
	Username      string        `json:"username"`
	UserID        string        `json:"user_id"`
	Authorization Authorization `json:"authorization"`
	jwt.StandardClaims
}

//Authorization contains all the roles of the user
type Authorization struct {
	//Roles which is used for RBAc
	Roles []string `json:"roles"`
}
