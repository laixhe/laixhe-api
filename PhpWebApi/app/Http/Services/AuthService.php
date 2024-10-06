<?php

namespace App\Http\Services;

use Throwable;
use RuntimeException;

use Carbon\Carbon;

use App\Http\Request\LoginRequest;
use App\Http\Request\RegisterRequest;
use App\Result\ResultCode;
use App\Models\User;

/**
 * 鉴权服务相关
 */
class AuthService
{
    /**
     * 注册
     * @param RegisterRequest $req
     * @return array
     *
     * @throws RuntimeException
     */
    public function register(RegisterRequest $req): array
    {
        if (!empty(User::query()->where('email', $req->email)->first())) {
            throw new RuntimeException('', ResultCode::EmailExist->value);
        }
        try {
            // sql
            // insert into `user` (`password`, `email`, `uname`, `age`, `login_at`, `updated_at`, `created_at`) values (?, ?, ?, ?, ?, ?, ?)

            // 方式1
//            $user = new User();
//            $user->password = password_hash($req->password, PASSWORD_BCRYPT);
//            $user->email = $req->email;
//            $user->uname = $req->uname;
//            $user->age = $req->age;
//            $user->login_at = Carbon::now();
//            if ($user->save()) {
//                return $user->toArray();
//            }
            // 方式2
            $user = User::query()->create([
                'password' => password_hash($req->password, PASSWORD_BCRYPT),
                'email' => $req->email,
                'uname' => $req->uname,
                'age' => $req->age,
                'login_at' => Carbon::now(),
            ]);
            if ($user !== null){
                return $user->toArray();
            }
        } catch (Throwable $e) {
            throw new RuntimeException($e->getMessage(), ResultCode::Service->value);
        }
        throw new RuntimeException('', ResultCode::Unknown->value);
    }

    /**
     * 登录
     * @param LoginRequest $req
     * @return array
     */
    public function login(LoginRequest $req): array
    {
        // select * from `user` where `email` = ? and `user`.`deleted_at` is null limit 1
        $user = User::query()->where('email', $req->email)->first();
        if (empty($user)) {
            return [];
        }
        // 判断密码是否匹配
        if (!password_verify($req->password, $user->password)) {
            return [];
        }
        return $user->toArray();
    }
}
