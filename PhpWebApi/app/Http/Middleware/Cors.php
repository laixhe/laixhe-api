<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

/**
 * CORS 中间件 (与 Go/Rust/TS 端 UseCors/响应头对齐)
 *
 * 前后端分离部署时允许跨域访问。为教学演示采用宽泛配置 (允许任意来源);
 * 生产环境请按需收紧 origin 白名单。
 */
class Cors
{
    public function handle(Request $request, Closure $next): Response
    {
        $response = $next($request);
        $response->headers->set('Access-Control-Allow-Origin', '*');
        $response->headers->set('Access-Control-Allow-Methods', 'GET, POST, PUT, PATCH, DELETE, OPTIONS');
        $response->headers->set('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Request-Id');

        return $response;
    }
}
