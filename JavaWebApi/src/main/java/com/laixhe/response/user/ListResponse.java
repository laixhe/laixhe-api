package com.laixhe.response.user;

import lombok.Data;
import lombok.ToString;

import java.util.List;

/**
 * @author laixhe
 */
@Data
@ToString
public class ListResponse {
    private List<UserResponse> list;
    private int total;
    private int page;
    private int size;
}
