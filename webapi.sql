
CREATE DATABASE `webapi` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci';

CREATE TABLE `user`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `password` varchar(120) NOT NULL DEFAULT '' COMMENT '用户密码',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `uname` varchar(100) NOT NULL DEFAULT '' COMMENT '用户名',
  `age` smallint UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户年龄',
  `score` decimal(10, 2) UNSIGNED NOT NULL DEFAULT 0.00 COMMENT '用户分数',
  `login_at` datetime NOT NULL COMMENT '登录时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_email`(`email` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;


INSERT INTO `user` (`password`,`email`,`uname`,`age`,`score`,`login_at`) VALUES ('$2y$10$bY4IoOWMVQkIg3ze.fUsOOCvkwr2oYsSpWQe.yDrJ4ZVPJODhtU8K', 'laixhe@laixhe.com', 'laixhe', 18, 100, NOW());

