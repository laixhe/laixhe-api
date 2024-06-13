package com.laixhe.request.user;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;
import lombok.ToString;

/**
 * @author laixhe
 */
@Data
@ToString
public class UpdateRequest {
    @NotBlank(message="用户名不能为空")
    @Size(min=2, max=30, message="用户名长度在2~30之间！")
    private String uname;
}
