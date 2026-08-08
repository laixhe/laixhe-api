//! 用户业务 (对应 Go 项目 app/services/user.go)

use crate::app::entity::user::{
    User, UserInfoRequest, UserListRequest, UserListResponse, UserUpdateRequest,
};
use crate::app::models::user as user_model;
use crate::app::models::USER_STATE_NORMAL;
use crate::error::ApiError;
use crate::log_elapsed;
use crate::logger::Timer;
use crate::state::AppState;

/// 分页计算 limit / offset (对应 orm.PageOffsetCalculation)
///
/// 多一层 max(1) 兜底: Go 侧由控制器归一化 page/page_size, 此处为双保险 (净效果一致)
fn page_offset_calculation(page: i32, page_size: i32) -> (u64, u64) {
    let page = page.max(1);
    let page_size = page_size.max(1);
    ((page_size) as u64, ((page - 1) * page_size) as u64)
}

/// 更新用户信息
pub async fn update(state: &AppState, req: &UserUpdateRequest) -> Result<User, ApiError> {
    let start = Timer::new();
    tracing::info!(uid = req.uid, "开始更新用户");
    let user = match user_model::find_by_id(&state.db, req.uid).await? {
        Some(user) => user,
        None => return Err(ApiError::param_error("用户不存在")),
    };
    if user.states != USER_STATE_NORMAL {
        return Err(ApiError::unauthorized());
    }
    let resp = User::from_model(&user, &req.nickname, &req.avatar_url);
    // 空字符串不更新，与 Go 的非零字段更新语义一致
    let data = user_model::UserUpdateData {
        nickname: if req.nickname.is_empty() {
            None
        } else {
            Some(req.nickname.clone())
        },
        avatar_url: if req.avatar_url.is_empty() {
            None
        } else {
            Some(req.avatar_url.clone())
        },
        ..Default::default()
    };
    user_model::update_user(&state.db, user.id, &data).await?;
    log_elapsed!(start, total_ms, debug, uid = req.uid, "更新用户完成");
    Ok(resp)
}

/// 获取用户信息
pub async fn info(state: &AppState, req: &UserInfoRequest) -> Result<User, ApiError> {
    let start = Timer::new();
    tracing::info!(uid = req.uid, "开始查询用户信息");
    let user = match user_model::find_by_id(&state.db, req.uid).await? {
        Some(user) => user,
        None => return Err(ApiError::param_error("用户不存在")),
    };
    log_elapsed!(start, total_ms, debug, uid = req.uid, "查询用户信息完成");
    Ok(User::from_model(&user, "", ""))
}

/// 获取用户列表
pub async fn list(state: &AppState, req: &UserListRequest) -> Result<UserListResponse, ApiError> {
    let start = Timer::new();
    tracing::info!(
        page = req.page,
        page_size = req.page_size,
        "开始查询用户列表"
    );
    let (limit, offset) = page_offset_calculation(req.page, req.page_size);
    let (users, total) = user_model::list_user(&state.db, limit, offset).await?;
    let list = users.iter().map(|u| User::from_model(u, "", "")).collect();
    log_elapsed!(
        start,
        total_ms,
        debug,
        page = req.page,
        page_size = req.page_size,
        total,
        "查询用户列表完成"
    );
    Ok(UserListResponse {
        total,
        page: req.page,
        page_size: req.page_size,
        list,
    })
}
