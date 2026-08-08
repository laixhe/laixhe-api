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
                    // 锁中毒时取回数据继续使用 (启动期单次写入, 数据一致性可接受)
                    let mut common = crate::sync::lock_unpoison(&state.common);
                    common.env = c.value;
                }
            }
        }
        Err(e) => tracing::error!("initConfigCommon failed: {:?}", e),
    }
    tracing::debug!("config http={:?}", state.config.http);
    tracing::debug!(
        "config common env={}",
        crate::sync::lock_unpoison(&state.common).env
    );
}
