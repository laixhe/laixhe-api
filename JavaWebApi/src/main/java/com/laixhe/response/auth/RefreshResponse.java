package com.laixhe.response.auth;

import lombok.Data;
import lombok.ToString;

/**
 * 登录响应参数
 *
 * @author laixhe
 */
@Data
@ToString
public class RefreshResponse {
    private String token;
}
