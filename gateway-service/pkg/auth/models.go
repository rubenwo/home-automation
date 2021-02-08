package auth

type DefaultLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
}
