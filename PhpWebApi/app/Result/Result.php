<?php

namespace App\Result;

/**
 * 响应请求
 */
class Result
{
    public ResultCode $code; // 响应码
    public string $message; // 响应错误信息

    public function __construct(ResultCode $code, string $msg)
    {
        $this->code = $code;
        $this->message = $msg;
    }
}
