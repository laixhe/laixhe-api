package com.laixhe.result;

/**
 * @author laixhe
 */
public interface IResult {
    /**
     * 获取响应码
     * @return int
     */
    int getCode();

    /**
     * 获取响应错误信息
     * @return String
     */
    String getMsg();
}
