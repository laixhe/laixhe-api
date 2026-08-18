<?php

namespace App\Http\Controllers;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Http\Request;

abstract class Controller
{
    /**
     * 校验顶层 body 必须是对象 (数组/标量/null 按请求格式错误返回 400, 与 Go/Rust 端绑定层行为一致)。
     * 空 body (无内容) 跳过, 由各 Request validator 的 required 兜底为 422 (与 TS 端缺 body 语义一致)。
     *
     * @return Result|null 校验失败返回 400 Result, 通过返回 null
     */
    protected function validateTopLevelBody(Request $request): ?Result
    {
        $rawBody = $request->getContent();
        if ($rawBody !== '' && !is_object(json_decode($rawBody))) {
            return new Result(ResultCode::BadRequest, 'Bad Request');
        }
        return null;
    }
}
