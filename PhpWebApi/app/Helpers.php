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
 * 校验失败 → 错误 Result: 字段类型非字符串 (如数字) 视为请求格式错误返回 400
 * (与 Go/Rust/TS 端绑定层行为一致), 其余业务校验失败为 422。
 * 供各 Request 类复用, 避免复制粘贴同一段依赖 Validator::failed() 内部结构的逻辑。
 */
function validation_error(\Illuminate\Validation\Validator $validator): Result
{
    $failed = $validator->failed();
    foreach ($failed as $rules) {
        if (isset($rules['String'])) {
            return new Result(ResultCode::BadRequest, 'Bad Request');
        }
    }
    return new Result(ResultCode::Param, $validator->errors()->first());
}

/**
 * 异常安全转错误响应 (替代直接 response_exception($e->getCode(), $e->getMessage())):
 * 已知业务错误码 (400/401/422/429) 的消息可透传客户端;
 * 其余异常 (如 QueryException 的完整 SQL 与绑定值) 只记服务端日志并返回默认 500 文案,
 * 避免把 SQL 等内部细节泄露给调用方 (与 bootstrap/app.php 未捕获异常固定文案的防护对齐)。
 *
 * @param \Throwable $e       被捕获的异常
 * @param string $context 日志上下文 (建议传"控制器/方法", 便于定位)
 */
function response_exception_safe(\Throwable $e, string $context = ''): JsonResponse
{
    $result = ResultCode::intToEnum((int)$e->getCode());
    if ($result !== ResultCode::Service) {
        return response_exception($e->getCode(), $e->getMessage());
    }
    logger()->error(($context !== '' ? "[$context] " : '') . '未处理异常', ['exception' => $e]);
    return response_exception(ResultCode::Service->value);
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
