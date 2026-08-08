<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;

use App\Http\Requests\LoginRequest;
use App\Http\Requests\RegisterRequest;
use App\Http\Services\AuthService;
use App\Http\Services\UserService;
use App\Result\ResultCode;
use App\Utils\JwtUtil;

/**
 * 鉴权相关 (与 Go 端 controllers/auth.go 对齐)
 */
class AuthController extends Controller
{
    /**
     * 注册
     *
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
        $authService = new AuthService();
        try {
            $user = $authService->register($registerRequest);
            $token = JwtUtil::getInstance()->createToken($user['id']);
            return response_success([
                'token' => $token,
                'user' => format_user($user),
            ]);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 登录
     *
     * 封禁账号与密码错误返回同一提示, 避免暴露账号状态 (可被探测)。
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function login(Request $request): JsonResponse
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

        $authService = new AuthService();
        $user = $authService->login($loginRequest);
        if (empty($user)) {
            return response_error(ResultCode::Param, '邮箱或密码错误');
        }
        // 判断密码是否匹配
        if (!password_verify($req['password'], $user['password'])) {
            return response_error(ResultCode::Param, '邮箱或密码错误');
        }

        $token = '';
        try {
            $token = JwtUtil::getInstance()->createToken($user['id']);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }

        return response_success([
            'token' => $token,
            'user' => format_user($user),
        ]);
    }

    /**
     * 刷新Jwt
     *
     * 用户不存在或账号被封禁时返回 401 (与 Go 端一致)。
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function refresh(Request $request): JsonResponse
    {
        // 获取登录用户ID (由 AuthJwt 中间件写入请求头)
        $uid = (int)$request->header('uid');
        if ($uid <= 0) {
            return response_error(ResultCode::AuthInvalid);
        }

        $userService = new UserService();
        $user = $userService->info($uid);
        // 只需要 states 与响应字段, 排除 password 减少不必要的列传输
        if (empty($user) || (int)($user['states'] ?? 0) !== 1) {
            return response_error(ResultCode::AuthInvalid);
        }

        $token = '';
        try {
            $token = JwtUtil::getInstance()->createToken($uid);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }

        return response_success([
            'token' => $token,
            'user' => format_user($user),
        ]);
    }
}
