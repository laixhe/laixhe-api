//! 健康检查控制器 (新增)
//!
//! GET /api/v1/health: 探测服务与数据库是否正常, 成功返回裸实体 JSON (与 Go 版一致)

use axum::extract::State;
use axum::Json;
use sea_orm::DatabaseConnection;
use serde::Serialize;
use std::sync::{LazyLock, Mutex};
use std::time::{Duration, Instant};
use utoipa::ToSchema;

use crate::error::ApiError;
use crate::state::AppState;

/// 健康检查响应体
#[derive(Debug, Clone, Serialize, ToSchema)]
pub struct HealthResponse {
    /// 服务状态 (固定 "ok"; DB 异常不输出 degraded, 直接走 503 统一错误体)
    pub status: &'static str,
    /// 数据库状态 (固定 "up"; DB 异常不输出 down, 直接走 503 统一错误体)
    pub database: &'static str,
    /// 服务版本
    pub version: &'static str,
    /// 服务启动时间 (服务器本地时区, 格式 Y-m-d H:i:s, 与 Go/PHP 版一致)
    pub started_at: String,
    /// 当前时间 (服务器本地时区, 格式 Y-m-d H:i:s)
    pub now: String,
}

/// GET /health 健康检查
///
/// 通过连接级 ping 探测数据库, 成功返回 200 + 健康信息。
/// 数据库不可用时返回 503 + 统一错误格式, 便于负载均衡探活。
#[utoipa::path(
    get,
    path = "/api/v1/health",
    tag = "Health",
    summary = "健康检查",
    responses(
        (status = 200, description = "OK", body = HealthResponse),
        (status = 503, description = "Service Unavailable", body = crate::docs::Error)
    )
)]
pub async fn health(
    State(state): State<AppState>,
) -> Result<Json<HealthResponse>, ApiError> {
    let resp = HealthResponse {
        status: "ok",
        database: "up",
        version: env!("CARGO_PKG_VERSION"),
        started_at: STARTED_AT.to_string(),
        now: now_local(),
    };
    if !ping_db(&state.db).await {
        // 数据库不可用 → 503 统一错误格式 (对齐 Go 版 fiber.NewError)
        return Err(ApiError {
            code: 503,
            message: "database unavailable".to_string(),
        });
    }
    Ok(Json(resp))
}

/// 当前服务器本地时间, 格式 "Y-m-d H:i:s" (对齐 Go time.DateTime / PHP date('Y-m-d H:i:s'))
fn now_local() -> String {
    jiff::Zoned::now()
        .strftime("%Y-%m-%d %H:%M:%S")
        .to_string()
}

/// 服务启动时间 (进程级常量, 服务器本地时区)
static STARTED_AT: std::sync::LazyLock<String> = std::sync::LazyLock::new(now_local);

/// 数据库探活结果缓存 (TTL 内复用, 避免被负载均衡秒级高频探活持续打数据库)
static DB_STATUS: LazyLock<Mutex<Option<(Instant, bool)>>> = LazyLock::new(|| Mutex::new(None));

/// 探活结果缓存有效期: 数据库短暂抖动时 LB 判活有 5 秒缓冲, 避免探活抖动误摘节点
const DB_PING_TTL: Duration = Duration::from_secs(5);

/// 探测数据库连接 (连接级 ping, 结果按 TTL 缓存)
async fn ping_db(db: &DatabaseConnection) -> bool {
    let now = Instant::now();
    // 命中缓存则直接复用 (锁在读后立即释放, 不跨 await 持锁)
    if let Some((at, ok)) = *crate::sync::lock_unpoison(&DB_STATUS) {
        if now.duration_since(at) < DB_PING_TTL {
            return ok;
        }
    }
    let ok = db.ping().await.is_ok();
    *crate::sync::lock_unpoison(&DB_STATUS) = Some((now, ok));
    ok
}
