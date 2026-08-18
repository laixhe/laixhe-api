<?php

namespace App\OpenApi\Schemas;

use OpenApi\Attributes as OA;

/**
 * 鉴权成功响应体 (注册/登录/刷新)
 */
#[OA\Schema(
    schema: 'AuthTokenResponse',
    description: '鉴权成功响应体',
    required: ['token', 'user'],
)]
class AuthTokenResponse
{
    #[OA\Property(description: 'jwt token', example: 'eyJhbGciOiJIUzI1NiIs...')]
    public string $token;

    #[OA\Property(ref: '#/components/schemas/User', description: '用户信息')]
    public object $user;
}
