//! 请求/响应实体 (对应 Go 项目 app/entity)
//!
//! 本目录存放接口层的请求/响应 DTO; 数据库 ORM 实体在 `crate::app::models`。
//! 注意: 该命名与 Rust 社区惯例 (sea-orm 生态通常把 ORM 实体放 `entity/`) 相反,
//! 为对齐 Go 原版分层, 详见 models/mod.rs 的说明。

use serde::Deserialize;

pub mod auth;
pub mod user;

/// serde 反序列化辅助: JSON null 视为缺失, 反序列化为空字符串
///
/// 与 Go 端 encoding/json 的 null→零值、PHP 端 null 走业务校验的语义一致
/// (serde 的 `#[serde(default)]` 只处理"字段缺失", 不处理"显式 null")。
pub fn de_null_to_empty<'de, D>(deserializer: D) -> Result<String, D::Error>
where
    D: serde::Deserializer<'de>,
{
    let opt = Option::<String>::deserialize(deserializer)?;
    Ok(opt.unwrap_or_default())
}
