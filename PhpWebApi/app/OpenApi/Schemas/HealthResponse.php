<?php

namespace App\OpenApi\Schemas;

use OpenApi\Attributes as OA;

/**
 * 健康检查响应体 (与 Go 端 controllers.HealthResponse 对齐)
 */
#[OA\Schema(
    schema: 'HealthResponse',
    description: '健康检查响应体',
    required: ['status', 'database', 'version', 'started_at', 'now'],
)]
class HealthResponse
{
    #[OA\Property(description: '服务状态 (固定 "ok")', example: 'ok')]
    public string $status;

    #[OA\Property(description: '数据库状态 (固定 "up")', example: 'up')]
    public string $database;

    #[OA\Property(description: '服务版本', example: '1.0.0')]
    public string $version;

    #[OA\Property(description: '服务启动时间 (服务器本地时区)', example: '2026-08-10 12:00:00')]
    public string $started_at;

    #[OA\Property(description: '当前时间 (服务器本地时区)', example: '2026-08-10 12:00:00')]
    public string $now;
}
