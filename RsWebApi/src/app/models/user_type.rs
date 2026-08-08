//! 用户类型 (对应 Go 项目 app/models/user_type.go)
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
