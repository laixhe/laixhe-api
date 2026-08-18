package com.laixhe.webapi.security;

import com.laixhe.webapi.common.ApiException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

/**
 * 从安全上下文取出当前登录用户信息
 */
public final class ClaimsHolder {

    private ClaimsHolder() {
    }

    /** 当前登录用户 UID, 未登录时抛 401 */
    public static int uid() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth != null && auth.getPrincipal() instanceof JwtClaims claims && claims.uid() > 0) {
            return claims.uid();
        }
        throw ApiException.unauthorized();
    }
}
