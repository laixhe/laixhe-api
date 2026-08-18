<?php

namespace App\OpenApi\Schemas;

use OpenApi\Attributes as OA;

/**
 * 用户列表响应体
 */
#[OA\Schema(
    schema: 'UserListResponse',
    description: '用户列表响应体',
    required: ['total', 'page', 'page_size', 'list'],
)]
class UserListResponse
{
    #[OA\Property(description: '总数', example: 100)]
    public int $total;

    #[OA\Property(description: '分页-当前页', example: 1)]
    public int $page;

    #[OA\Property(description: '分页-每页数量', example: 12)]
    public int $page_size;

    #[OA\Property(type: 'array', items: new OA\Items(ref: '#/components/schemas/User'), description: '列表')]
    public array $list;
}
