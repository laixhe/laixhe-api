<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

use App\Result\ResultCode;
use App\Http\Services\UserService;

class UserController extends Controller
{
    public function info(Request $request)
    {
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

    public function list(Request $request)
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
