package com.laixhe.webapi.dto;

import com.laixhe.webapi.entity.User;
import io.swagger.v3.oas.annotations.media.Schema;

import java.time.format.DateTimeFormatter;

/**
 * 用户信息 (对应 swagger entity.User, 不包含 password)
 */
@Schema(description = "用户信息")
public record UserResponse(
        @Schema(description = "用户id") int uid,
        @Schema(name = "type_id", description = "UserType: 1 - 普通用户") int typeId,
        @Schema(description = "账号") String account,
        @Schema(description = "手机号") String mobile,
        @Schema(description = "邮箱") String email,
        @Schema(description = "昵称") String nickname,
        @Schema(name = "avatar_url", description = "头像地址") String avatarUrl,
        @Schema(description = "UserSex: 0 - 未填写, 1 - 男, 2 - 女") int sex,
        @Schema(description = "UserState: 0 - 禁用, 1 - 正常") int states,
        @Schema(name = "created_at", description = "创建时间") String createdAt
) {

    private static final DateTimeFormatter TIME_FORMAT = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    public static UserResponse from(User user) {
        return from(user, user.getNickname(), user.getAvatarUrl());
    }

    /**
     * overrideNickname/overrideAvatarUrl 不为空时覆盖对应字段 (与 Go 版 NewUserFromModel 对齐)
     */
    public static UserResponse from(User user, String overrideNickname, String overrideAvatarUrl) {
        String nickname = overrideNickname != null && !overrideNickname.isEmpty() ? overrideNickname : user.getNickname();
        String avatar = overrideAvatarUrl != null && !overrideAvatarUrl.isEmpty() ? overrideAvatarUrl : user.getAvatarUrl();
        return new UserResponse(
                user.getId(),
                user.getTypeId(),
                user.getAccount(),
                user.getMobile(),
                user.getEmail(),
                nickname,
                avatar,
                user.getSex(),
                user.getStates(),
                user.getCreatedAt() != null ? user.getCreatedAt().format(TIME_FORMAT) : ""
        );
    }
}
