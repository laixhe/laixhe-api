//! 用户信息实体 (对应 Go 项目 app/entity/user.go)

use chrono::{Datelike, Timelike};
use serde::{Deserialize, Serialize};

use crate::app::models::user::Model as UserModel;

/// 用户信息
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    /// 用户id
    pub uid: u32,
    /// 类型 (1 - 普通用户)
    pub type_id: i32,
    /// 账号
    pub account: String,
    /// 手机号
    pub mobile: String,
    /// 邮箱
    pub email: String,
    /// 昵称
    pub nickname: String,
    /// 头像地址
    pub avatar_url: String,
    /// 性别 (0 - 未填写, 1 - 男, 2 - 女)
    pub sex: i32,
    /// 状态 (0 - 禁用, 1 - 正常)
    pub states: i32,
    /// 创建时间
    pub created_at: String,
}

impl User {
    /// 从 DB 模型转换为响应实体
    /// override_nickname / override_avatar_url 不为空时覆盖对应字段
    pub fn from_model(m: &UserModel, override_nickname: &str, override_avatar_url: &str) -> Self {
        let nickname = if override_nickname.is_empty() {
            m.nickname.clone()
        } else {
            override_nickname.to_string()
        };
        let avatar_url = if override_avatar_url.is_empty() {
            m.avatar_url.clone()
        } else {
            override_avatar_url.to_string()
        };
        User {
            uid: m.id,
            type_id: m.type_id,
            account: m.account.clone(),
            mobile: m.mobile.clone(),
            email: m.email.clone(),
            nickname,
            avatar_url,
            sex: m.sex,
            states: m.states,
            created_at: format_created_at(&m.created_at),
        }
    }
}

/// 格式化创建时间 (业务统一使用 jiff, 对应 Go time.DateTime "2006-01-02 15:04:05")
///
/// chrono 仅用于读取 sea-orm 实体字段 (NaiveDateTime) 的数值, 格式化交给 jiff 完成。
fn format_created_at(ndt: &chrono::NaiveDateTime) -> String {
    jiff::civil::DateTime::new(
        ndt.year() as i16,
        ndt.month() as i8,
        ndt.day() as i8,
        ndt.hour() as i8,
        ndt.minute() as i8,
        ndt.second() as i8,
        0,
    )
    .map(|dt| jiff::fmt::strtime::format("%Y-%m-%d %H:%M:%S", dt).unwrap_or_default())
    .unwrap_or_default()
}

/// 请求-更新用户信息 (Uid 由 JWT 提供, 禁止 body 传入; 对齐 Go 的 json:"-")
///
/// nickname / avatar_url 缺省为空字符串 (对齐 Go: 未注册 validator, validate:"required" 不生效, 由业务校验兜底)
#[derive(Debug, Deserialize)]
pub struct UserUpdateRequest {
    /// 用户id (由 JWT 提供, 反序列化时忽略请求体中的 uid)
    #[serde(skip)]
    pub uid: u32,
    /// 昵称
    #[serde(default)]
    pub nickname: String,
    /// 头像地址
    #[serde(default)]
    pub avatar_url: String,
}

/// 请求-获取用户信息
#[derive(Debug, Deserialize)]
pub struct UserInfoRequest {
    /// 用户id
    #[serde(default)]
    pub uid: u32,
}

/// 请求-获取用户列表
#[derive(Debug, Deserialize)]
pub struct UserListRequest {
    /// 分页-当前页(默认 1)
    #[serde(default)]
    pub page: i32,
    /// 分页-每页数量(默认 12)
    #[serde(default)]
    pub page_size: i32,
}

/// 响应-获取用户列表
#[derive(Debug, Serialize)]
pub struct UserListResponse {
    /// 总数
    pub total: i32,
    /// 分页-当前页
    pub page: i32,
    /// 分页-每页数量
    pub page_size: i32,
    /// 列表
    pub list: Vec<User>,
}
