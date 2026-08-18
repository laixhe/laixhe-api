<?php

namespace App\OpenApi\Schemas;

use OpenApi\Attributes as OA;

/**
 * 统一错误响应体 (与 Go 端 core.Error 对齐)
 */
#[OA\Schema(
    schema: 'Error',
    description: '统一错误响应体',
    required: ['code', 'message'],
)]
class Error
{
    #[OA\Property(description: '错误码 (与 HTTP 状态码一致)', example: 422)]
    public int $code;

    #[OA\Property(description: '错误描述', example: '参数错误')]
    public string $message;
}
