//! 用户相关路由 (对应 Go 项目 routers/user.go)

use axum::middleware::from_fn_with_state;
use axum::routing::{get, post};
use axum::Router;

use crate::app::controllers::user as ctrl;
use crate::middleware::jwt::jwt_middleware;
use crate::state::AppState;

/// 注册用户路由
pub fn router(state: AppState) -> Router<AppState> {
    let user_router = Router::new()
        // 公开用户接口
        .route("/user/info", get(ctrl::info)) // 获取用户信息
        .route("/user/list", get(ctrl::list)) // 获取用户列表
        // 需要 JWT 鉴权的路由
        .route(
            "/user/update",
            post(ctrl::update).route_layer(from_fn_with_state(state, jwt_middleware)), // 更新用户信息
        );
    user_router
}
