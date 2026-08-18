//! 应用聚合 (对应 Go 项目 app/app.go)
//!
//! Controller / Service 层通过 AppState 共享 Config、数据库连接与运行时配置。

pub mod controllers;
pub mod entity;
pub mod models;
pub mod services;
pub mod util;
