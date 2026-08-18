package com.laixhe.webapi.entity;

/**
 * 用户状态 (对应 Go models/user_states.go)
 */
public final class UserState {

    /** 禁用 */
    public static final int BANNED = 0;
    /** 正常 */
    public static final int NORMAL = 1;

    private UserState() {
    }
}
