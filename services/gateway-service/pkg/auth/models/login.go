package models
type DefaultLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username           string `json:"username"`
	UserID             string `json:"user_id"`
	AuthorizationToken string `json:"authorization_token"`
	RefreshToken       string `json:"refresh_token"`
}

type LogoutRequest struct {
	AuthorizationToken string `json:"authorization_token"`
	RefreshToken       string `json:"refresh_token"`
}
