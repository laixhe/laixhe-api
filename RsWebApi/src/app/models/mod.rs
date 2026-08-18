//! 数据模型 (对应 Go 项目 app/models)
//!
//! 使用 sea-orm 定义数据库实体，实体类型 (Entity/ActiveModel/Model/Column)
//! 直接位于各文件模块层级，通过 `crate::app::models::user` 等路径访问。
//!
//! 注意: 本目录命名与 Rust 社区惯例 (sea-orm 生态通常把 ORM 实体放 `entity/`) 相反,
//! 这是为对齐 Go 原版 `app/models` (表结构) 与 `app/entity` (接口 DTO) 的分层,
//! 跨语言对照阅读时保持一致; 若脱离本仓库单独阅读, 请先看本说明避免混淆。

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
