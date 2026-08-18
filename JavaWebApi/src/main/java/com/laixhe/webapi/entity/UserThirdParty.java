package com.laixhe.webapi.entity;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.UniqueConstraint;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

/**
 * 用户第三方 (对应 Go models/user_third_party.go)
 */
@Entity
@Table(name = "user_third_party", uniqueConstraints = @UniqueConstraint(name = "user_third_party_uid_idx", columnNames = "uid"))
@Getter
@Setter
@NoArgsConstructor
public class UserThirdParty {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    @Column(name = "uid", nullable = false)
    private Integer uid;

    /** 微信unionid */
    @Column(name = "wechat_unionid", nullable = false, length = 200)
    private String wechatUnionid = "";

    /** 微信openid */
    @Column(name = "wechat_openid", nullable = false, length = 200)
    private String wechatOpenid = "";
}
