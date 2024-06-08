<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

use App\Result\ResultCode;
use App\Http\Services\UserService;

/**
 * 用户相关
 */
class UserController extends Controller
{
    /**
     * 登录用户信息
     * @param Request $request
     * @return JsonResponse
     */
    public function info(Request $request): JsonResponse
    {
        // 获取登录用户ID
        $uid = (int)$request->header('uid');

        $userService = new UserService();
        $user = $userService->info($uid);
        if (empty($user)) {
            return response_error(ResultCode::AuthNotLogin, '');
        }

        return response_success([
            'user' => [
                'uid' => $user['id'],
                'uname' => $user['uname'],
                'email' => $user['email'],
                'created_at' => $user['created_at'],
            ],
        ]);
    }

    /**
     * 查询用户列表
     * @param Request $request
     * @return JsonResponse
     */
    public function list(Request $request): JsonResponse
    {
        $userService = new UserService();
        $users = $userService->list();
        $data = [];
        foreach ($users as $user) {
            $data[] = [
                'uid' => $user['id'],
                'uname' => $user['uname'],
                'email' => $user['email'],
                'created_at' => $user['created_at'],
            ];
        }
        return response_success($data);
    }

}
