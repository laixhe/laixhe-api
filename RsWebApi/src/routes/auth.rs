//! 鉴权相关路由 (对应 Go 项目 routers/auth.go)

use axum::middleware::from_fn_with_state;
use axum::routing::post;
use axum::Router;

use crate::app::controllers::auth as ctrl;
use crate::middleware::jwt::jwt_middleware;
use crate::state::AppState;

/// 注册鉴权路由
pub fn router(state: AppState) -> Router<AppState> {
    let auth_router = Router::new()
        .route("/auth/register", post(ctrl::register)) // 注册
        .route("/auth/login", post(ctrl::login)) // 登录
        // 需要 JWT 鉴权的路由
        .route(
            "/auth/refresh",
            post(ctrl::refresh).route_layer(from_fn_with_state(state, jwt_middleware)), // 刷新 Jwt
        );
    auth_router
}
