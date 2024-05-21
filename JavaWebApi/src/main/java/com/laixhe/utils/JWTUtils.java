package com.laixhe.utils;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTCreator;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.exceptions.AlgorithmMismatchException;
import com.auth0.jwt.exceptions.IncorrectClaimException;
import com.auth0.jwt.exceptions.SignatureVerificationException;
import com.auth0.jwt.exceptions.TokenExpiredException;
import com.auth0.jwt.interfaces.DecodedJWT;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;

import java.time.Instant;
import java.util.Map;

import com.laixhe.config.JwtConfig;
import com.laixhe.exception.BusinessException;
import com.laixhe.result.ResultCode;

/**
 * @author laixhe
 */
@Slf4j
@Data
@Component
public class JWTUtils {

    private static JwtConfig jwtConfig;

    @Autowired
    public JWTUtils(JwtConfig config) {
        jwtConfig = config;
    }

    // 常见的异常:
    // TokenExpiredException           令牌过期
    // SignatureVerificationException  签名不一致
    // AlgorithmMismatchException      算法不匹配
    // IncorrectClaimException         定义的有效期内不可用

    /**
     * 生成
     *
     * @param map payload 数据
     * @return String
     */
    public static String createToken(Map<String, Object> map) {
        // 创建 jwt 构造器
        JWTCreator.Builder builder = JWT.create();
        // 设置 payload
        map.forEach((key, value) -> {
            if (value instanceof String) {
                builder.withClaim(key, (String) value);
            } else if (value instanceof Integer) {
                builder.withClaim(key, (Integer) value);
            }
        });
        Instant instant = Instant.now();
        Instant instantExp = instant.plusSeconds(jwtConfig.getExpire());
        // 设置过期时间
        builder.withExpiresAt(instantExp)
                .withIssuedAt(instant)
                .withNotBefore(instant);
        // 生成 token
        return builder.sign(Algorithm.HMAC256(jwtConfig.getSecret()));
    }

    /**
     * 验证并返回 token 信息
     *
     * @param token 待校验的 token
     * @return DecodedJWT
     */
    public static DecodedJWT verifyToken(String token) {
        try {
            return JWT.require(Algorithm.HMAC256(jwtConfig.getSecret()))
                    .build().verify(token);
        } catch (SignatureVerificationException e) {
            // 签名无效
            throw new BusinessException(ResultCode.AuthInvalid, e.getMessage());
        } catch (TokenExpiredException e2) {
            // token 过期
            throw new BusinessException(ResultCode.AuthExpire, e2.getMessage());
        } catch (AlgorithmMismatchException | IncorrectClaimException e3) {
            // token 无效
            throw new BusinessException(ResultCode.AuthInvalid, e3.getMessage());
        } catch (Exception e4) {
            // 签名无效
            throw new BusinessException(ResultCode.AuthInvalid, e4.getMessage());
        }
    }

}
