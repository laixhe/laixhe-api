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
        $authorization = $request->header('Authorization');
        if (strlen($authorization) <= 7) {
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
