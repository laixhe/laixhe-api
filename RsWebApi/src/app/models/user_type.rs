//! 用户类型 (对应 Go 项目 app/models/user_type.go)
//!
//! 以下函数当前业务未使用, 保留以与 Go 版代码结构对齐 (Go 端同样提供这些辅助函数);
//! 由 `#![allow(dead_code)]` 抑制编译告警。
#![allow(dead_code)]

/// 普通用户
pub const USER_TYPE_ORDINARY: i32 = 1;

/// 判断用户类型是否有效
pub fn is_user_type_valid(t: i32) -> bool {
    matches!(t, USER_TYPE_ORDINARY)
}

/// 获取用户类型中文描述
pub fn get_user_type_text(t: i32) -> &'static str {
    match t {
        USER_TYPE_ORDINARY => "普通用户",
        _ => "",
    }
}
