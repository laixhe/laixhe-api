<?php

namespace App\Http\Services;

use RuntimeException;

use Illuminate\Pagination\LengthAwarePaginator;
use Illuminate\Support\Facades\Cache;

use App\Models\User;
use App\Result\ResultCode;
use App\Http\Requests\UserUpdateRequest;

/**
 * 用户服务相关 (与 Go 端 services/user.go 对齐)
 */
class UserService
{
    /**
     * 查询用户信息 (排除 password)
     *
     * @param int $uid
     * @return array 未找到时返回空数组
     */
    public function info(int $uid): array
    {
        $user = User::query()
            ->select(User::noPassword())
            ->where('id', $uid)
            ->first();
        if (empty($user)) {
            return [];
        }
        return $user->toArray();
    }

    /**
     * 查询用户列表 (排除 password, 按 ID 降序)
     *
     * @param int $page 分页-当前页(默认 1)
     * @param int $pageSize 分页-每页数量(默认 12)
     * @return LengthAwarePaginator
     */
    public function list(int $page, int $pageSize): LengthAwarePaginator
    {
        // total 加 5s 短缓存 (与 Go 端 ListUser 注释、TS 端 getTotalUserCount 同理):
        // count(*) 在 InnoDB 下为全表扫描, 高频翻页时避免重复全表 count;
        // 代价是新增用户后最多延迟 5s 反映到列表 total。
        // 缓存驱动由 CACHE_STORE 决定 (默认 file, 生产高并发建议 redis), 多实例共享同一份缓存。
        $total = (int) Cache::remember('user_total_count', 5, static fn (): int => User::query()->count());

        $users = User::query()
            ->select(User::noPassword())
            ->orderByDesc('id')
            ->forPage($page, $pageSize)
            ->get();

        // 手动构造分页器 (与 paginate() 返回类型一致, 控制器无需改动):
        // forPage 等价于 skip((page-1)*pageSize)->take(pageSize)
        return new LengthAwarePaginator($users, $total, $pageSize, $page);
    }

    /**
     * 修改用户信息 (与 Go 端 User.Update 对齐)
     *
     * 先查后改: 查询 (排除 password) 用于 states 校验与响应组装。
     *
     * @param int $uid 用户id (由 JWT 提供)
     * @param UserUpdateRequest $req
     * @return array 更新后的用户信息
     *
     * @throws RuntimeException 用户不存在 (422) / 账号被封禁 (401)
     */
    public function update(int $uid, UserUpdateRequest $req): array
    {
        $user = User::query()
            ->select(User::noPassword())
            ->where('id', $uid)
            ->first();
        if (empty($user)) {
            throw new RuntimeException('用户不存在', ResultCode::Param->value);
        }
        // 账号状态: 1=正常 0=封禁 (见迁移文件 comment); 封禁返回 401 (与 Go 端一致)
        if ((int)$user->states !== 1) {
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }
        $user->nickname = $req->nickname;
        // 头像地址为空时不更新 (与 Go 端 UpdateUser 按非零字段更新一致)
        if ($req->avatar_url !== '') {
            $user->avatar_url = $req->avatar_url;
        }
        $user->save();
        return $user->toArray();
    }
}
