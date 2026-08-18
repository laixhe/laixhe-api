//! 公共参数校验 (提取自 controllers 中重复的校验逻辑, 对应 Go 控制器内联校验)

use crate::error::ApiError;

/// 校验昵称格式: 长度 2~20 位 (按"字符"计数, 对齐 Go RuneCountInString;
/// 使用 chars().count() 而非 len(), 否则中文等多字节字符会被误判)
///
/// 错误消息与 Go 控制器保持一致。
pub fn validate_nickname(nickname: &str) -> Result<(), ApiError> {
    let chars = nickname.chars().count();
    if chars < 2 {
        return Err(ApiError::param_error("昵称长度不能小于2位"));
    }
    if chars > 20 {
        return Err(ApiError::param_error("昵称长度不能超过20位"));
    }
    Ok(())
}

/// 校验头像地址格式: 长度不超过 255, 且须精确以 http:// 或 https:// 开头 (空字符串允许)
///
/// 不用 starts_with("http"), 否则 httpxxx:// 之类的畸形 scheme 也能通过 (与 Go/PHP/TS 端保持一致)
pub fn validate_avatar_url(avatar_url: &str) -> Result<(), ApiError> {
    if avatar_url.len() > 255 {
        return Err(ApiError::param_error("头像地址长度不能超过255位"));
    }
    if !avatar_url.is_empty()
        && !avatar_url.starts_with("http://")
        && !avatar_url.starts_with("https://")
    {
        return Err(ApiError::param_error("头像地址必须以http或https开头"));
    }
    Ok(())
}
