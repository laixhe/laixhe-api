<?php

namespace App\Http\Services;

use RuntimeException;

use Illuminate\Contracts\Pagination\LengthAwarePaginator;

use App\Models\User;
use App\Result\ResultCode;
use App\Http\Request\UserUpdateRequest;

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
        // select * from `user` where `id` = ? limit 1
        $user = User::query()
            ->where('id', $uid)
            ->first();
        if (empty($user)) {
            return [];
        }
        return $user->toArray();
    }

    /**
     * 查询用户列表
     * @param int $size 每页页数(数量)
     * @return LengthAwarePaginator
     */
    public function list(int $limit): LengthAwarePaginator
    {
        // select count(*) as aggregate from `user`
        // select * from `user` order by `id` desc limit 2 offset 0
        return User::query()
            ->orderByDesc('id')
            ->paginate($limit);
    }

    /**
     * 修改用户信息
     * @param int $uid
     * @param UserUpdateRequest $req
     * @return void
     */
    public function update(int $uid, UserUpdateRequest $req): void
    {
        // select `id` from `user` where `nickname` = ? limit 1
        $userID = User::query()->where('nickname', $req->nickname)->value('id');
        if (!empty($userID)) {
            $userID = (int)$userID;
            if ($userID === $uid){
                return;
            }
            throw new RuntimeException('用户名已存在！', ResultCode::Param->value);
        }
        // update `user` set `nickname` = ?, `avatar_url` = ?, `user`.`updated_at` = ? where `id` = ? limit 1
        User::query()->where('id', $uid)->limit(1)->update([
            'nickname' => $req->nickname,
            'avatar_url' => $req->avatar_url
        ]);
    }
}
