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
    // 仅 debug 级才为日志字段分配字符串: 默认 info 级下避免每请求 2 次无谓堆分配。
    // 不能直接传 &str 引用: req 在 next.run 时被 move, 借用无法跨 await 存活;
    // info 级下这些 debug! 不输出, 空字符串无任何影响。
    let method = if tracing::level_enabled!(tracing::Level::DEBUG) {
        req.method().as_str().to_string()
    } else {
        String::new()
    };
    let path = if tracing::level_enabled!(tracing::Level::DEBUG) {
        req.uri().path().to_string()
    } else {
        String::new()
    };
    // 生成或透传 request id
    let request_id = resolve_request_id(req.headers());
    // debug 级: 逐请求日志, 高 QPS 下 info 级别会产生大量日志开销 (见 jwt.rs 同策略)
    tracing::debug!(
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
        debug,
        request_id = %request_id,
        method = %method,
        path = %path,
        status,
        "请求处理完成"
    );

    // 回写 X-Request-ID 响应头。
    // request_id 要么是已通过 to_str() 校验的客户端头, 要么是 xid::new() 生成的纯 ASCII
    // 字母数字串 (见 resolve_request_id), 二者均为合法 HeaderValue, from_str 不可能失败,
    // 故直接 expect 而非保留不可能触发的 unwrap_or_else 兜底分支
    let mut resp = resp;
    resp.headers_mut().insert(
        X_REQUEST_ID.clone(),
        HeaderValue::from_str(&request_id)
            .expect("request_id 必为合法 ASCII HeaderValue (见 resolve_request_id)"),
    );
    resp
}
