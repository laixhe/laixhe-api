<?php

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Symfony\Component\HttpKernel\Exception\HttpExceptionInterface;

use App\Result\Result;
use App\Result\ResultCode;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__.'/../routes/web.php',
        api: __DIR__.'/../routes/api.php',
        commands: __DIR__.'/../routes/console.php',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        // 全局 IP 限流中间件 (健康检查路径已豁免)
        $middleware->append(\App\Http\Middleware\RateLimit::class);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        // 自定义错误处理 (与 Go 端 core.ErrorHandler 对齐):
        // 未处理的未知异常 (如数据库异常) 统一返回固定 500 JSON,
        // 避免将内部实现细节泄露给客户端; HTTP 异常 (404/405 等) 保留默认状态码。
        $exceptions->render(function (Throwable $e) {
            if ($e instanceof HttpExceptionInterface) {
                return null; // 交给默认处理器
            }
            return response()->json(new Result(ResultCode::Service, 'internal server error'), 500);
        });
    })->create();
