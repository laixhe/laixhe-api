package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;

import java.util.List;

/**
 * 用户列表响应 (对应 swagger entity.UserListResponse)
 */
@Schema(description = "用户列表响应")
public record UserListResponse(
        @Schema(description = "总数") int total,
        @Schema(description = "分页-当前页") int page,
        @Schema(name = "page_size", description = "分页-每页数量") int pageSize,
        @Schema(description = "列表") List<UserResponse> list
) {
}
