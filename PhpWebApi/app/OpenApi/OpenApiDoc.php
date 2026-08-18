<?php

namespace App\OpenApi;

use OpenApi\Attributes as OA;

/**
 * OpenAPI 根注解 (与 Go 端 `swag init` 生成的 docs 对齐)
 *
 * 生成方式: `composer swagger` (见 scripts/generate-swagger.php)
 */
#[OA\OpenApi(
    security: [['BearerAuth' => []]],
)]
#[OA\Info(
    version: '1.0',
    title: 'API接口',
    description: '用户认证与用户管理 API 服务。注册/登录成功后返回 JWT 令牌 (HS256 签名, 含 uid/exp/iat/nbf 声明), 过期时长由 JWT_EXPIRE_TIME 控制。受保护接口需在请求头携带 Authorization: Bearer <token>。',
)]
#[OA\SecurityScheme(
    securityScheme: 'BearerAuth',
    type: 'http',
    scheme: 'bearer',
    bearerFormat: 'JWT',
    description: '在请求头携带 Authorization: Bearer <token>',
)]
class OpenApiDoc
{
}
