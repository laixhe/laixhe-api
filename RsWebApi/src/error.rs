//! 错误类型 (对应 Go 项目 core/error.go 与 gonet/xfiber 的错误约定)
//!
//! 响应格式与 fiber.Error 一致: `{"code": <int>, "message": "<string>"}`
//!
//! # 错误码约定
//!
//! | code | HTTP 状态 | 含义 |
//! | ---- | --------- | ---- |
//! | 400  | 400       | 请求体 / Query 解析失败（对应 Go fiber bind 的 "Bad request: ..."） |
//! | 422  | 422       | 参数错误（业务校验失败，对应 Go `xfiber.ParamError`） |
//! | 401  | 401       | 未授权（缺少 / 无效 JWT，或用户被禁用） |
//! | 405  | 405       | 方法不允许（由路由层 method_not_allowed_fallback 返回） |
//! | 408  | 408       | 请求超时（超过 `http.timeout` 秒，由路由层超时中间件返回） |
//! | 404  | 404       | 路由不存在（由路由层 fallback 兜底） |
//! | 413  | 413       | 请求体超限（超过 4MB，由 JsonBody 提取器保留原始状态码） |
//! | 429  | 429       | 触发接口限流 |
//! | 500  | 500       | 服务器内部错误（数据库 / bcrypt / JWT 签发等，固定文案） |
//! | 503  | 503       | 服务不可用（健康检查数据库探测失败） |

use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::Json;
use sea_orm::DbErr;
use serde::Serialize;

/// API 错误
///
/// 控制器 / 中间件返回的统一错误类型，`?` 运算符可直接传播
/// （[`From<DbErr>`] 会把数据库错误转换为 500）。
#[derive(Debug, Clone, Serialize)]
pub struct ApiError {
    /// 错误码 (与 HTTP 状态码一致)
    pub code: u16,
    /// 错误描述
    pub message: String,
}

impl ApiError {
    /// 请求体 / Query 解析失败 400 (对应 Go fiber bind 的 "Bad request: ...")
    pub fn bad_request(message: impl Into<String>) -> Self {
        ApiError {
            code: StatusCode::BAD_REQUEST.as_u16(),
            message: message.into(),
        }
    }

    /// 参数错误 422 (对应 Go `xfiber.ParamError` = fiber.StatusUnprocessableEntity)
    ///
    /// # 示例
    ///
    /// ```ignore
    /// return Err(ApiError::param_error("邮箱格式错误"));
    /// ```
    pub fn param_error(message: impl Into<String>) -> Self {
        ApiError {
            code: StatusCode::UNPROCESSABLE_ENTITY.as_u16(),
            message: message.into(),
        }
    }

    /// 未授权 401 (对应 xfiber.AuthorizedError)
    pub fn unauthorized() -> Self {
        ApiError {
            code: StatusCode::UNAUTHORIZED.as_u16(),
            message: "Unauthorized".to_string(),
        }
    }

    /// 服务器内部错误 500
    pub fn internal(message: impl Into<String>) -> Self {
        ApiError {
            code: StatusCode::INTERNAL_SERVER_ERROR.as_u16(),
            message: message.into(),
        }
    }
}

impl IntoResponse for ApiError {
    fn into_response(self) -> Response {
        let body = Json(serde_json::json!({
            "code": self.code,
            "message": self.message,
        }));
        (
            StatusCode::from_u16(self.code).unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
            body,
        )
            .into_response()
    }
}

impl std::fmt::Display for ApiError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for ApiError {}

/// 数据库错误统一转换为 500 (收敛为固定文案, 避免 DB 错误细节泄露给客户端,
/// 对齐 Go fork 的 DefaultErrorHandler; 原始错误记录到服务端日志)
impl From<DbErr> for ApiError {
    fn from(err: DbErr) -> Self {
        tracing::error!(%err, "数据库操作失败");
        ApiError::internal("internal server error")
    }
}
