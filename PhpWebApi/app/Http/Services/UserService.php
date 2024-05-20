<?php

namespace App\Http\Services;

use App\Models\User;

class UserService
{
    public function info(int $uid): array
    {
        // select * from `user` where `id` = ? and `user`.`deleted_at` is null limit 1
        $user = User::query()->where('id', $uid)->first();
        if (empty($user)) {
            return [];
        }
        return $user->toArray();
    }

    public function list(): array
    {
        // select * from `user` where `user`.`deleted_at` is null order by `id` desc
        return User::query()->orderByDesc('id')->get()->toArray();
    }
}
