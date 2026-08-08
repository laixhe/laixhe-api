//! 鉴权业务 (对应 Go 项目 app/services/auth.go)

use sea_orm::ActiveValue::Set;

use crate::app::entity::auth::{
    AuthLoginRequest, AuthLoginResponse, AuthRefreshRequest, AuthRefreshResponse,
    AuthRegisterRequest, AuthRegisterResponse,
};
use crate::app::entity::user::User;
use crate::app::models::user as user_model;
use crate::app::models::{USER_SEX_UNKNOWN, USER_STATE_NORMAL, USER_TYPE_ORDINARY};
use crate::error::ApiError;
use crate::log_elapsed;
use crate::logger::Timer;
use crate::middleware::jwt::{gen_token, new_jwt_claims};
use crate::state::AppState;

/// 注册
pub async fn register(
    state: &AppState,
    req: &AuthRegisterRequest,
) -> Result<AuthRegisterResponse, ApiError> {
    let start = Timer::new();
    tracing::info!(email = %req.email, "开始注册用户");

    // 先检查邮箱是否已注册，避免无效的 bcrypt 计算 (仅查 id 列, 对应 Go Select("id"))
    let step = Timer::new();
    if user_model::find_id_by_email(&state.db, &req.email)
        .await?
        .is_some()
    {
        log_elapsed!(step, elapsed_ms, debug, email = %req.email, "邮箱已存在，注册失败");
        return Err(ApiError::param_error("邮箱已存在"));
    }
    log_elapsed!(step, elapsed_ms, debug, email = %req.email, "邮箱未注册，通过检查");

    // 密码 bcrypt 加密 (cost=10, 对齐 Go 的 bcrypt.DefaultCost)
    let step = Timer::new();
    let password = bcrypt::hash(&req.password, 10).map_err(|e| {
        // 收敛为固定文案, 避免错误细节泄露 (对齐 Go fork DefaultErrorHandler); 原始错误记服务端日志
        tracing::error!(%e, "bcrypt 密码加密失败");
        ApiError::internal("internal server error")
    })?;
    log_elapsed!(step, elapsed_ms, debug, email = %req.email, "bcrypt 密码加密完成");

    let now = user_model::now_local_naive();
    let account = xid::new().to_string();
    let user = user_model::ActiveModel {
        type_id: Set(USER_TYPE_ORDINARY),
        account: Set(account.clone()),
        mobile: Set(String::new()),
        nickname: Set(req.nickname.clone()),
        email: Set(req.email.clone()),
        password: Set(password.clone()),
        avatar_url: Set(String::new()),
        sex: Set(USER_SEX_UNKNOWN),
        states: Set(USER_STATE_NORMAL),
        created_at: Set(now),
        updated_at: Set(now),
        ..Default::default()
    };

    // 事务创建用户 + 扩展信息 + 第三方关联
    let step = Timer::new();
    let uid = user_model::create_user(&state.db, user).await?;
    log_elapsed!(
        step,
        elapsed_ms,
        debug,
        email = %req.email,
        uid,
        "事务创建用户成功 (user + user_extend + user_third_party)"
    );

    // 签发 JWT
    let step = Timer::new();
    let claims = new_jwt_claims(uid, state.config.jwt.expire_time);
    let token = gen_token(&state.config.jwt.secret_key, &claims)?;
    log_elapsed!(step, elapsed_ms, debug, email = %req.email, uid, "JWT 签发完成");

    // 组装响应模型 (字段与写入一致)
    let model = user_model::Model {
        id: uid,
        type_id: USER_TYPE_ORDINARY,
        account,
        mobile: String::new(),
        email: req.email.clone(),
        nickname: req.nickname.clone(),
        avatar_url: String::new(),
        password,
        sex: USER_SEX_UNKNOWN,
        states: USER_STATE_NORMAL,
        created_at: now,
        updated_at: now,
    };
    log_elapsed!(start, total_ms, debug, email = %req.email, uid, "注册完成");
    Ok(AuthRegisterResponse {
        token,
        user: User::from_model(&model, "", ""),
    })
}

/// 登录
pub async fn login(
    state: &AppState,
    req: &AuthLoginRequest,
) -> Result<AuthLoginResponse, ApiError> {
    let start = Timer::new();
    tracing::info!(email = %req.email, "开始登录");

    // 按邮箱查询用户
    let step = Timer::new();
    let user = match user_model::find_by_email(&state.db, &req.email).await? {
        Some(user) => user,
        None => {
            log_elapsed!(step, elapsed_ms, debug, email = %req.email, "用户不存在，登录失败");
            return Err(ApiError::param_error("邮箱或密码错误"));
        }
    };
    log_elapsed!(step, elapsed_ms, debug, email = %req.email, uid = user.id, "查询用户成功");

    if user.states != USER_STATE_NORMAL {
        tracing::warn!(
            email = %req.email,
            uid = user.id,
            states = user.states,
            "用户状态异常，拒绝登录"
        );
        return Err(ApiError::unauthorized());
    }

    // 校验密码
    let step = Timer::new();
    if !bcrypt::verify(&req.password, &user.password).map_err(|e| {
        // 收敛为固定文案, 避免错误细节泄露 (对齐 Go fork DefaultErrorHandler); 原始错误记服务端日志
        tracing::error!(%e, "bcrypt 密码校验失败");
        ApiError::internal("internal server error")
    })? {
        log_elapsed!(step, elapsed_ms, debug, email = %req.email, uid = user.id, "密码校验失败，登录失败");
        return Err(ApiError::param_error("邮箱或密码错误"));
    }
    log_elapsed!(step, elapsed_ms, debug, email = %req.email, uid = user.id, "密码校验通过");

    // 签发 JWT
    let step = Timer::new();
    let claims = new_jwt_claims(user.id, state.config.jwt.expire_time);
    let token = gen_token(&state.config.jwt.secret_key, &claims)?;
    log_elapsed!(step, elapsed_ms, debug, email = %req.email, uid = user.id, "JWT 签发完成");

    log_elapsed!(start, total_ms, debug, email = %req.email, uid = user.id, "登录成功");
    Ok(AuthLoginResponse {
        token,
        user: User::from_model(&user, "", ""),
    })
}

/// 刷新Jwt
pub async fn refresh(
    state: &AppState,
    req: &AuthRefreshRequest,
) -> Result<AuthRefreshResponse, ApiError> {
    let start = Timer::new();
    tracing::info!(uid = req.uid, "开始刷新Jwt");
    let user = match user_model::find_by_id(&state.db, req.uid).await? {
        Some(user) => user,
        None => return Err(ApiError::unauthorized()),
    };
    if user.states != USER_STATE_NORMAL {
        return Err(ApiError::unauthorized());
    }
    let claims = new_jwt_claims(user.id, state.config.jwt.expire_time);
    let token = gen_token(&state.config.jwt.secret_key, &claims)?;
    log_elapsed!(start, total_ms, debug, uid = req.uid, "刷新Jwt完成");
    Ok(AuthRefreshResponse {
        token,
        user: User::from_model(&user, "", ""),
    })
}
