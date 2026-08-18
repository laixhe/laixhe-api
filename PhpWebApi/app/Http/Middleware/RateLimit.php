<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Symfony\Component\HttpFoundation\Response;

use App\Result\ResultCode;

/**
 * IP 接口限流中间件 (与 Go 端 rate_limit 对齐)
 *
 * 基于滑动窗口统计每个 key (IP) 在窗口内的请求次数,
 * 超过阈值时返回 429 统一 JSON。
 *
 * 性能注意: 每个请求都会对窗口内时间戳数组做一次缓存读+写 (O(窗口内请求数)),
 * 开发环境用默认 file 缓存即可; 生产高并发建议将 CACHE_STORE 设为 redis,
 * 避免限流本身成为性能瓶颈。
 *
 * 正确性注意: file 缓存驱动的 get/put 不是原子的 (read-modify-write 无锁),
 * 并发请求可能互相覆盖计数, 极端情况下可绕过限流阈值;
 * 生产环境应使用 redis 等支持原子操作 (或原子自增) 的缓存驱动。
 */
class RateLimit
{
    // 健康检查路径 (限流豁免, 避免负载均衡探活被误伤)
    private const HEALTH_PATH = 'api/v1/health';

    /**
     * Handle an incoming request.
     */
    public function handle(Request $request, Closure $next): Response
    {
        // 配置关闭限流时直接放行
        if (!config('rate_limit.enable', true)) {
            return $next($request);
        }
        // 健康检查路径豁免限流
        if ($request->path() === self::HEALTH_PATH) {
            return $next($request);
        }
        $ip = $this->resolveClientIP($request);
        if (!$this->check($ip)) {
            // {"code": 429, "message": "请求过于频繁，请稍后再试"}
            return response_error(ResultCode::TooManyRequests);
        }
        return $next($request);
    }

    /**
     * 滑动窗口计数: key 在窗口内的请求次数是否未超限
     *
     * @param string $key 客户端 IP
     * @return bool true 表示未超限 (计数已入窗), false 表示超限
     */
    private function check(string $key): bool
    {
        $max = (int)config('rate_limit.max', 1000);
        $window = (int)config('rate_limit.window', 60);

        $cacheKey = 'rate_limit:' . $key;
        $times = Cache::get($cacheKey, []);
        $now = time();
        // 移除窗口外的旧记录
        $times = array_values(array_filter($times, static fn (int $t): bool => ($now - $t) < $window));

        if (count($times) >= $max) {
            Cache::put($cacheKey, $times, $window);
            return false;
        }
        $times[] = $now;
        Cache::put($cacheKey, $times, $window);

        return true;
    }

    /**
     * 解析客户端 IP: 优先代理头 X-Forwarded-For, 其次真实连接地址, 兜底 "unknown"
     *
     * 注意: 直接信任 X-Forwarded-For 时客户端可伪造该头绕过限流,
     * 生产环境建议仅在有可信反向代理时启用, 或直接只使用真实 IP。
     */
    private function resolveClientIP(Request $request): string
    {
        $xff = trim((string)$request->header('X-Forwarded-For'));
        if ($xff !== '') {
            $first = trim(explode(',', $xff)[0]);
            if ($first !== '') {
                return $first;
            }
        }
        $ip = $request->ip();
        if ($ip !== null && $ip !== '') {
            return $ip;
        }
        return 'unknown';
    }
}
