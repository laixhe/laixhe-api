<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Symfony\Component\HttpFoundation\Response;
use Godruoyi\Snowflake\Snowflake;

/**
 * 为每个请求分配唯一 request_id:
 * 写入日志上下文 (便于日志串联排查), 并回写 X-Request-Id 响应头 (与 Go 端 requestId 中间件对齐)。
 */
class AssignRequestId
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        // Snowflake 生成全局唯一请求 ID
        $requestId = (new Snowflake())->id();
        Log::withContext([
            'request_id' => $requestId
        ]);

        $response = $next($request);
        $response->headers->set('X-Request-Id', $requestId);

        return $response;
    }
}
