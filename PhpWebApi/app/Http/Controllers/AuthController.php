<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

use App\Http\Request\LoginRequest;
use App\Http\Request\RegisterRequest;
use App\Http\Services\AuthService;
use App\Result\ResultCode;
use App\Utils\JwtUtil;

class AuthController extends Controller
{
    // 注册
    public function register(Request $request)
    {
        $req = $request->only([
            'email',
            'password',
            'uname',
            'age',
        ]);
        $registerRequest = new RegisterRequest();
        $error = $registerRequest->validator($req);
        if ($error !== null) {
            return response_result($error);
        }

        $loginService = new AuthService();
        try {
            $user = $loginService->register($registerRequest);
            $token = JwtUtil::getInstance()->createToken(['uid' => $user['id']]);
            return response_success([
                'token' => $token,
                'user' => [
                    'uid' => $user['id'],
                    'uname' => $user['uname'],
                    'email' => $user['email'],
                    'created_at' => $user['created_at'],
                ],
            ]);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
    }

    // 登录
    public function Login(Request $request)
    {
        $req = $request->only([
            'email',
            'password',
        ]);
        $loginRequest = new LoginRequest();
        $error = $loginRequest->validator($req);
        if ($error !== null) {
            return response_result($error);
        }

        $loginService = new AuthService();
        $user = $loginService->login($loginRequest);
        if (empty($user)) {
            return response_error(ResultCode::PasswordError, '');
        }

        $token = '';
        try {
            $token = JwtUtil::getInstance()->createToken(['uid' => $user['id']]);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }

        return response_success([
            'token' => $token,
            'user' => [
                'uid' => $user['id'],
                'uname' => $user['uname'],
                'email' => $user['email'],
                'created_at' => $user['created_at'],
            ],
        ]);
    }
}
