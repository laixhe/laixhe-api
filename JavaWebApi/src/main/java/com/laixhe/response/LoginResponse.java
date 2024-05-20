package com.laixhe.response;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Data;
import lombok.ToString;

import java.time.LocalDateTime;

/**
 * @author Administrator
 */
@Data
@ToString
public class LoginResponse {
    private Integer uid;
    private String email;
    private String uname;
    private double score;
    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createdAt;
    private String token;
}
