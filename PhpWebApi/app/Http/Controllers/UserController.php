<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

use App\Result\ResultCode;
use App\Http\Services\UserService;
use App\Http\Request\UserUpdateRequest;

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
        /**
         * 采用 laravel 的 paginate 分页机制，会自动获取请求参数 page
         * GET http://webapi.laixhe.com/api/user/list?page=2
         * POST Content-Type: application/x-www-form-urlencoded
         * POST Content-Type: application/json
         * 只要有 page 参数都可以被  paginate 分页机制获取到
         */
        // 分页当前页数
        //$page = (int) $request->input('page', 0);
        // 每页页数(数量)
        $size = (int)$request->input('size', 0);
        if ($size <= 0) {
            $size = 20;
        }

        $userService = new UserService();
        $dbData = $userService->list($size);

        $users = $dbData->items();
        $data = [];
        foreach ($users as $user) {
            $data[] = [
                'uid' => $user['id'],
                'uname' => $user['uname'],
                'email' => $user['email'],
                'created_at' => $user['created_at'],
            ];
        }
        $result = [
            'list' => $data,
            'total' => $dbData->total(),
            'page' => $dbData->currentPage(),
            'size' => $dbData->perPage(),
        ];

        return response_success($result);
    }

    /**
     * 修改用户信息
     * @param Request $request
     * @return JsonResponse
     */
    public function update(Request $request): JsonResponse
    {
        // 获取登录用户ID
        $uid = (int)$request->header('uid');
        // 获取想要的请求参数
        $req = $request->only([
            'uname',
            'login_at',
        ]);
        $userUpdateRequest = new UserUpdateRequest();
        $error = $userUpdateRequest->validator($req);
        if ($error !== null) {
            return response_result($error);
        }
        //
        $userService = new UserService();
        try {
            $userService->update($uid, $userUpdateRequest);
        } catch (\Throwable $e) {
            return response_exception($e->getCode(), $e->getMessage());
        }
        return response_success();
    }

}
