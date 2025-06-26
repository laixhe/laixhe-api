package entity

// AuthRegisterRequest 请求-注册
type AuthRegisterRequest struct {
	Nickname string `json:"nickname"` // 昵称
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// AuthRegisterResponse 响应-注册
type AuthRegisterResponse struct {
	Token string `json:"token"` // jwt token
	User  *User  `json:"user"`  // 用户信息
}

// AuthLoginRequest 请求-登录
type AuthLoginRequest struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// AuthLoginResponse 响应-登录
type AuthLoginResponse struct {
	Token string `json:"token"` // jwt token
	User  *User  `json:"user"`  // 用户信息
}

// AuthRefreshRequest 请求-刷新Jwt
type AuthRefreshRequest struct {
	Uid int `json:"uid"` // 用户id
}

// AuthRefreshResponse 响应-刷新Jwt
type AuthRefreshResponse struct {
	Token string `json:"token"` // jwt token
	User  *User  `json:"user"`  // 用户信息
}
