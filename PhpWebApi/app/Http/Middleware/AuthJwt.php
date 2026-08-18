<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

use App\Result\ResultCode;
use App\Utils\JwtUtil;

/**
 * JWT 鉴权中间件 (与 Go 端 JWT 中间件对齐)
 *
 * 严格校验 Bearer 前缀 → 验签与有效期校验 → uid 必须为正数,
 * 任一环节失败统一返回 401 "Unauthorized" (不暴露具体失败原因, 防账号探测)。
 */
class AuthJwt
{
    public function handle(Request $request, Closure $next): Response
    {
        $authorization = (string)$request->header('Authorization');
        // 严格校验 Bearer 前缀, 避免非 Bearer 方案被误解析
        if (!str_starts_with($authorization, 'Bearer ')) {
            return response_error(ResultCode::AuthInvalid);
        }
        $token = substr($authorization, 7);
        try {
            $claims = JwtUtil::getInstance()->validatorToken($token);
            $uid = (int)$claims->get('uid');
            if ($uid <= 0) {
                return response_error(ResultCode::AuthInvalid);
            }

            $request->attributes->set('uid', $uid);
        } catch (\Throwable $e) {
            // catch Throwable 而非 Exception: 防御 JWT 库内部抛出的 TypeError 等 Error
            // 401 为鉴权失败原样返回; 其它异常 (如 JWT_SECRET 未配置) 收敛为固定 500 文案,
            // 避免把内部配置信息透传给客户端
            if ((int)$e->getCode() === ResultCode::AuthInvalid->value) {
                return response_exception(ResultCode::AuthInvalid->value, 'Unauthorized');
            }
            return response_exception(ResultCode::Service->value, 'internal server error');
        }
        return $next($request);
    }
}
