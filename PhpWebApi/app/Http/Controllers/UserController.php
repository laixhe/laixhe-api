<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

use App\Result\ResultCode;
use App\Http\Services\UserService;
use App\Http\Requests\UserUpdateRequest;
use OpenApi\Attributes as OA;

/**
 * 用户相关 (与 Go 端 controllers/user.go 对齐)
 */
class UserController extends Controller
{
    /**
     * 获取用户信息 (公开接口, 不受 JWT 保护)
     */
    #[OA\Get(
        path: '/api/v1/user/info',
        summary: '获取用户信息',
        tags: ['User'],
        parameters: [
            new OA\QueryParameter(name: 'uid', description: '用户id', required: true, schema: new OA\Schema(type: 'integer')),
        ],
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/User')),
            new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 422, description: '参数错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function info(Request $request): JsonResponse
    {
        // 非数字 uid 按请求格式错误返回 400 (与 Go/Rust 端绑定层行为一致);
        // 数字但非法 (<=0) 仍走下方 422 "无效的用户ID"
        $uidRaw = $request->input('uid', 0);
        if ($uidRaw !== 0 && filter_var($uidRaw, FILTER_VALIDATE_INT) === false) {
            return response_error(ResultCode::BadRequest, '无效的用户ID');
        }
        $uid = (int)$uidRaw;
        if ($uid <= 0) {
            return response_error(ResultCode::Param, '无效的用户ID');
        }

        $userService = new UserService();
        $user = $userService->info($uid);
        if (empty($user)) {
            return response_error(ResultCode::Param, '用户不存在');
        }

        return response_success(format_user($user));
    }

    /**
     * 获取用户列表 (公开接口, 不受 JWT 保护)
     */
    #[OA\Get(
        path: '/api/v1/user/list',
        summary: '获取用户列表',
        tags: ['User'],
        parameters: [
            new OA\QueryParameter(name: 'page', description: '分页-当前页(默认 1)', required: false, schema: new OA\Schema(type: 'integer')),
            new OA\QueryParameter(name: 'page_size', description: '分页-每页数量(默认 12)', required: false, schema: new OA\Schema(type: 'integer')),
        ],
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/UserListResponse')),
            new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 422, description: '参数错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function list(Request $request): JsonResponse
    {
        // 非数字分页参数按请求格式错误返回 400 (与 Go/Rust 端绑定层行为一致);
        // 仅缺省/0 走下方归一化钳制
        $pageRaw = $request->input('page', 0);
        $pageSizeRaw = $request->input('page_size', 0);
        if ($pageRaw !== 0 && filter_var($pageRaw, FILTER_VALIDATE_INT) === false) {
            return response_error(ResultCode::BadRequest, '无效的分页参数');
        }
        if ($pageSizeRaw !== 0 && filter_var($pageSizeRaw, FILTER_VALIDATE_INT) === false) {
            return response_error(ResultCode::BadRequest, '无效的分页参数');
        }

        // 分页-当前页(默认 1)
        $page = (int)$pageRaw;
        if ($page <= 0) {
            $page = 1;
        }
        // 分页-每页数量(默认 12)
        $page_size = (int)$pageSizeRaw;
        if ($page_size <= 0) {
            $page_size = 12;
        }
        // 上限保护: 防止超大 page_size 触发全量查询
        if ($page_size > 100) {
            $page_size = 100;
        }

        $userService = new UserService();
        $dbData = $userService->list($page, $page_size);
        //
        $data = [];
        foreach ($dbData->items() as $user) {
            $data[] = format_user($user->toArray());
        }
        $result = [
            'total' => $dbData->total(),
            'page' => $dbData->currentPage(),
            'page_size' => $dbData->perPage(),
            'list' => $data,
        ];
        return response_success($result);
    }

    /**
     * 更新用户信息 (需要 JWT, Uid 由 JWT 提供)
     */
    #[OA\Post(
        path: '/api/v1/user/update',
        summary: '更新用户信息',
        tags: ['User'],
        security: [['BearerAuth' => []]],
        requestBody: new OA\RequestBody(
            required: true,
            content: new OA\JsonContent(ref: '#/components/schemas/UserUpdateRequest')
        ),
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/User')),
            new OA\Response(response: 401, description: '未授权', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 422, description: '参数错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
            new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function update(Request $request): JsonResponse
    {
        // 顶层 body 类型校验 (数组/标量/null → 400; 空 body → 422 兜底): 见 Controller::validateTopLevelBody
        $topLevelError = $this->validateTopLevelBody($request);
        if ($topLevelError !== null) {
            return response_result($topLevelError);
        }

        // 获取登录用户ID (由 AuthJwt 中间件写入请求属性)
        $uid = (int)$request->attributes->get('uid');
        if ($uid <= 0) {
            return response_error(ResultCode::AuthInvalid);
        }
        // 获取想要的请求参数
        $req = $request->only([
            'nickname',
            'avatar_url',
        ]);
        $userUpdateRequest = new UserUpdateRequest();
        $error = $userUpdateRequest->validator($req);
        if ($error !== null) {
            return response_result($error);
        }
        // 头像地址非空时必须以 http:// 或 https:// 开头 (与 Go 端一致;
        // 用精确前缀匹配, 避免 httpxxx:// 之类的畸形 scheme 通过)
        $avatarUrl = $req['avatar_url'] ?? '';
        if ($avatarUrl !== '' && !str_starts_with($avatarUrl, 'http://') && !str_starts_with($avatarUrl, 'https://')) {
            return response_error(ResultCode::Param, '头像地址必须以http或https开头');
        }
        //
        $userService = new UserService();
        try {
            $user = $userService->update($uid, $userUpdateRequest);
        } catch (\Throwable $e) {
            // 已知业务错误 (如"用户不存在") 透传; 未知异常 (含 SQL) 记日志后返回固定文案 (见 Helpers)
            return response_exception_safe($e, 'user/update');
        }
        return response_success(format_user($user));
    }

}
