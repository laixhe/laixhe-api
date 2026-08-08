//! 全局请求日志中间件 (X-Request-ID + 请求/响应耗时)

use axum::extract::Request;
use axum::http::{header::HeaderName, HeaderMap, HeaderValue};
use axum::middleware::Next;
use axum::response::Response;

use crate::log_elapsed;
use crate::logger::Timer;

/// X-Request-ID 请求/响应头名称
pub static X_REQUEST_ID: HeaderName = HeaderName::from_static("x-request-id");

/// 请求 ID，由全局请求日志中间件生成并放入 request extensions，
/// 供 JWT 中间件等提取用于日志关联。
#[derive(Debug, Clone)]
pub struct RequestId(pub String);

/// 从请求中提取 Request ID (若不存在则生成)
///
/// 注意: 透传的 X-Request-ID 可被客户端伪造, 仅用于日志关联无安全影响;
/// 生成时使用 xid 保证全局唯一 (Go fiber requestid 中间件使用 UUID, 算法不同但唯一性一致)。
fn resolve_request_id(headers: &HeaderMap) -> String {
    headers
        .get(X_REQUEST_ID.clone())
        .and_then(|v| v.to_str().ok())
        .filter(|s| !s.is_empty())
        .map(|s| s.to_string())
        .unwrap_or_else(|| xid::new().to_string())
}

/// 全局请求日志中间件：
/// 生成/透传 X-Request-ID，记录请求入口与响应完成的耗时，并回写响应头
pub async fn request_logger(req: Request, next: Next) -> Response {
    let start = Timer::new();
    let method = req.method().as_str().to_string();
    let path = req.uri().path().to_string();
    // 生成或透传 request id
    let request_id = resolve_request_id(req.headers());
    tracing::info!(
        request_id = %request_id,
        method = %method,
        path = %path,
        "收到请求"
    );

    let mut req = req;
    req.extensions_mut().insert(RequestId(request_id.clone()));

    let resp = next.run(req).await;
    let status = resp.status().as_u16();
    log_elapsed!(
        start,
        total_ms,
        info,
        request_id = %request_id,
        method = %method,
        path = %path,
        status,
        "请求处理完成"
    );

    // 回写 X-Request-ID 响应头
    let mut resp = resp;
    resp.headers_mut().insert(
        X_REQUEST_ID.clone(),
        HeaderValue::from_str(&request_id).unwrap_or_else(|_| HeaderValue::from_static("")),
    );
    resp
}
