//! 数据模型 (对应 Go 项目 app/models)
//!
//! 使用 sea-orm 定义数据库实体，实体类型 (Entity/ActiveModel/Model/Column)
//! 直接位于各文件模块层级，通过 `crate::app::models::user` 等路径访问。

pub mod config_common;
pub mod config_common_key;
pub mod user;
pub mod user_extend;
pub mod user_sex;
pub mod user_states;
pub mod user_third_party;
pub mod user_type;

// 常量与工具函数
pub use config_common_key::CONFIG_COMMON_ENV;
pub use user_sex::USER_SEX_UNKNOWN;
pub use user_states::USER_STATE_NORMAL;
pub use user_type::USER_TYPE_ORDINARY;
