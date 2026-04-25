use axum::{
    routing::{get, post},
    Router,
};
use tower_http::trace::TraceLayer;

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
    Router::new()
        .route("/", get(|| async { "☺ WebApi to Rust" }))
        .nest("/api/auth", auth_router)
        .nest("/api/user", user_router)
        .layer(TraceLayer::new_for_http())
}
