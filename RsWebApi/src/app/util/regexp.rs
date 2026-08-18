//! 正则校验 (对应 Go 项目 app/util/regexp.go + gonet/utils.IsEmail)

use std::sync::OnceLock;

use regex::Regex;

/// 密码正则：字母、数字、_、@、$，长度 6~64 位
/// 长度上限 64: bcrypt 只取前 72 字节, 上限保证任何合法密码都不会被静默截断产生碰撞
fn matching_password() -> &'static Regex {
    static RE: OnceLock<Regex> = OnceLock::new();
    RE.get_or_init(|| Regex::new(r"^[a-zA-Z0-9_@$]{6,64}$").unwrap())
}

/// 密码是否合法
pub fn is_password(password: &str) -> bool {
    matching_password().is_match(password)
}

/// 邮箱正则
fn matching_email() -> &'static Regex {
    static RE: OnceLock<Regex> = OnceLock::new();
    RE.get_or_init(|| Regex::new(r"^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$").unwrap())
}

/// 邮箱是否合法 (对应 gonet/utils.IsEmail)
pub fn is_email(email: &str) -> bool {
    matching_email().is_match(email)
}
