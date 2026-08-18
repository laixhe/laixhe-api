//! 用户账号状态 (对应 Go 项目 app/models/user_states.go)
//!
//! 以下函数当前业务未使用, 保留以与 Go 版代码结构对齐 (Go 端同样提供这些辅助函数);
//! 由 `#![allow(dead_code)]` 抑制编译告警。
#![allow(dead_code)]

/// 禁用
pub const USER_STATE_BANNED: i32 = 0;
/// 正常
pub const USER_STATE_NORMAL: i32 = 1;

/// 判断用户状态值是否有效
pub fn is_user_state_valid(s: i32) -> bool {
    matches!(s, USER_STATE_BANNED | USER_STATE_NORMAL)
}

/// 获取用户状态中文描述
pub fn get_user_state_text(s: i32) -> &'static str {
    match s {
        USER_STATE_BANNED => "禁用",
        USER_STATE_NORMAL => "正常",
        _ => "",
    }
}
