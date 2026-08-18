//! 用户控制器 (对应 Go 项目 app/controllers/user.go)

use axum::extract::State;
use axum::Json;

use crate::app::controllers::{JsonBody, QueryParams};
use crate::app::entity::user::{
    User, UserInfoRequest, UserListRequest, UserListResponse, UserUpdateRequest,
};
use crate::app::services;
use crate::app::util::validate::{validate_avatar_url, validate_nickname};
use crate::error::ApiError;
use crate::log_elapsed;
use crate::logger::Timer;
use crate::middleware::jwt::JwtAuth;
use crate::state::AppState;

/// 更新用户信息 (成功返回裸实体 JSON, 与 Go 版 ctx.JSON(resp) 一致)
#[utoipa::path(
    post,
    path = "/api/v1/user/update",
    tag = "User",
    summary = "更新用户信息",
    security(("BearerAuth" = [])),
    request_body = UserUpdateRequest,
    responses(
        (status = 200, description = "OK", body = User),
        (status = 401, description = "未授权", body = crate::docs::Error),
        (status = 400, description = "请求格式错误", body = crate::docs::Error),
        (status = 422, description = "参数错误", body = crate::docs::Error),
        (status = 500, description = "Internal Server Error", body = crate::docs::Error)
    )
)]
pub async fn update(
    State(state): State<AppState>,
    JwtAuth(claims): JwtAuth,
    JsonBody(mut req): JsonBody<UserUpdateRequest>,
) -> Result<Json<User>, ApiError> {
    let start = Timer::new();
    tracing::debug!(uid = claims.uid, "收到更新用户请求");
    // 参数校验 (昵称 / 头像地址格式)
    validate_nickname(&req.nickname)?;
    validate_avatar_url(&req.avatar_url)?;
    // Uid 由 JWT 提供
    req.uid = claims.uid;
    let resp = services::user::update(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        debug,
        uid = claims.uid,
        "更新用户接口处理完成"
    );
    Ok(Json(resp))
}

/// 获取用户信息 (公开接口, 返回完整用户实体 JSON)
#[utoipa::path(
    get,
    path = "/api/v1/user/info",
    tag = "User",
    summary = "获取用户信息",
    params(("uid" = i32, Query, description = "用户id")),
    responses(
        (status = 200, description = "OK", body = User),
        (status = 400, description = "请求格式错误", body = crate::docs::Error),
        (status = 422, description = "参数错误", body = crate::docs::Error),
        (status = 500, description = "Internal Server Error", body = crate::docs::Error)
    )
)]
pub async fn info(
    State(state): State<AppState>,
    QueryParams(req): QueryParams<UserInfoRequest>,
) -> Result<Json<User>, ApiError> {
    let start = Timer::new();
    tracing::debug!(uid = req.uid, "收到获取用户信息请求");
    if req.uid <= 0 {
        return Err(ApiError::param_error("无效的用户ID"));
    }
    let resp = services::user::info(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        debug,
        uid = req.uid,
        "获取用户信息接口处理完成"
    );
    Ok(Json(resp))
}

/// 获取用户列表 (成功返回裸实体 JSON)
#[utoipa::path(
    get,
    path = "/api/v1/user/list",
    tag = "User",
    summary = "获取用户列表",
    params(
        ("page" = Option<i32>, Query, description = "分页-当前页(默认 1)"),
        ("page_size" = Option<i32>, Query, description = "分页-每页数量(默认 12)")
    ),
    responses(
        (status = 200, description = "OK", body = UserListResponse),
        (status = 400, description = "请求格式错误", body = crate::docs::Error),
        (status = 422, description = "参数错误", body = crate::docs::Error),
        (status = 500, description = "Internal Server Error", body = crate::docs::Error)
    )
)]
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
    tracing::debug!(
        page = req.page,
        page_size = req.page_size,
        "收到获取用户列表请求"
    );
    let resp = services::user::list(&state, &req).await?;
    log_elapsed!(
        start,
        total_ms,
        debug,
        page = req.page,
        page_size = req.page_size,
        "获取用户列表接口处理完成"
    );
    Ok(Json(resp))
}
