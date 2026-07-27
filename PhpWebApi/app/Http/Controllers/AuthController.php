<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;

use App\Http\Request\LoginRequest;
use App\Http\Request\RegisterRequest;
use App\Http\Services\AuthService;
use App\Http\Services\UserService;
use App\Result\ResultCode;
use App\Utils\JwtUtil;

/**
 * 鉴权相关
 */
class AuthController extends Controller
{
    /**
     * 注册
     * @param Request $request
     * @return JsonResponse
     */
    public function register(Request $request): JsonResponse
    {
        // 获取想要的请求参数
        $req = $request->only([
            'nickname',
            'email',
            'password',
        ]);
        $registerRequest = new RegisterRequest();
        $error = $registerRequest->validator($req);
        if ($error !== null) {
            return response_result($error);
        }
        //
        $loginService = new AuthService();
        try {
            $user = $loginService->register($registerRequest);
            $token = JwtUtil::getInstance()->createToken($user['id']);
            return response_success([
                'token' => $token,
                'user' => [
                    'uid' => $user['id'],
                    'type_id' => $user['type_id'],
                    'nickname' => $user['nickname'],
                    'avatar_url' => $user['avatar_url'],
                    'states' => $user['states'],
                    'created_at' => $user['created_at'],
                ],
            ]);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 登录
     * @param Request $request
     * @return JsonResponse
     */
    public function Login(Request $request): JsonResponse
    {
        // 获取想要的请求参数
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
            return response_error(ResultCode::Param, '邮箱或密码不正确');
        }
        // 判断密码是否匹配
        if (!password_verify($req['password'], $user['password'])) {
            return response_error(ResultCode::Param, '邮箱或密码不正确');
        }

        $token = '';
        try {
            $token = JwtUtil::getInstance()->createToken($user['id']);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }

        return response_success([
            'token' => $token,
            'user' => [
                'uid' => $user['id'],
                'type_id' => $user['type_id'],
                'nickname' => $user['nickname'],
                'avatar_url' => $user['avatar_url'],
                'states' => $user['states'],
                'created_at' => $user['created_at'],
            ],
        ]);
    }

    /**
     * 刷新Jwt
     * @param Request $request
     * @return JsonResponse
     */
    public function refresh(Request $request): JsonResponse
    {
        // 获取登录用户ID
        $uid = (int)$request->header('uid');

        $token = '';
        try {
            $token = JwtUtil::getInstance()->createToken($uid);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
        //
        $userService = new UserService();
        $user = $userService->info($uid);
        if (empty($user)) {
            return response_error(ResultCode::Param, '用户不存在');
        }
        return response_success([
            'token' => $token,
            'user' => [
                'uid' => $user['id'],
                'type_id' => $user['type_id'],
                'nickname' => $user['nickname'],
                'avatar_url' => $user['avatar_url'],
                'states' => $user['states'],
                'created_at' => $user['created_at'],
            ],
        ]);
    }
}
