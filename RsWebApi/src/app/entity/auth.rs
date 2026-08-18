//! 鉴权相关实体 (对应 Go 项目 app/entity/auth.go)

use serde::{Deserialize, Serialize};
use utoipa::ToSchema;

use super::user::User;

/// 请求-注册
///
/// 字段缺省为空字符串 (对齐 Go: 未注册 validator, validate:"required" 不生效, 由业务校验兜底)
#[derive(Debug, Deserialize, ToSchema)]
pub struct AuthRegisterRequest {
    /// 昵称
    #[serde(default, deserialize_with = "super::de_null_to_empty")]
    #[schema(required)]
    pub nickname: String,
    /// 邮箱
    #[serde(default, deserialize_with = "super::de_null_to_empty")]
    #[schema(required)]
    pub email: String,
    /// 密码
    #[serde(default, deserialize_with = "super::de_null_to_empty")]
    #[schema(required)]
    pub password: String,
}

// 说明: 下面 4 个响应类型字段完全相同 ({token, user}), 但刻意分开定义:
// - AuthRegisterResponse / AuthLoginResponse / AuthRefreshResponse: 分别对应三个接口
//   handler 的实际返回类型, 与 Go 端 entity/auth.go 的同名类型一一对应, 便于各自独立演进
//   (如注册未来多返回一个字段, 不影响登录/刷新);
// - AuthTokenResponse: 仅用于 OpenAPI 文档 (docs.rs 中三个接口的响应体 schema 统一引用它),
//   对齐 PHP/TS 端共用的 AuthTokenResponse, 使"文档 schema"与"接口实现"解耦。

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
#[derive(Debug, Deserialize, ToSchema)]
pub struct AuthLoginRequest {
    /// 邮箱
    #[serde(default, deserialize_with = "super::de_null_to_empty")]
    #[schema(required)]
    pub email: String,
    /// 密码
    #[serde(default, deserialize_with = "super::de_null_to_empty")]
    #[schema(required)]
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
    pub uid: i32,
}

/// 响应-刷新Jwt
#[derive(Debug, Serialize)]
pub struct AuthRefreshResponse {
    /// jwt token
    pub token: String,
    /// 用户信息
    pub user: User,
}

/// 响应-鉴权成功 (注册/登录/刷新共用, 与 PHP/TS 端 AuthTokenResponse 对齐)
#[derive(Debug, Serialize, ToSchema)]
pub struct AuthTokenResponse {
    /// jwt token
    pub token: String,
    /// 用户信息
    pub user: User,
}

#[cfg(test)]
mod tests {
    use super::*;

    /// null 字段反序列化为空字符串 (与 Go/PHP 端 null→空值语义一致, 走业务校验 422)
    #[test]
    fn register_null_field_becomes_empty() {
        let req: AuthRegisterRequest =
            serde_json::from_str(r#"{"nickname":null,"email":"a@b.com","password":"pass123"}"#)
                .unwrap();
        assert_eq!(req.nickname, "");
    }

    /// 字段缺失仍默认空字符串 (serde(default) 行为不受 deserialize_with 影响)
    #[test]
    fn login_missing_field_defaults_empty() {
        let req: AuthLoginRequest = serde_json::from_str(r#"{"email":"a@b.com"}"#).unwrap();
        assert_eq!(req.password, "");
    }
}
