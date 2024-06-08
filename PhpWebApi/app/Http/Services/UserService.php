<?php

namespace App\Http\Services;

use App\Models\User;

/**
 * 用户服务相关
 */
class UserService
{
    /**
     * 查询用户信息
     * @param int $uid
     * @return array
     */
    public function info(int $uid): array
    {
        // select * from `user` where `id` = ? and `user`.`deleted_at` is null limit 1
        $user = User::query()->where('id', $uid)->first();
        if (empty($user)) {
            return [];
        }
        return $user->toArray();
    }

    /**
     * 查询用户列表
     * @return array
     */
    public function list(): array
    {
        // select * from `user` where `user`.`deleted_at` is null order by `id` desc
        return User::query()->orderByDesc('id')->get()->toArray();
    }
}
