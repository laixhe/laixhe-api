<?php

namespace App\OpenApi\Schemas;

use OpenApi\Attributes as OA;

/**
 * 用户信息响应体 (与 Go 端 entity.User 对齐)
 */
#[OA\Schema(
    schema: 'User',
    description: '用户信息',
    required: ['uid', 'type_id', 'account', 'mobile', 'email', 'nickname', 'avatar_url', 'sex', 'states', 'created_at'],
)]
class User
{
    #[OA\Property(description: '用户id', example: 1)]
    public int $uid;

    #[OA\Property(description: '类型 1-普通用户', example: 1, enum: [1])]
    public int $type_id;

    #[OA\Property(description: '账号', example: 'xid00000000000000000000')]
    public string $account;

    #[OA\Property(description: '手机号', example: '')]
    public string $mobile;

    #[OA\Property(description: '邮箱', example: 'user@example.com')]
    public string $email;

    #[OA\Property(description: '昵称', example: '张三')]
    public string $nickname;

    #[OA\Property(description: '头像地址', example: '')]
    public string $avatar_url;

    #[OA\Property(description: '性别 (0-未填写 1-男 2-女)', example: 1, enum: [0, 1, 2])]
    public int $sex;

    #[OA\Property(description: '状态 (0-禁用 1-正常)', example: 1, enum: [0, 1])]
    public int $states;

    #[OA\Property(description: '创建时间, 格式 "Y-m-d H:i:s"', example: '2026-08-10 12:00:00')]
    public string $created_at;
}
