//! 用户性别 (对应 Go 项目 app/models/user_sex.go)
//!
//! 以下函数当前业务未使用, 保留以与 Go 版代码结构对齐 (Go 端同样提供这些辅助函数);
//! 由 `#![allow(dead_code)]` 抑制编译告警。
#![allow(dead_code)]

/// 未填写
pub const USER_SEX_UNKNOWN: i32 = 0;
/// 男
pub const USER_SEX_MALE: i32 = 1;
/// 女
pub const USER_SEX_FEMALE: i32 = 2;

/// 判断性别值是否有效（仅男/女为有效值）
pub fn is_user_sex_valid(s: i32) -> bool {
    matches!(s, USER_SEX_MALE | USER_SEX_FEMALE)
}

/// 获取性别中文描述
pub fn get_user_sex_text(s: i32) -> &'static str {
    match s {
        USER_SEX_MALE => "男",
        USER_SEX_FEMALE => "女",
        _ => "",
    }
}
