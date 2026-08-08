//! 业务服务逻辑 (对应 Go 项目 app/services)

pub mod auth;
pub mod user;

use crate::app::models::config_common as config_common_model;
use crate::app::models::CONFIG_COMMON_ENV;
use crate::state::AppState;

/// 从数据库 config_common 表加载运行时配置（如环境标识 env）
/// 对应 Service.initConfigCommon
pub async fn init_config_common(state: &AppState) {
    match config_common_model::list(&state.db, &[]).await {
        Ok(configs) => {
            for c in configs {
                if c.key == CONFIG_COMMON_ENV {
                    // 锁中毒时取回数据继续使用 (启动期单次写入, 数据一致性可接受; 与 rate_limit.rs 处理一致)
                    let mut common = state
                        .common
                        .lock()
                        .unwrap_or_else(|poisoned| poisoned.into_inner());
                    common.env = c.value;
                }
            }
        }
        Err(e) => tracing::error!("initConfigCommon failed: {:?}", e),
    }
    tracing::debug!("config http={:?}", state.config.http);
    tracing::debug!(
        "config common env={}",
        state
            .common
            .lock()
            .map(|c| c.env.clone())
            .unwrap_or_default()
    );
}
