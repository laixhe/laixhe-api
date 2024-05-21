package com.laixhe.response.auth;

import lombok.Data;
import lombok.ToString;

/**
 * @author laixhe
 */
@Data
@ToString
public class RefreshResponse {
    private String token;
}
