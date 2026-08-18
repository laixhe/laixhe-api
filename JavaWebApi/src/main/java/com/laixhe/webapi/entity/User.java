package com.laixhe.webapi.entity;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.Table;
import jakarta.persistence.UniqueConstraint;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;

/**
 * 用户 (对应 Go models/user.go, 表结构见 webapi.sql)
 */
@Entity
@Table(name = "user", uniqueConstraints = {
        @UniqueConstraint(name = "user_account_idx", columnNames = "account"),
        @UniqueConstraint(name = "user_email_idx", columnNames = "email")
}, indexes = @Index(name = "user_mobile_idx", columnList = "mobile"))
@Getter
@Setter
@NoArgsConstructor
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    @Column(name = "type_id", nullable = false)
    private Integer typeId = 0;

    @Column(name = "account", nullable = false, length = 120)
    private String account = "";

    @Column(name = "mobile", nullable = false, length = 120)
    private String mobile = "";

    @Column(name = "email", nullable = false, length = 120)
    private String email = "";

    /** 密码哈希, 不参与任何对外响应 */
    @Column(name = "password", nullable = false, length = 120)
    private String password = "";

    @Column(name = "nickname", nullable = false, length = 120)
    private String nickname = "";

    @Column(name = "avatar_url", nullable = false, length = 255)
    private String avatarUrl = "";

    @Column(name = "sex", nullable = false)
    private Integer sex = 0;

    @Column(name = "states", nullable = false)
    private Integer states = 0;

    @Column(name = "created_at", nullable = false)
    private LocalDateTime createdAt;

    @Column(name = "updated_at", nullable = false)
    private LocalDateTime updatedAt;
}
