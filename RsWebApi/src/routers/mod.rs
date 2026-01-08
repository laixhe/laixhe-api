use axum::{
    routing::{get, post},
    Router,
};

use crate::{controllers::auth, controllers::user};

pub async fn init() -> Router {
    // 鉴权
    let auth_router = Router::new()
        .route("/register", post(auth::register))
        .route("/login", post(auth::login))
        .route("/refresh", post(auth::refresh));
    // 用户
    let user_router = Router::new()
        .route("/info", get(user::info))
        .route("/list", get(user::list))
        .route("/update", post(user::update));
    //
    return Router::new()
        .route("/", get(|| async { "☺ webapi to Rust" }))
        .nest("/api/auth", auth_router)
        .nest("/api/user", user_router);
}
