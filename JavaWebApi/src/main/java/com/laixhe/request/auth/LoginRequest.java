package com.laixhe.request.auth;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;
import lombok.ToString;

/**
 * 登录请求参数
 *
 * @author laixhe
 */
@Data
@ToString
public class LoginRequest {

    @NotBlank(message = "邮箱不能为空")
    @Email(message = "邮箱格式不正确")
    private String email;
    @NotBlank(message = "密码不能为空")
    @Size(min = 6, max = 20, message = "密码长度在6~20之间！")
    private String password;

}
