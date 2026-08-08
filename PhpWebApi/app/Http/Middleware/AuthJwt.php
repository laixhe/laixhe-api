<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

use App\Result\ResultCode;
use App\Utils\JwtUtil;

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
            if($uid <= 0) {
                return response_error(ResultCode::AuthInvalid);
            }

            $request->headers->set('uid', $uid);
        } catch (\Exception $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
        return $next($request);
    }
}
