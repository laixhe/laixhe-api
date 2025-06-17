package entity

// AuthLoginRequest 请求-登录
type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthLoginResponse 响应-登录
type AuthLoginResponse struct {
	Token string `json:"token"`
}
