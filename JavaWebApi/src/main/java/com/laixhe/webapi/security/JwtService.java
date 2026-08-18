package com.laixhe.webapi.security;

import com.laixhe.webapi.config.AppProperties;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.Date;

/**
 * JWT 生成与解析 (HS256, 与 Go 版 gonet/jwt 对齐)
 */
@Service
public class JwtService {

    private final SecretKey key;
    private final long expireSeconds;

    public JwtService(AppProperties props) {
        this.key = Keys.hmacShaKeyFor(props.getJwt().getSecretKey().getBytes(StandardCharsets.UTF_8));
        this.expireSeconds = props.getJwt().getExpireSeconds();
    }

    /**
     * 生成 JWT: 载荷含 uid, 并设置签发/生效/过期时间 (与 Go 版 NewJwtClaims 对齐)
     */
    public String generateToken(int uid) {
        Instant now = Instant.now();
        return Jwts.builder()
                .claim("uid", uid)
                .issuedAt(Date.from(now))        // 签发时间
                .notBefore(Date.from(now))       // 生效时间
                .expiration(Date.from(now.plusSeconds(expireSeconds))) // 过期时间
                .signWith(key, Jwts.SIG.HS256)
                .compact();
    }

    /**
     * 解析并校验 JWT (签名/过期校验失败抛 JwtException), 并校验 uid 非零
     */
    public JwtClaims parse(String token) {
        Claims payload = Jwts.parser().verifyWith(key).build().parseSignedClaims(token).getPayload();
        Integer uid = payload.get("uid", Integer.class);
        if (uid == null) {
            throw new io.jsonwebtoken.JwtException("missing uid claim");
        }
        return new JwtClaims(uid);
    }
}
