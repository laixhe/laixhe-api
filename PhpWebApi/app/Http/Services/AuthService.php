<?php

namespace App\Http\Services;

use Throwable;
use RuntimeException;

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
            throw new RuntimeException('邮箱已存在', ResultCode::Param->value);
        }
        try {
            $user = User::query()->create([
                'type_id' => 1,
                'account' => '',
                'mobile' => '',
                'email' => $req->email,
                'password' => password_hash($req->password, PASSWORD_BCRYPT),
                'nickname' => $req->nickname,
                'avatar_url' => '',
                'sex' => 2,
                'states' => 1,
            ]);
            if ($user !== null){
                return $user->toArray();
            }
        } catch (Throwable $e) {
            throw new RuntimeException($e->getMessage(), ResultCode::Service->value);
        }
        throw new RuntimeException('', ResultCode::Service->value);
    }

    /**
     * 登录
     * @param LoginRequest $req
     * @return array
     */
    public function login(LoginRequest $req): array
    {
        // select * from `user` where `email` = ? limit 1
        $user = User::query()->where('email', $req->email)->first();
        if (empty($user)) {
            return [];
        }
        return $user->toArray();
    }
}
