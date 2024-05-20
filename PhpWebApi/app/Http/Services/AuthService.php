<?php

namespace App\Http\Services;

use Carbon\Carbon;
use Throwable;

use App\Http\Request\LoginRequest;
use App\Http\Request\RegisterRequest;
use App\Result\ResultCode;
use App\Models\User;

class AuthService
{
    public function register(RegisterRequest $req): array
    {
        if (!empty(User::query()->where('email', $req->email)->first())) {
            throw new \RuntimeException('', ResultCode::EmailExist->value);
        }
        try {
            // insert into `user` (`password`, `email`, `uname`, `age`, `login_at`, `updated_at`, `created_at`) values (?, ?, ?, ?, ?, ?, ?)
            $user = new User();
            $user->password = password_hash($req->password, PASSWORD_BCRYPT);
            $user->email = $req->email;
            $user->uname = $req->uname;
            $user->age = $req->age;
            $user->login_at = Carbon::now()->toDateTimeString();
            if ($user->save()) {
                return $user->toArray();
            }
        } catch (Throwable $e) {
            throw new \RuntimeException($e->getMessage(), ResultCode::Service->value);
        }
        throw new \RuntimeException('', ResultCode::Unknown->value);
    }

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
