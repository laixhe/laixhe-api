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
 * 用户扩展 (对应 Go models/user_extend.go)
 */
@Entity
@Table(name = "user_extend", uniqueConstraints = @UniqueConstraint(name = "user_extend_uid_idx", columnNames = "uid"))
@Getter
@Setter
@NoArgsConstructor
public class UserExtend {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    @Column(name = "uid", nullable = false)
    private Integer uid;

    /** 生日(年月日) */
    @Column(name = "birthday", nullable = false)
    private Integer birthday = 0;

    /** 身高(cm) */
    @Column(name = "height", nullable = false)
    private Integer height = 0;

    /** 体重(kg) */
    @Column(name = "weight", nullable = false)
    private Integer weight = 0;
}
