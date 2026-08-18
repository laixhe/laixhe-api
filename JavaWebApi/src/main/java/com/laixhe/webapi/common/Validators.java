package com.laixhe.webapi.common;

/**
 * 手写参数校验 (与 Go 版 controllers/controller.go 对齐, 错误文案完全一致)
 */
public final class Validators {

    private Validators() {
    }

    /**
     * 昵称长度校验 (注册与更新用户信息共用)。
     * 按 Unicode 码点统计, 与 Go RuneCount 行为一致 (中文等按 1 字符计)。
     */
    public static void validateNickname(String nickname) {
        int len = nickname == null ? 0 : nickname.codePointCount(0, nickname.length());
        if (len < 2) {
            throw ApiException.paramError("昵称长度不能小于2位");
        }
        if (len > 20) {
            throw ApiException.paramError("昵称长度不能超过20位");
        }
    }

    /**
     * 头像地址校验: 长度不超过 255 位; 非空时必须精确以 http:// 或 https:// 开头
     */
    public static void validateAvatarUrl(String avatarUrl) {
        if (avatarUrl != null && avatarUrl.length() > 255) {
            throw ApiException.paramError("头像地址长度不能超过255位");
        }
        if (avatarUrl != null && !avatarUrl.isEmpty()) {
            if (!avatarUrl.startsWith("http://") && !avatarUrl.startsWith("https://")) {
                throw ApiException.paramError("头像地址必须以http或https开头");
            }
        }
    }
}
