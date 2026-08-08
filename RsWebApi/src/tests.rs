//! 集成测试: 使用 axum oneshot 直接调用路由, 覆盖核心接口
//!
//! 依赖本机 MySQL (config.yaml 配置), 运行: cargo test

use axum::body::{to_bytes, Body};
use axum::http::{header, Request, StatusCode};
use axum::routing::get;
use axum::Router;
use sea_orm::{ColumnTrait, EntityTrait, QueryFilter};
use serde_json::Value;
use std::sync::Arc;
use std::time::Duration;
use tower::ServiceExt;
use tower_http::timeout::TimeoutLayer;

use crate::app::models::{user as user_model, user_extend, user_third_party};
use crate::middleware::rate_limit::RateLimiter;
use crate::routes;
use crate::state::AppState;

/// 构建测试应用 (返回路由与状态, 便于清理测试数据)
async fn build_app() -> (Router, AppState) {
    let state = AppState::new("./config.yaml").await;
    let router = routes::build(state.clone());
    (router, state)
}

/// 发送请求并返回 (status, body_json)
async fn send(app: &Router, req: Request<Body>) -> (StatusCode, Value) {
    let resp = app.clone().oneshot(req).await.expect("request failed");
    let status = resp.status();
    let body = to_bytes(resp.into_body(), 1024 * 1024)
        .await
        .expect("read response body failed");
    let json: Value = serde_json::from_slice(&body).unwrap_or(Value::Null);
    (status, json)
}

/// POST JSON 请求快捷构造
fn post_json(uri: &str, body: impl Into<String>) -> Request<Body> {
    Request::builder()
        .method("POST")
        .uri(uri)
        .header(header::CONTENT_TYPE, "application/json")
        .body(Body::from(body.into()))
        .unwrap()
}

/// 生成唯一测试邮箱
fn unique_email() -> String {
    format!(
        "test_{}@example.com",
        std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos()
    )
}

/// 注册一个测试用户, 返回 (uid, token)
async fn register_user(app: &Router, email: &str) -> (u32, String) {
    let req = post_json(
        "/api/v1/auth/register",
        format!(r#"{{"nickname":"testuser","email":"{email}","password":"abc123"}}"#),
    );
    let (status, json) = send(app, req).await;
    assert_eq!(status, StatusCode::OK, "register should succeed: {json}");
    // 成功响应为裸实体 JSON (与 Go 版一致), 无 code/message 包裹
    let uid = json["user"]["uid"].as_i64().unwrap() as u32;
    let token = json["token"].as_str().unwrap().to_string();
    (uid, token)
}

/// 清理测试用户数据 (user + user_extend + user_third_party)
async fn cleanup_user(state: &AppState, uid: u32) {
    let db = &state.db;
    user_third_party::Entity::delete_many()
        .filter(user_third_party::Column::Uid.eq(uid))
        .exec(db)
        .await
        .unwrap();
    user_extend::Entity::delete_many()
        .filter(user_extend::Column::Uid.eq(uid))
        .exec(db)
        .await
        .unwrap();
    user_model::Entity::delete_many()
        .filter(user_model::Column::Id.eq(uid))
        .exec(db)
        .await
        .unwrap();
}

#[tokio::test]
async fn swagger_json_ok() {
    let (app, _state) = build_app().await;
    let req = Request::builder()
        .uri("/api/v1/swagger.json")
        .body(Body::empty())
        .unwrap();
    let resp = app.clone().oneshot(req).await.unwrap();
    assert_eq!(resp.status(), StatusCode::OK);
    // 全局请求日志中间件应回写 X-Request-ID 响应头
    assert!(resp.headers().contains_key("X-Request-ID"));
}

#[tokio::test]
async fn swagger_ui_page_ok() {
    let (app, _state) = build_app().await;
    let req = Request::builder()
        .uri("/api/v1/swagger")
        .body(Body::empty())
        .unwrap();
    let resp = app.clone().oneshot(req).await.unwrap();
    assert_eq!(resp.status(), StatusCode::OK);
    // 应为 HTML 页面
    assert!(resp
        .headers()
        .get(header::CONTENT_TYPE)
        .and_then(|v| v.to_str().ok())
        .map(|ct| ct.contains("text/html"))
        .unwrap_or(false));
    let body = to_bytes(resp.into_body(), 1024 * 1024).await.unwrap();
    let html = String::from_utf8_lossy(&body);
    assert!(
        html.contains("swagger-ui"),
        "页面应包含 swagger-ui 资源引用"
    );
}

#[tokio::test]
async fn health_check_ok() {
    let (app, _state) = build_app().await;
    let req = Request::builder()
        .uri("/api/v1/health")
        .body(Body::empty())
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK);
    // 成功响应为裸实体 JSON, 无 code/message 包裹
    assert_eq!(json["status"], "ok");
    assert_eq!(json["database"], "up");
}

#[tokio::test]
async fn rate_limit_returns_429() {
    // 使用阈值=2 的限流器构造测试状态 (其余字段复用默认初始化)
    // 注意: /api/v1/health 已豁免限流, 此处用 swagger.json 路径验证
    let state = AppState::new("./config.yaml").await;
    let state = AppState {
        limiter: Arc::new(RateLimiter::new(2, Duration::from_secs(60))),
        ..state
    };
    let app = routes::build(state);
    let req = || {
        Request::builder()
            .uri("/api/v1/swagger.json")
            .body(Body::empty())
            .unwrap()
    };
    // 前 2 次允许通过
    let (status, _) = send(&app, req()).await;
    assert_eq!(status, StatusCode::OK);
    let (status, _) = send(&app, req()).await;
    assert_eq!(status, StatusCode::OK);
    // 第 3 次触发限流, 返回 429 统一 JSON
    let (status, json) = send(&app, req()).await;
    assert_eq!(status, StatusCode::TOO_MANY_REQUESTS);
    assert_eq!(json["code"], 429);
}

#[tokio::test]
async fn gzip_compression_applied() {
    // 携带 Accept-Encoding: gzip 时, 响应应带有 Content-Encoding: gzip
    let (app, _state) = build_app().await;
    let req = Request::builder()
        .uri("/api/v1/health")
        .header(header::ACCEPT_ENCODING, "gzip")
        .body(Body::empty())
        .unwrap();
    let resp = app.clone().oneshot(req).await.unwrap();
    assert_eq!(resp.status(), StatusCode::OK);
    assert_eq!(
        resp.headers()
            .get(header::CONTENT_ENCODING)
            .and_then(|v| v.to_str().ok()),
        Some("gzip")
    );
}

#[tokio::test]
async fn request_timeout_returns_408() {
    // 验证超时中间件: 超过时限的请求返回 408 (独立路由, 不依赖数据库)
    let app = axum::Router::new()
        .route(
            "/slow",
            get(|| async {
                tokio::time::sleep(Duration::from_millis(500)).await;
                "done"
            }),
        )
        .layer(TimeoutLayer::with_status_code(
            StatusCode::REQUEST_TIMEOUT,
            Duration::from_millis(50),
        ));
    let req = Request::builder().uri("/slow").body(Body::empty()).unwrap();
    let resp = app.oneshot(req).await.unwrap();
    assert_eq!(resp.status(), StatusCode::REQUEST_TIMEOUT);
}

#[tokio::test]
async fn panic_recover_returns_500_json() {
    // 验证 panic 恢复中间件 (对应 Fiber UseRecover): handler panic 返回统一 JSON 500
    use std::any::Any;
    use tower_http::catch_panic::CatchPanicLayer;

    fn handle_panic(_err: Box<dyn Any + Send + 'static>) -> axum::response::Response {
        axum::response::Response::builder()
            .status(StatusCode::INTERNAL_SERVER_ERROR)
            .body(axum::body::Body::from(
                r#"{"code":500,"message":"internal server error"}"#,
            ))
            .unwrap()
    }

    let app = axum::Router::new()
        .route(
            "/boom",
            get(|| async {
                panic!("boom");
                #[allow(unreachable_code)]
                ""
            }),
        )
        .layer(CatchPanicLayer::custom(handle_panic));
    let req = Request::builder().uri("/boom").body(Body::empty()).unwrap();
    let resp = app.oneshot(req).await.unwrap();
    assert_eq!(resp.status(), StatusCode::INTERNAL_SERVER_ERROR);
    let body = to_bytes(resp.into_body(), 1024 * 1024).await.unwrap();
    let json: Value = serde_json::from_slice(&body).unwrap();
    assert_eq!(json["code"], 500);
}

#[tokio::test]
async fn not_found_returns_unified_json() {
    let (app, _state) = build_app().await;
    let req = Request::builder()
        .uri("/api/v1/no-such-route")
        .body(Body::empty())
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::NOT_FOUND);
    assert_eq!(json["code"], 404);
    assert_eq!(json["message"], "Not Found");
}

#[tokio::test]
async fn register_invalid_email() {
    let (app, _state) = build_app().await;
    let req = post_json(
        "/api/v1/auth/register",
        r#"{"nickname":"abcd","email":"not-an-email","password":"abc123"}"#,
    );
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNPROCESSABLE_ENTITY);
    assert_eq!(json["code"], 422);
    assert_eq!(json["message"], "邮箱格式错误");
}

#[tokio::test]
async fn register_short_nickname() {
    let (app, _state) = build_app().await;
    let req = post_json(
        "/api/v1/auth/register",
        r#"{"nickname":"a","email":"x@y.com","password":"abc123"}"#,
    );
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNPROCESSABLE_ENTITY);
    assert_eq!(json["message"], "昵称长度不能小于2位");
}

#[tokio::test]
async fn refresh_without_token_unauthorized() {
    let (app, _state) = build_app().await;
    let req = Request::builder()
        .method("POST")
        .uri("/api/v1/auth/refresh")
        .body(Body::empty())
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
    assert_eq!(json["code"], 401);
}

#[tokio::test]
async fn update_without_token_unauthorized() {
    let (app, _state) = build_app().await;
    let req = post_json(
        "/api/v1/user/update",
        r#"{"nickname":"abc","avatar_url":""}"#,
    );
    let (status, _json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
}

/// 完整链路: 注册 → 登录 → 刷新 → 更新 (并清理测试数据)
#[tokio::test]
async fn auth_roundtrip() {
    let (app, state) = build_app().await;
    let email = format!(
        "test_{}@example.com",
        std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos()
    );

    // 注册
    let req = post_json(
        "/api/v1/auth/register",
        format!(r#"{{"nickname":"testuser","email":"{email}","password":"abc123"}}"#),
    );
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK, "register should succeed: {json}");
    // 成功响应为裸实体 JSON, 无 code/message 包裹
    let uid = json["user"]["uid"].as_i64().unwrap() as u32;
    let token = json["token"].as_str().unwrap().to_string();
    assert!(!token.is_empty(), "register should return token");

    // 登录
    let req = post_json(
        "/api/v1/auth/login",
        format!(r#"{{"email":"{email}","password":"abc123"}}"#),
    );
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK);
    let login_token = json["token"].as_str().unwrap().to_string();
    assert!(!login_token.is_empty());

    // 刷新 (带 JWT)
    let req = Request::builder()
        .method("POST")
        .uri("/api/v1/auth/refresh")
        .header(header::AUTHORIZATION, format!("Bearer {login_token}"))
        .body(Body::empty())
        .unwrap();
    let (status, _json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK);

    // 更新 (带 JWT)
    let req = Request::builder()
        .method("POST")
        .uri("/api/v1/user/update")
        .header(header::AUTHORIZATION, format!("Bearer {login_token}"))
        .header(header::CONTENT_TYPE, "application/json")
        .body(Body::from(
            r#"{"nickname":"updated_name","avatar_url":"https://example.com/a.png"}"#,
        ))
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(json["nickname"], "updated_name");

    // 清理测试数据
    let db = &state.db;
    user_third_party::Entity::delete_many()
        .filter(user_third_party::Column::Uid.eq(uid))
        .exec(db)
        .await
        .unwrap();
    user_extend::Entity::delete_many()
        .filter(user_extend::Column::Uid.eq(uid))
        .exec(db)
        .await
        .unwrap();
    user_model::Entity::delete_many()
        .filter(user_model::Column::Id.eq(uid))
        .exec(db)
        .await
        .unwrap();
}

#[tokio::test]
async fn login_wrong_password_returns_422() {
    let (app, state) = build_app().await;
    let email = unique_email();
    let (uid, _token) = register_user(&app, &email).await;
    // 错误密码登录 → 422 "邮箱或密码错误"
    let req = post_json(
        "/api/v1/auth/login",
        format!(r#"{{"email":"{email}","password":"wrong123"}}"#),
    );
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNPROCESSABLE_ENTITY);
    assert_eq!(json["message"], "邮箱或密码错误");
    cleanup_user(&state, uid).await;
}

#[tokio::test]
async fn register_duplicate_email_returns_422() {
    let (app, state) = build_app().await;
    let email = unique_email();
    let (uid, _token) = register_user(&app, &email).await;
    // 同邮箱重复注册 → 422 "邮箱已存在"
    let req = post_json(
        "/api/v1/auth/register",
        format!(r#"{{"nickname":"another","email":"{email}","password":"abc123"}}"#),
    );
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNPROCESSABLE_ENTITY);
    assert_eq!(json["message"], "邮箱已存在");
    cleanup_user(&state, uid).await;
}

#[tokio::test]
async fn update_with_invalid_token_unauthorized() {
    let (app, _state) = build_app().await;
    // 无效 JWT → 401
    let req = Request::builder()
        .method("POST")
        .uri("/api/v1/user/update")
        .header(header::AUTHORIZATION, "Bearer invalid.token.value")
        .header(header::CONTENT_TYPE, "application/json")
        .body(Body::from(r#"{"nickname":"abc","avatar_url":""}"#))
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
    assert_eq!(json["code"], 401);
}

#[tokio::test]
async fn user_info_ok() {
    let (app, state) = build_app().await;
    let email = unique_email();
    let (uid, _token) = register_user(&app, &email).await;
    // 查询用户信息 → 200, 公开视图不含 email 等敏感字段, 校验 uid 与昵称
    let req = Request::builder()
        .uri(format!("/api/v1/user/info?uid={uid}"))
        .body(Body::empty())
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(json["uid"], uid as i64);
    assert_eq!(json["nickname"], "testuser");
    // 敏感字段不应出现在公开接口响应中
    assert!(json.get("email").is_none());
    assert!(json.get("mobile").is_none());
    assert!(json.get("account").is_none());
    cleanup_user(&state, uid).await;
}

#[tokio::test]
async fn user_list_ok() {
    let (app, state) = build_app().await;
    let email = unique_email();
    let (uid, _token) = register_user(&app, &email).await;
    // 分页查询用户列表 → 200, total 至少包含刚注册的用户
    let req = Request::builder()
        .uri("/api/v1/user/list?page=1&page_size=12")
        .body(Body::empty())
        .unwrap();
    let (status, json) = send(&app, req).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(json["page"], 1);
    assert!(json["total"].as_i64().unwrap() >= 1);
    cleanup_user(&state, uid).await;
}
