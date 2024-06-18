package com.laixhe.request.auth;

import jakarta.validation.constraints.*;
import lombok.Data;
import lombok.ToString;

/**
 * 注册请求参数
 *
 * @author laixhe
 */
@Data
@ToString
public class RegisterRequest {

    @NotBlank(message = "邮箱不能为空")
    @Email(message = "邮箱格式不正确")
    private String email;
    @NotBlank(message = "密码不能为空")
    @Size(min = 6, max = 20, message = "密码长度在6~20之间！")
    private String password;
    @NotBlank(message = "用户名不能为空")
    @Size(min = 2, max = 30, message = "用户名长度在2~30之间！")
    private String uname;
    @Min(value = 0, message = "年龄在0~200之间！")
    @Max(value = 200, message = "年龄在0~200之间！")
    private int age;

}
