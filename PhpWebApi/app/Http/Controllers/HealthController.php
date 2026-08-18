<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\DB;

use App\Result\Result;
use App\Result\ResultCode;
use OpenApi\Attributes as OA;

/**
 * 健康检查相关 (与 Go 端 controllers/health.go 对齐)
 */
class HealthController extends Controller
{
    // 服务版本
    private const VERSION = '1.0.0';

    // 健康检查中数据库探测结果的缓存时长(秒)。
    // 探活请求可能非常频繁, 缓存一段时间可显著降低对数据库的压力;
    // 代价是数据库故障后最多延迟该时长才会反映到健康检查结果上。
    private const PING_INTERVAL = 5;

    /**
     * 健康检查
     *
     * 通过 SELECT 1 探测数据库连接, 正常返回 200 + 健康信息;
     * 数据库不可用时返回 503 + 统一错误格式, 便于负载均衡探活。
     */
    #[OA\Get(
        path: '/api/v1/health',
        summary: '健康检查',
        tags: ['Health'],
        responses: [
            new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/HealthResponse')),
            new OA\Response(response: 503, description: 'Service Unavailable', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
        ],
    )]
    public function index(): JsonResponse
    {
        $now = now()->format('Y-m-d H:i:s');
        // 服务启动时间: 进程内首次请求时写入缓存 (缓存被清除后重新填充, 故实际是"缓存初始化时间",
        // 服务器本地时区, 与 created_at 格式保持一致)
        $startedAt = Cache::rememberForever('health_started_at', static fn (): string => now()->format('Y-m-d H:i:s'));

        $database = $this->dbHealthy();
        if ($database !== 'up') {
            return response()->json(new Result(ResultCode::ServiceUnavailable, 'database unavailable'), 503);
        }

        return response()->json([
            'status' => 'ok',      // 服务状态 (固定 "ok")
            'database' => 'up',    // 数据库状态 (固定 "up")
            'version' => self::VERSION, // 服务版本
            'started_at' => $startedAt, // 服务启动时间 (服务器本地时区)
            'now' => $now,             // 当前时间 (服务器本地时区)
        ]);
    }

    /**
     * 探测数据库连接, 结果缓存 PING_INTERVAL 时长。
     *
     * @return string 'up' / 'down'
     */
    private function dbHealthy(): string
    {
        return Cache::remember('health_db', self::PING_INTERVAL, static function (): string {
            try {
                DB::select('SELECT 1');
                return 'up';
            } catch (\Throwable) {
                return 'down';
            }
        });
    }
}
