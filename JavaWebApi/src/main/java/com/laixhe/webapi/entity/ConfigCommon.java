package com.laixhe.webapi.entity;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.Table;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

/**
 * 通用配置 (对应 Go models/config_common.go, 表结构见 webapi.sql)
 */
@Entity
@Table(name = "config_common", indexes = @Index(name = "config_common_key_idx", columnList = "`key`"))
@Getter
@Setter
@NoArgsConstructor
public class ConfigCommon {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    /** 配置键 */
    @Column(name = "`key`", nullable = false, length = 255)
    private String key = "";

    /** 配置值 */
    @Column(name = "`value`", nullable = false, length = 512)
    private String value = "";

    /** 描述 */
    @Column(name = "`describe`", nullable = false, length = 255)
    private String describe = "";
}
