package com.laixhe.response.user;

import lombok.Data;
import lombok.ToString;

/**
 * 查询响应信息响应参数
 *
 * @author laixhe
 */
@Data
@ToString
public class InfoResponse {
    private UserResponse user;
}
