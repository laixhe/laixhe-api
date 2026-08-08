//! 鉴权控制器 (对应 Go 项目 app/controllers/auth.go)

use axum::extract::State;
use axum::Json;

use crate::app::controllers::JsonBody;
use crate::app::entity::auth::{
    AuthLoginRequest, AuthLoginResponse, AuthRefreshRequest, AuthRefreshResponse,
    AuthRegisterRequest, AuthRegisterResponse,
};
use crate::app::services;
use crate::app::util::regexp::{is_email, is_password};
use crate::app::util::validate::validate_nickname;
use crate::error::ApiError;
use crate::log_elapsed;
use crate::logger::Timer;
use crate::middleware::jwt::JwtAuth;
use crate::state::AppState;

/// 验证邮箱和密码格式
fn validate_email_and_password(email: &str, password: &str) -> Result<(), ApiError> {
    if !is_email(email) {
        return Err(ApiError::param_error("邮箱格式错误"));
    }
    if password.len() < 6 {
        return Err(ApiError::param_error("密码长度不能小于6位"));
    }
    if !is_password(password) {
        return Err(ApiError::param_error(
            "密码格式错误，只能包含字母 数字 _ @ $",
        ));
    }
    Ok(())
}

/// 注册 (成功返回裸实体 JSON, 与 Go 版 ctx.JSON(resp) 一致)
pub async fn register(
    State(state): State<AppState>,
    JsonBody(req): JsonBody<AuthRegisterRequest>,
) -> Result<Json<AuthRegisterResponse>, ApiError> {
    let start = Timer::new();
    tracing::info!(email = %req.email, "收到注册请求");

    // 参数校验 (昵称 / 邮箱 / 密码)
    let step = Timer::new();
    validate_nickname(&req.nickname)?;
    validate_email_and_password(&req.email, &req.password)?;
    log_elapsed!(step, elapsed_ms, info, email = %req.email, "注册参数校验通过");

    let resp = services::auth::register(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        info,
        email = %req.email,
        uid = resp.user.uid,
        "注册接口处理完成"
    );
    Ok(Json(resp))
}

/// 登录 (成功返回裸实体 JSON)
pub async fn login(
    State(state): State<AppState>,
    JsonBody(req): JsonBody<AuthLoginRequest>,
) -> Result<Json<AuthLoginResponse>, ApiError> {
    let start = Timer::new();
    tracing::info!(email = %req.email, "收到登录请求");

    // 参数校验
    let step = Timer::new();
    validate_email_and_password(&req.email, &req.password)?;
    log_elapsed!(step, elapsed_ms, info, email = %req.email, "登录参数校验通过");

    let resp = services::auth::login(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        info,
        email = %req.email,
        uid = resp.user.uid,
        "登录接口处理完成"
    );
    Ok(Json(resp))
}

/// 刷新Jwt (成功返回裸实体 JSON)
pub async fn refresh(
    State(state): State<AppState>,
    JwtAuth(claims): JwtAuth,
) -> Result<Json<AuthRefreshResponse>, ApiError> {
    let start = Timer::new();
    tracing::info!(uid = claims.uid, "收到刷新Jwt请求");
    let req = AuthRefreshRequest { uid: claims.uid };
    let resp = services::auth::refresh(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        info,
        uid = claims.uid,
        "刷新Jwt接口处理完成"
    );
    Ok(Json(resp))
}
