package com.laixhe.webapi.entity;

/**
 * 用户性别 (对应 Go models/user_sex.go)
 */
public final class UserSex {

    /** 未填写 */
    public static final int UNKNOWN = 0;
    /** 男 */
    public static final int MALE = 1;
    /** 女 */
    public static final int FEMALE = 2;

    private UserSex() {
    }
}
