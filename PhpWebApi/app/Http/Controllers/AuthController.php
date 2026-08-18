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
use OpenApi\Attributes as OA;

/**
 * 鉴权相关 (与 Go 端 controllers/auth.go 对齐)
 */
class AuthController extends Controller
{
    /**
     * 注册
     */
    #[OA\Post(
        path: '/api/v1/auth/register',
        summary: '注册',
        tags: ['Auth'],
        requestBody: new OA\RequestBody(
            required: true,
            content: new OA\JsonContent(ref: '#/components/schemas/AuthRegisterRequest')
        ),
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/AuthTokenResponse')),
            new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 422, description: '参数错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function register(Request $request): JsonResponse
    {
        // 顶层 body 类型校验 (数组/标量/null → 400; 空 body → 422 兜底): 见 Controller::validateTopLevelBody
        $topLevelError = $this->validateTopLevelBody($request);
        if ($topLevelError !== null) {
            return response_result($topLevelError);
        }

        // 获取想要的请求参数
        $req = $request->only([
            'nickname',
            'email',
            'password',
        ]);
        // 教学取舍说明: 这里用 new + 手动 validator() 替代 Laravel FormRequest 自动校验,
        // Service 也用 new 实例化不走容器 DI — 为与 Go/Rust/TS 端的绑定层/服务层语义逐行对齐,
        // 避免 Laravel 的隐式依赖注入魔法干扰跨语言对照; 生产项目建议回归 FormRequest + 构造函数 DI。
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
            // 已知业务错误 (如"邮箱已存在") 透传; 未知异常 (含 SQL) 记日志后返回固定文案 (见 Helpers)
            return response_exception_safe($e, 'auth/register');
        }
    }

    /**
     * 登录
     *
     * 封禁账号与密码错误返回同一提示, 避免暴露账号状态 (可被探测)。
     */
    #[OA\Post(
        path: '/api/v1/auth/login',
        summary: '登录',
        tags: ['Auth'],
        requestBody: new OA\RequestBody(
            required: true,
            content: new OA\JsonContent(ref: '#/components/schemas/AuthLoginRequest')
        ),
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/AuthTokenResponse')),
            new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 422, description: '参数错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function login(Request $request): JsonResponse
    {
        // 顶层 body 类型校验 (数组/标量/null → 400; 空 body → 422 兜底): 见 Controller::validateTopLevelBody
        $topLevelError = $this->validateTopLevelBody($request);
        if ($topLevelError !== null) {
            return response_result($topLevelError);
        }

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
            // 登录失败统一返回 422 (参数错误) 而非 401: 与 Go/PHP/Rust/TS 四端对齐,
            // 不区分"邮箱不存在/密码错误/账号封禁", 避免暴露账号状态 (防账号探测);
            // 401 语义只留给"未携带/无效 JWT"的鉴权场景 (见 AuthJwt 中间件)
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
            return response_exception_safe($e, 'auth/login');
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
     */
    #[OA\Post(
        path: '/api/v1/auth/refresh',
        summary: '刷新Jwt',
        tags: ['Auth'],
        security: [['BearerAuth' => []]],
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/AuthTokenResponse')),
            new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 401, description: '未授权', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function refresh(Request $request): JsonResponse
    {
        // 获取登录用户ID (由 AuthJwt 中间件写入请求属性)
        $uid = (int)$request->attributes->get('uid');
        if ($uid <= 0) {
            return response_error(ResultCode::AuthInvalid);
        }

        $userService = new UserService();
        $user = $userService->info($uid);
        // 只需要 states 与响应字段, 排除 password 减少不必要的列传输;
        // 账号状态: 1=正常 0=封禁 (见迁移文件 comment)
        if (empty($user) || (int)($user['states'] ?? 0) !== 1) {
            return response_error(ResultCode::AuthInvalid);
        }

        $token = '';
        try {
            $token = JwtUtil::getInstance()->createToken($uid);
        } catch (\Throwable $e) {
            return response_exception_safe($e, 'auth/refresh');
        }

        return response_success([
            'token' => $token,
            'user' => format_user($user),
        ]);
    }
}
