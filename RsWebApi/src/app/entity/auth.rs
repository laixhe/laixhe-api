//! 鉴权相关实体 (对应 Go 项目 app/entity/auth.go)

use serde::{Deserialize, Serialize};

use super::user::User;

/// 请求-注册
///
/// 字段缺省为空字符串 (对齐 Go: 未注册 validator, validate:"required" 不生效, 由业务校验兜底)
#[derive(Debug, Deserialize)]
pub struct AuthRegisterRequest {
    /// 昵称
    #[serde(default)]
    pub nickname: String,
    /// 邮箱
    #[serde(default)]
    pub email: String,
    /// 密码
    #[serde(default)]
    pub password: String,
}

/// 响应-注册
#[derive(Debug, Serialize)]
pub struct AuthRegisterResponse {
    /// jwt token
    pub token: String,
    /// 用户信息
    pub user: User,
}

/// 请求-登录
///
/// 字段缺省为空字符串 (对齐 Go: 未注册 validator, validate:"required" 不生效, 由业务校验兜底)
#[derive(Debug, Deserialize)]
pub struct AuthLoginRequest {
    /// 邮箱
    #[serde(default)]
    pub email: String,
    /// 密码
    #[serde(default)]
    pub password: String,
}

/// 响应-登录
#[derive(Debug, Serialize)]
pub struct AuthLoginResponse {
    /// jwt token
    pub token: String,
    /// 用户信息
    pub user: User,
}

/// 请求-刷新Jwt (Uid 由 JWT 提供，不参与反序列化)
#[derive(Debug)]
pub struct AuthRefreshRequest {
    /// 用户id
    pub uid: u32,
}

/// 响应-刷新Jwt
#[derive(Debug, Serialize)]
pub struct AuthRefreshResponse {
    /// jwt token
    pub token: String,
    /// 用户信息
    pub user: User,
}
