//! 路由 (对应 Go 项目 routers)
//!
//! 中间件执行顺序 (从外到内):
//! 请求日志/Request-ID → panic 恢复 → CORS → 响应 gzip 压缩 → 请求超时(408) → 请求体大小限制 → IP 限流(429) → 业务路由
//!
//! 注意: axum Router::layer 与 Fiber Use 相反, 是"后注册的在外层",
//! 因此下面 layer 需按上述顺序的逆序注册 (限流最先, 日志最后)。

pub mod auth;
pub mod user;

use std::any::Any;
use std::time::Duration;

use axum::body::Body;
use axum::extract::{DefaultBodyLimit, Request, State};
use axum::http::{header, StatusCode};
use axum::middleware::{from_fn, from_fn_with_state, Next};
use axum::response::{IntoResponse, Json, Response};
use axum::routing::get;
use axum::Router;
use tower_http::catch_panic::CatchPanicLayer;
use tower_http::compression::CompressionLayer;
use tower_http::cors::CorsLayer;
use utoipa::OpenApi;

use crate::app::controllers::health;
use crate::docs::ApiDoc;
use crate::error::ApiError;
use crate::middleware::rate_limit::rate_limit_middleware;
use crate::middleware::request_log::{request_logger, X_REQUEST_ID};
use crate::state::AppState;

/// 由 utoipa 注解生成的 OpenAPI 文档 (首次访问时生成并缓存, 与 Go 端 `swag init` 产物等价)
static SWAGGER_JSON: std::sync::LazyLock<String> = std::sync::LazyLock::new(|| {
    ApiDoc::openapi().to_json().expect("openapi to_json failed")
});
static SWAGGER_YAML: std::sync::LazyLock<String> = std::sync::LazyLock::new(|| {
    ApiDoc::openapi().to_yaml().expect("openapi to_yaml failed")
});

/// Swagger UI 页面 (HTML, 引用 jsdelivr CDN 的 swagger-ui 资源)
const SWAGGER_UI_HTML: &str = r#"<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>API 接口文档</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
  <style>html{box-sizing:border-box;overflow-y:scroll}body{margin:0;background:#fafafa}</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/api/v1/swagger.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset]
      });
    };
  </script>
</body>
</html>"#;

/// 构建总路由 (对应 Router.init)
pub fn build(state: AppState) -> Router {
    let api_v1 = Router::new()
        .route("/swagger.json", get(swagger_json))
        .route("/swagger.yaml", get(swagger_yaml))
        // Swagger UI 可视化页面 (浏览器访问 /api/v1/swagger)
        .route("/swagger", get(swagger_ui))
        // 健康检查 (含数据库探测)
        .route("/health", get(health::health))
        .merge(auth::router(state.clone()))
        .merge(user::router(state.clone()));
    Router::new()
        .nest("/api/v1", api_v1)
        // 全局 404 兜底, 返回统一 JSON 格式
        .fallback(not_found)
        // 405 方法不允许兜底: 返回统一 JSON (对齐 Go fiber 的 405 JSON 输出)
        .method_not_allowed_fallback(method_not_allowed)
        // 全局限流中间件 (基于客户端 IP, 超过阈值返回 429 统一 JSON)
        // 最内层: 需最先注册 (axum layer 后注册的在外层), 使其最贴近业务路由
        .layer(from_fn_with_state(state.clone(), rate_limit_middleware))
        // 请求体大小限制 4MB (Go 侧无显式中间件, 依赖 fiber 内建默认 BodyLimit=4MB, 行为一致)
        .layer(DefaultBodyLimit::max(4 * 1024 * 1024))
        // 请求超时中间件 (超过 http.timeout 秒未完成返回 408 统一 JSON, 对齐 Go fiber timeout.OnTimeout)
        .layer(from_fn_with_state(state.clone(), timeout_middleware))
        // 响应 gzip 压缩 (客户端携带 Accept-Encoding: gzip 时启用)
        .layer(CompressionLayer::new())
        // 跨域 CORS (位于限流外层, 保证 429 等错误响应也带 CORS 头)
        .layer(cors_layer())
        // panic 恢复 (对应 Fiber UseRecover; handler panic 时返回统一 JSON 500)
        .layer(CatchPanicLayer::custom(handle_panic))
        // 全局请求日志中间件 (X-Request-ID + 请求/响应耗时, 最外层)
        // 最后注册使其包裹所有中间件, 记录被限流/超时拦截的请求
        .layer(from_fn(request_logger))
        .with_state(state)
}

/// 请求超时中间件: 超过 http.timeout 秒未完成返回 408 统一 JSON (对应 Go fiber timeout.OnTimeout)
async fn timeout_middleware(
    State(state): State<AppState>,
    req: Request,
    next: Next,
) -> Response {
    // 请求超时时间: 优先读配置 http.timeout (秒), 缺省或配置 <= 0 时取 30 秒 (对齐 Go config.Check)
    let timeout = Duration::from_secs(state.config.http.timeout.filter(|&t| t > 0).unwrap_or(30) as u64);
    match tokio::time::timeout(timeout, next.run(req)).await {
        Ok(resp) => resp,
        // 超时: 返回统一 JSON 408 (对齐 Go fiber timeout.OnTimeout)
        Err(_) => ApiError {
            code: StatusCode::REQUEST_TIMEOUT.as_u16(),
            message: "Request Timeout".to_string(),
        }
        .into_response(),
    }
}

/// panic 恢复处理: 将处理器 panic 转为统一 JSON 500 (对应 Fiber 的 UseRecover)
fn handle_panic(_err: Box<dyn Any + Send + 'static>) -> Response {
    Response::builder()
        .status(StatusCode::INTERNAL_SERVER_ERROR)
        .header(header::CONTENT_TYPE, "application/json")
        .body(Body::from(
            r#"{"code":500,"message":"internal server error"}"#,
        ))
        .unwrap()
}

/// CORS 配置: 允许任意来源/方法/头 (适用于前后端分离场景, 可按需收紧)
fn cors_layer() -> CorsLayer {
    CorsLayer::new()
        .allow_origin(tower_http::cors::Any)
        .allow_methods(tower_http::cors::Any)
        .allow_headers([
            header::CONTENT_TYPE,
            header::AUTHORIZATION,
            header::ACCEPT,
            X_REQUEST_ID.clone(),
        ])
        .expose_headers([X_REQUEST_ID.clone()])
}

/// GET /api/v1/swagger.json
async fn swagger_json() -> impl IntoResponse {
    (
        [
            (header::CONTENT_TYPE, "application/json"),
            // 静态文档, 允许浏览器/代理短时缓存, 减少重复传输
            (header::CACHE_CONTROL, "public, max-age=300"),
        ],
        SWAGGER_JSON.as_str(),
    )
}

/// GET /api/v1/swagger.yaml
async fn swagger_yaml() -> impl IntoResponse {
    (
        [
            (header::CONTENT_TYPE, "application/x-yaml"),
            (header::CACHE_CONTROL, "public, max-age=300"),
        ],
        SWAGGER_YAML.as_str(),
    )
}

/// GET /api/v1/swagger (Swagger UI 可视化页面)
async fn swagger_ui() -> impl IntoResponse {
    (
        [
            (header::CONTENT_TYPE, "text/html; charset=utf-8"),
            (header::CACHE_CONTROL, "public, max-age=300"),
        ],
        SWAGGER_UI_HTML,
    )
}

/// 404 兜底, 与 ApiError 格式一致: {"code":404,"message":"Not Found"}
async fn not_found() -> impl IntoResponse {
    (
        StatusCode::NOT_FOUND,
        Json(serde_json::json!({ "code": 404, "message": "Not Found" })),
    )
}

/// 405 兜底, 与统一错误格式一致: {"code":405,"message":"Method Not Allowed"}
async fn method_not_allowed() -> impl IntoResponse {
    (
        StatusCode::METHOD_NOT_ALLOWED,
        Json(serde_json::json!({ "code": 405, "message": "Method Not Allowed" })),
    )
}
