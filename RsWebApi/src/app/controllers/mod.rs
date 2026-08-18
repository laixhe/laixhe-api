//! 业务控制器 (对应 Go 项目 app/controllers)
//!
//! 控制器职责：参数校验 → 调用 Service 业务逻辑 → 直接返回业务实体
//! (`Result<Json<T>, ApiError>`, 成功为裸实体 JSON, 与 Go 版 `ctx.JSON(resp)` 一致)。
//!
//! 接口文档: 通过 utoipa 注解 (`#[utoipa::path]`) 自动生成,
//! 运行 `make docs` 或 `webapi --gen-docs` 输出 `docs/swagger.json|yaml`,
//! 服务运行期由 `docs::ApiDoc::openapi()` 序列化提供 (`/api/v1/swagger.json|yaml`)。

pub mod auth;
pub mod health;
pub mod user;

use axum::extract::{FromRequest, FromRequestParts, Query, Request};
use axum::http::request::Parts;
use axum::http::StatusCode;
use axum::Json;
use serde::de::DeserializeOwned;

use crate::error::ApiError;

/// JSON 请求体提取器，解码失败返回统一 JSON 错误 (400)
///
/// 相比 `axum::extract::Json`，失败时不会返回默认的纯文本错误，
/// 而是转换为 [`ApiError::bad_request`]，保证响应格式统一。
///
/// # 示例
///
/// ```ignore
/// pub async fn login(
///     State(state): State<AppState>,
///     JsonBody(req): JsonBody<AuthLoginRequest>,
/// ) -> Result<Json<AuthLoginResponse>, ApiError> {
///     // ...
/// }
/// ```
pub struct JsonBody<T>(pub T);

impl<S, T> FromRequest<S> for JsonBody<T>
where
    S: Send + Sync,
    T: DeserializeOwned,
{
    type Rejection = ApiError;

    async fn from_request(req: Request, state: &S) -> Result<Self, Self::Rejection> {
        let Json(value) = Json::<T>::from_request(req, state)
            .await
            .map_err(|e| {
                // 请求体超限(4MB): 保留 413 状态码, 消息对齐 Go fiber 内建 BodyLimit 输出
                if e.status() == StatusCode::PAYLOAD_TOO_LARGE {
                    ApiError {
                        code: 413,
                        message: "Request Entity Too Large".to_string(),
                    }
                } else {
                    ApiError::bad_request(format!("Bad request: {}", e.body_text()))
                }
            })?;
        Ok(JsonBody(value))
    }
}

/// Query 参数提取器，解析失败返回统一 JSON 错误 (400)
///
/// 相比 `axum::extract::Query`，解析失败（如类型不匹配）时同样返回
/// [`ApiError::bad_request`]，保持统一响应格式。
///
/// # 示例
///
/// ```ignore
/// pub async fn info(
///     State(state): State<AppState>,
///     QueryParams(req): QueryParams<UserInfoRequest>,
/// ) -> Result<Json<User>, ApiError> {
///     // ...
/// }
/// ```
pub struct QueryParams<T>(pub T);

impl<S, T> FromRequestParts<S> for QueryParams<T>
where
    S: Send + Sync,
    T: DeserializeOwned,
{
    type Rejection = ApiError;

    async fn from_request_parts(parts: &mut Parts, state: &S) -> Result<Self, Self::Rejection> {
        let Query(value) = axum::extract::Query::<T>::from_request_parts(parts, state)
            .await
            .map_err(|e| ApiError::bad_request(format!("Bad request: {}", e.body_text())))?;
        Ok(QueryParams(value))
    }
}
