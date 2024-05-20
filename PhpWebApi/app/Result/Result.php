<?php

namespace App\Result;

/**
 * 响应请求
 */
class Result
{
    public ResultCode $code; // 响应码
    public string $msg; // 响应信息
    public array $data; // 数据

    public function __construct(ResultCode $code, string $msg, array $data)
    {
        $this->code = $code;
        $this->msg = $msg;
        $this->data = $data;
    }
}
