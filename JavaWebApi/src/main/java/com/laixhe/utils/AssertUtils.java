package com.laixhe.utils;

import com.laixhe.exception.BusinessException;
import com.laixhe.result.ResultCode;

/**
 * @author laixhe
 */
public class AssertUtils {

    public static void isTrue(boolean expression, ResultCode code) {
        if (!expression) {
            throw new BusinessException(code, code.getMsg());
        }
    }

    public static void isTrue(boolean expression, ResultCode code, String msg) {
        if (!expression) {
            throw new BusinessException(code, msg);
        }
    }

}
