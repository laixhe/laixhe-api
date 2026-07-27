<?php

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Http\JsonResponse;

/**
 * 抛出异常并结束程序
 *
 * @param bool $condition 判断条件，判断结果为 true 时生效，false时继续业务流程
 * @param string $message
 * @param int $code
 * @return void
 * @throws Throwable
 */
function throw_if_fail(bool $condition, string $message, int $code = 0): void
{
    throw_if($condition, 'RuntimeException', $message, $code);
}

function response_result(Result $result): JsonResponse
{
    return response()->json($result);
}

function response_success($data = []): JsonResponse
{
    return response()->json($data);
}

function response_error(ResultCode $code, string $msg = ''): JsonResponse
{
    return response()->json(new Result($code, $msg ?: $code->text()));
}

function response_exception(int $code, string $msg = ''): JsonResponse
{
    $result = ResultCode::intToEnum($code);
    return response()->json(new Result($result, $msg ?: $result->text()));
}
