//! 公共参数校验 (提取自 controllers 中重复的校验逻辑, 对应 Go 控制器内联校验)

use crate::error::ApiError;

/// 校验昵称格式: 长度 2~20 位
///
/// 错误消息与 Go 控制器保持一致。
pub fn validate_nickname(nickname: &str) -> Result<(), ApiError> {
    if nickname.len() < 2 {
        return Err(ApiError::param_error("昵称长度不能小于2位"));
    }
    if nickname.len() > 20 {
        return Err(ApiError::param_error("昵称长度不能超过20位"));
    }
    Ok(())
}

/// 校验头像地址格式: 长度不超过 255, 且以 http/https 开头 (空字符串允许)
pub fn validate_avatar_url(avatar_url: &str) -> Result<(), ApiError> {
    if avatar_url.len() > 255 {
        return Err(ApiError::param_error("头像地址长度不能超过255位"));
    }
    if !avatar_url.is_empty() && !avatar_url.starts_with("http") {
        return Err(ApiError::param_error("头像地址必须以http或https开头"));
    }
    Ok(())
}
