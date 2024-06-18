package com.laixhe.request.user;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.laixhe.validator.DateTimeStr;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;
import lombok.ToString;

/**
 * 修改用户信息请求参数
 *
 * @author laixhe
 */
@Data
@ToString
public class UpdateRequest {
    @NotBlank(message = "用户名不能为空！")
    @Size(min = 2, max = 30, message = "用户名长度在2~30之间！")
    private String uname;

    @JsonProperty("login_at")  // 前端字段名为 login_at (前端字段名与后端不对应)
    @NotBlank(message = "登录时间不能为空！")
    @DateTimeStr(format = "yyyy-MM-dd HH:mm:ss", message = "登录时间格式不对！")
    private String loginAt;
}
