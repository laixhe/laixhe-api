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

/**
 * 直接返回一个 Result 对象作为 JSON (一般用于请求校验失败的结果透传)。
 * HTTP 状态码与 Result.code 同步 (0=成功转 200, 其余如 422 直接作为状态码)。
 */
function response_result(Result $result): JsonResponse
{
    return response()->json($result, $result->code->value === 0 ? 200 : $result->code->value);
}

/**
 * 成功响应: 直接返回业务数据 JSON (如 {"token":..., "user":...})
 */
function response_success($data = []): JsonResponse
{
    return response()->json($data);
}

/**
 * 已知业务错误响应: 返回 {"code": <HTTP状态码>, "message": ...}, HTTP 状态码同步
 */
function response_error(ResultCode $code, string $msg = ''): JsonResponse
{
    return response()->json(new Result($code, $msg ?: $code->text()), $code->value);
}

/**
 * 异常转错误响应: 按异常 code 映射到 ResultCode 后返回, 未知 code 归为 500
 */
function response_exception(int $code, string $msg = ''): JsonResponse
{
    $result = ResultCode::intToEnum($code);
    return response()->json(new Result($result, $msg ?: $result->text()), $result->value);
}

/**
 * 用户信息响应结构 (与 Go 端 entity.User 对齐)
 * @param array $user 用户记录 (含 id/type_id/account/mobile/email/nickname/avatar_url/sex/states/created_at)
 * @return array
 */
function format_user(array $user): array
{
    return [
        'uid' => (int)($user['id'] ?? 0),
        'type_id' => (int)($user['type_id'] ?? 0),
        'account' => (string)($user['account'] ?? ''),
        'mobile' => (string)($user['mobile'] ?? ''),
        'email' => (string)($user['email'] ?? ''),
        'nickname' => (string)($user['nickname'] ?? ''),
        'avatar_url' => (string)($user['avatar_url'] ?? ''),
        'sex' => (int)($user['sex'] ?? 0),
        'states' => (int)($user['states'] ?? 0),
        'created_at' => (string)($user['created_at'] ?? ''),
    ];
}
