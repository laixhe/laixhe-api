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
        // select `id`, `uname`, `email`, `created_at` from `user` where `id` = ? and `deleted_at` is null limit 1
        $user = User::query()
            ->select(['id', 'uname', 'email', 'created_at'])
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
    public function list(int $size): LengthAwarePaginator
    {
        // select * from `user` where `deleted_at` is null order by `id` desc
        //return User::query()->orderByDesc('id')->get()->toArray();

        // select count(*) as aggregate from `user` where `deleted_at` is null
        // select `id`, `uname`, `email`, `created_at` from `user` where `deleted_at` is null order by `id` desc limit 20 offset 0
        return User::query()
            ->select(['id', 'uname', 'email', 'created_at'])
            ->orderByDesc('id')
            ->paginate($size);
    }

    /**
     * 修改用户信息
     * @param int $uid
     * @param UserUpdateRequest $req
     * @return void
     */
    public function update(int $uid, UserUpdateRequest $req): void
    {
        // select `id` from `user` where `uname` = ? and `deleted_at` is null limit 1
        $userID = User::query()->where('uname', $req->uname)->value('id');
        if (!empty($userID)) {
            $userID = (int)$userID;
            if ($userID === $uid){
                return;
            }
            throw new RuntimeException('用户名已存在！', ResultCode::Param->value);
        }

        // update `user` set `uname` = ?, `login_at` = ?, `updated_at` = ? where `id` = ? and `deleted_at` is null limit 1
        User::query()->where('id', $uid)->limit(1)->update([
            'uname' => $req->uname,
            'login_at' => $req->login_at
        ]);
    }
}
