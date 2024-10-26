package com.laixhe.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.Data;
import lombok.ToString;

import java.time.LocalDateTime;

/**
 * 用户表模型
 *
 * @author laixhe
 */
@Data
@ToString
@Table("user")
public class User {
    // 用户ID
    @Id(keyType = KeyType.Auto)
    private Integer id;
    // 用户密码
    private String password;
    // 用户邮箱
    private String email;
    // 用户名
    private String uname;
    // 用户年龄
    private Integer age;
    // 用户分数
    private Double score;
    // 登录时间
    private LocalDateTime loginAt;
    // 创建时间
    @Column(onInsertValue = "now()")
    private LocalDateTime createdAt;
    // 更新时间
    @Column(onUpdateValue = "now()", onInsertValue = "now()")
    private LocalDateTime updatedAt;
    // 删除时间
    private LocalDateTime deletedAt;
}
