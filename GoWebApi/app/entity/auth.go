package entity

// 提示: 各结构体中的 validate tag 仅用于 swag 生成 API 文档,
// 项目未注册 validator, tag 不参与请求校验; 真实校验在控制器层手写完成。

// AuthRegisterRequest 请求-注册
type AuthRegisterRequest struct {
	Nickname string `json:"nickname" validate:"required"` // 昵称
	Email    string `json:"email" validate:"required"`    // 邮箱
	Password string `json:"password" validate:"required"` // 密码
}

// AuthRegisterResponse 响应-注册
type AuthRegisterResponse struct {
	Token string `json:"token" validate:"required"` // jwt token
	User  *User  `json:"user" validate:"required"`  // 用户信息
}

// AuthLoginRequest 请求-登录
type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required"`    // 邮箱
	Password string `json:"password" validate:"required"` // 密码
}

// AuthLoginResponse 响应-登录
type AuthLoginResponse struct {
	Token string `json:"token" validate:"required"` // jwt token
	User  *User  `json:"user" validate:"required"`  // 用户信息
}

// AuthRefreshRequest 请求-刷新Jwt
type AuthRefreshRequest struct {
	Uid int `json:"uid"` // 用户id
}

// AuthRefreshResponse 响应-刷新Jwt
type AuthRefreshResponse struct {
	Token string `json:"token" validate:"required"` // jwt token
	User  *User  `json:"user" validate:"required"`  // 用户信息
}
