//! 用户控制器 (对应 Go 项目 app/controllers/user.go)

use axum::extract::State;
use axum::Json;

use crate::app::controllers::{JsonBody, QueryParams};
use crate::app::entity::user::{
    User, UserInfoRequest, UserListRequest, UserListResponse, UserPublic, UserUpdateRequest,
};
use crate::app::services;
use crate::app::util::validate::{validate_avatar_url, validate_nickname};
use crate::error::ApiError;
use crate::log_elapsed;
use crate::logger::Timer;
use crate::middleware::jwt::JwtAuth;
use crate::state::AppState;

/// 更新用户信息 (成功返回裸实体 JSON, 与 Go 版 ctx.JSON(resp) 一致)
pub async fn update(
    State(state): State<AppState>,
    JwtAuth(claims): JwtAuth,
    JsonBody(mut req): JsonBody<UserUpdateRequest>,
) -> Result<Json<User>, ApiError> {
    let start = Timer::new();
    tracing::info!(uid = claims.uid, "收到更新用户请求");
    // 参数校验 (昵称 / 头像地址格式)
    validate_nickname(&req.nickname)?;
    validate_avatar_url(&req.avatar_url)?;
    // Uid 由 JWT 提供
    req.uid = claims.uid;
    let resp = services::user::update(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        info,
        uid = claims.uid,
        "更新用户接口处理完成"
    );
    Ok(Json(resp))
}

/// 获取用户信息 (公开接口, 成功返回脱敏后的公开实体 JSON)
pub async fn info(
    State(state): State<AppState>,
    QueryParams(req): QueryParams<UserInfoRequest>,
) -> Result<Json<UserPublic>, ApiError> {
    let start = Timer::new();
    tracing::info!(uid = req.uid, "收到获取用户信息请求");
    if req.uid == 0 {
        return Err(ApiError::param_error("无效的用户ID"));
    }
    let resp = services::user::info(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        info,
        uid = req.uid,
        "获取用户信息接口处理完成"
    );
    Ok(Json(resp))
}

/// 获取用户列表 (成功返回裸实体 JSON)
pub async fn list(
    State(state): State<AppState>,
    QueryParams(mut req): QueryParams<UserListRequest>,
) -> Result<Json<UserListResponse>, ApiError> {
    let start = Timer::new();
    // 分页参数缺省值
    if req.page <= 0 {
        req.page = 1;
    }
    if req.page_size <= 0 {
        req.page_size = 12;
    }
    // 上限保护: 防止超大 page_size 触发全量查询 (对齐 Go 版本)
    if req.page_size > 100 {
        req.page_size = 100;
    }
    tracing::info!(
        page = req.page,
        page_size = req.page_size,
        "收到获取用户列表请求"
    );
    let resp = services::user::list(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        info,
        page = req.page,
        page_size = req.page_size,
        "获取用户列表接口处理完成"
    );
    Ok(Json(resp))
}
