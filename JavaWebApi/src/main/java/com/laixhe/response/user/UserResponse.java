package com.laixhe.response.user;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Data;
import lombok.ToString;

import java.time.LocalDateTime;

/**
 * @author laixhe
 */
@Data
@ToString
public class UserResponse {
    private Integer uid;
    private String email;
    private String uname;
    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createdAt;
}
