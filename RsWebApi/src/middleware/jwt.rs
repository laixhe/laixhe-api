//! JWT 令牌载荷 / 签发 / 校验 / 中间件 (对应 Go 项目 core/middlewares/jwt.go)
//!
//! 使用 jiff 生成时间戳，jsonwebtoken 完成 HS256 签名与验签。

use axum::extract::{FromRequestParts, Request, State};
use axum::http::header;
use axum::http::request::Parts;
use axum::middleware::Next;
use axum::response::Response;
use jsonwebtoken::{decode, encode, Algorithm, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};
use std::sync::OnceLock;

use crate::error::ApiError;
use crate::log_elapsed;
use crate::logger::Timer;
use crate::middleware::request_log::RequestId;
use crate::state::AppState;

/// JWT 令牌载荷，存储用户 UID
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct JwtClaims {
    pub uid: i32,
    /// 过期时间 (Unix 秒)
    pub exp: i64,
    /// 发布时间 (Unix 秒)
    pub iat: i64,
    /// 生效时间 (Unix 秒)
    pub nbf: i64,
}

/// 创建 JWT 载荷并设置过期时间、发布时间、生效时间
pub fn new_jwt_claims(uid: i32, expire_time: i64) -> JwtClaims {
    let now = jiff::Timestamp::now();
    let now_secs = now.as_second();
    JwtClaims {
        uid,
        exp: now_secs + expire_time, // 过期时间
        iat: now_secs,               // 发布时间
        nbf: now_secs,               // 生效时间
    }
}

/// 缓存的 JWT 密钥与校验配置 (构建一次, 避免每次请求重建堆分配;
/// 密钥来自 config.jwt.secret_key, 运行期不变; 只缓存首个 secret_key, 更换密钥不刷新缓存)
struct JwtKeys {
    validation: Validation,
    decoding_key: DecodingKey,
    encoding_key: EncodingKey,
}

/// OnceLock 只初始化一次: 认证热路径上读取为原子加载, 无锁竞争
static JWT_KEYS: OnceLock<JwtKeys> = OnceLock::new();

/// 获取缓存的 JWT 密钥对象, 首次调用时按 secret_key 构建 (线程安全)
fn jwt_keys(secret_key: &str) -> &'static JwtKeys {
    JWT_KEYS.get_or_init(|| {
        let mut validation = Validation::new(Algorithm::HS256);
        // 校验 nbf (生效时间), 对齐 Go jwtv5 默认行为
        validation.validate_nbf = true;
        // 对齐 Go jwtv5: leeway=0 (无签发/过期宽限期), exp 可选 (缺失时不强制校验)
        validation.leeway = 0;
        validation.set_required_spec_claims(&[] as &[String]);
        JwtKeys {
            validation,
            decoding_key: DecodingKey::from_secret(secret_key.as_bytes()),
            encoding_key: EncodingKey::from_secret(secret_key.as_bytes()),
        }
    })
}

/// 签发 JWT 令牌
pub fn gen_token(secret_key: &str, claims: &JwtClaims) -> Result<String, ApiError> {
    let keys = jwt_keys(secret_key);
    encode(
        &Header::new(Algorithm::HS256),
        claims,
        &keys.encoding_key,
    )
    .map_err(|e| {
        // 收敛为固定文案, 避免错误细节泄露 (对齐 Go fork DefaultErrorHandler); 原始错误记服务端日志
        tracing::error!(%e, "JWT 签发失败");
        ApiError::internal("internal server error")
    })
}

/// 校验并解析 JWT 令牌
pub fn parse_token(secret_key: &str, token: &str) -> Result<JwtClaims, ApiError> {
    let keys = jwt_keys(secret_key);
    let data = decode::<JwtClaims>(token, &keys.decoding_key, &keys.validation)
        .map_err(|_| ApiError::unauthorized())?;
    // uid 从 1 起, 0 视为无效载荷: 防御伪造 {"uid":0} 的 token (对齐 Go 侧 Uid > 0 判断)
    if data.claims.uid == 0 {
        return Err(ApiError::unauthorized());
    }
    Ok(data.claims)
}

/// 强制 JWT 校验中间件，无 Token 或校验失败返回 401
pub async fn jwt_middleware(
    State(state): State<AppState>,
    req: Request,
    next: Next,
) -> Result<Response, ApiError> {
    let start = Timer::new();
    // 仅 debug 级才为日志字段分配字符串, 避免默认 info 级下每请求白做堆分配
    // (req 在 next.run 时被 move, &str 借用无法跨 await 存活; warn 分支在 move 前内联 &str)
    let method = if tracing::level_enabled!(tracing::Level::DEBUG) {
        req.method().as_str().to_string()
    } else {
        String::new()
    };
    let path = if tracing::level_enabled!(tracing::Level::DEBUG) {
        req.uri().path().to_string()
    } else {
        String::new()
    };
    // 由全局请求日志中间件注入的 Request ID, 用于日志关联;
    // 与上方 method/path 同理, 仅 debug 级才 clone: 默认 info 级下 debug 日志不输出,
    // 常见成功请求路径不为此做堆分配 (缺 Token 的 warn 分支在下方按需单独 clone)
    let request_id = if tracing::level_enabled!(tracing::Level::DEBUG) {
        req.extensions()
            .get::<RequestId>()
            .map(|r| r.0.clone())
            .unwrap_or_default()
    } else {
        String::new()
    };
    tracing::debug!(
        request_id = %request_id,
        method = %method,
        path = %path,
        "收到请求，开始 JWT 鉴权"
    );

    let auth = req
        .headers()
        .get(header::AUTHORIZATION)
        .and_then(|v| v.to_str().ok())
        .unwrap_or("");
    // 严格校验大写 "Bearer " 前缀 (与 Go/PHP/TS 端一致, 避免小写变体被误解析)
    let token = auth.strip_prefix("Bearer ").unwrap_or("");
    if token.is_empty() {
        // 该 warn 分支在默认 info 级也会输出, 因此在此按需 clone Request ID,
        // 不占用上方常见的成功请求路径 (见 request_id 注释)
        let request_id = req
            .extensions()
            .get::<RequestId>()
            .map(|r| r.0.clone())
            .unwrap_or_default();
        log_elapsed!(
            start,
            elapsed_ms,
            warn,
            request_id = %request_id,
            method = %req.method().as_str(), // req 尚未 move, 直接内联 &str 避免额外分配
            path = %req.uri().path(),
            "缺少 Authorization Token，拒绝访问"
        );
        return Err(ApiError::unauthorized());
    }

    // 校验并解析 JWT 令牌
    let step = Timer::new();
    let claims = parse_token(&state.config.jwt.secret_key, token)?;
    log_elapsed!(
        step,
        elapsed_ms,
        debug,
        request_id = %request_id,
        method = %method,
        path = %path,
        uid = claims.uid,
        "JWT 校验通过"
    );

    // 将载荷放入 request extensions，供 JwtAuth 提取器使用
    let mut req = req;
    req.extensions_mut().insert(claims);
    let resp = next.run(req).await;
    log_elapsed!(
        start,
        total_ms,
        debug,
        request_id = %request_id,
        method = %method,
        path = %path,
        "请求处理完成 (JWT 中间件)"
    );
    Ok(resp)
}

/// 从请求上下文中提取已验证的 JWT 载荷 (对应 middlewares.GetJwtClaims)
pub struct JwtAuth(pub JwtClaims);

impl FromRequestParts<AppState> for JwtAuth {
    type Rejection = ApiError;

    async fn from_request_parts(
        parts: &mut Parts,
        _state: &AppState,
    ) -> Result<Self, Self::Rejection> {
        parts
            .extensions
            .get::<JwtClaims>()
            .cloned()
            .map(JwtAuth)
            .ok_or_else(ApiError::unauthorized)
    }
}
