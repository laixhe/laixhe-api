<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

use App\Result\ResultCode;
use App\Http\Services\UserService;
use App\Http\Requests\UserUpdateRequest;

/**
 * 用户相关 (与 Go 端 controllers/user.go 对齐)
 */
class UserController extends Controller
{
    /**
     * 获取用户信息 (公开接口, 不受 JWT 保护)
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function info(Request $request): JsonResponse
    {
        $uid = (int)$request->input('uid', 0);
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
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function list(Request $request): JsonResponse
    {
        // 分页-当前页(默认 1)
        $page = (int)$request->input('page', 0);
        if ($page <= 0) {
            $page = 1;
        }
        // 分页-每页数量(默认 12)
        $page_size = (int)$request->input('page_size', 0);
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
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function update(Request $request): JsonResponse
    {
        // 获取登录用户ID (由 AuthJwt 中间件写入请求头)
        $uid = (int)$request->header('uid');
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
        // 头像地址非空时必须以 http 或 https 开头 (与 Go 端一致)
        $avatarUrl = $req['avatar_url'] ?? '';
        if ($avatarUrl !== '' && !str_starts_with($avatarUrl, 'http')) {
            return response_error(ResultCode::Param, '头像地址必须以http或https开头');
        }
        //
        $userService = new UserService();
        try {
            $user = $userService->update($uid, $userUpdateRequest);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
        return response_success(format_user($user));
    }

}
