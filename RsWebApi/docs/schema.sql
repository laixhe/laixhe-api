-- WebApi 数据库初始化脚本 (MySQL)
-- 与 sea-orm 实体模型 (src/app/models) 保持一致

CREATE TABLE IF NOT EXISTS `user` (
    `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `type_id` INT NOT NULL DEFAULT 0 COMMENT '类型 1普通',
    `account` VARCHAR(120) NOT NULL DEFAULT '' COMMENT '账号',
    `mobile` VARCHAR(120) NOT NULL DEFAULT '' COMMENT '手机号',
    `email` VARCHAR(120) NOT NULL DEFAULT '' COMMENT '邮箱',
    `password` VARCHAR(120) NOT NULL DEFAULT '' COMMENT '密码',
    `nickname` VARCHAR(120) NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar_url` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像地址',
    `sex` INT NOT NULL DEFAULT 0 COMMENT '性别 0未填写 1男 2女',
    `states` INT NOT NULL DEFAULT 0 COMMENT '状态 0封禁 1正常',
    `created_at` DATETIME NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_account` (`account`),
    KEY `idx_user_mobile` (`mobile`),
    -- email 唯一索引 (与 webapi.sql 一致; 注册先查后插 + 数据库唯一约束双重防重)
    UNIQUE KEY `idx_user_email` (`email`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户表';

CREATE TABLE IF NOT EXISTS `user_extend` (
    `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `uid` INT NOT NULL COMMENT '用户UID',
    `birthday` INT NOT NULL DEFAULT 0 COMMENT '生日(年月日)',
    `height` INT NOT NULL DEFAULT 0 COMMENT '身高(cm)',
    `weight` INT NOT NULL DEFAULT 0 COMMENT '体重(kg)',
    PRIMARY KEY (`id`),
    -- 一对一约束: 与 user_third_party.uid 一致 (对齐 webapi.sql / Go / TS 端)
    UNIQUE KEY `idx_user_extend_uid` (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户扩展表';

CREATE TABLE IF NOT EXISTS `user_third_party` (
    `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `uid` INT NOT NULL COMMENT '用户UID',
    `wechat_unionid` VARCHAR(200) NOT NULL DEFAULT '' COMMENT '微信unionid',
    `wechat_openid` VARCHAR(200) NOT NULL DEFAULT '' COMMENT '微信openid',
    PRIMARY KEY (`id`),
    -- 一对一约束: 一个用户只保留一条第三方绑定记录 (对齐 webapi.sql / Go / TS 端)
    UNIQUE KEY `idx_user_third_party_uid` (`uid`),
    KEY `idx_user_third_party_wechat_openid` (`wechat_openid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户第三方表';

CREATE TABLE IF NOT EXISTS `config_common` (
    `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `key` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '配置键',
    `value` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '配置值',
    `describe` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
    PRIMARY KEY (`id`),
    KEY `config_common_key_idx` (`key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '通用配置表';
