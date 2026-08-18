<?php

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Symfony\Component\HttpFoundation\Response;
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
        // CORS (与 Go/Rust/TS 端 UseCors 对齐)
        $middleware->append(\App\Http\Middleware\Cors::class);
        // 全局 IP 限流中间件 (健康检查路径已豁免)
        $middleware->append(\App\Http\Middleware\RateLimit::class);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        // 自定义错误处理 (与 Go 端 core.ErrorHandler 对齐):
        // HTTP 异常 (404 等) 返回统一 JSON 格式 {code,message} (与 Rust/TS 端对齐);
        // 未处理的未知异常 (如数据库异常) 统一返回固定 500 JSON, 避免泄露内部细节。
        $exceptions->render(function (Throwable $e) {
            if ($e instanceof HttpExceptionInterface) {
                // 使用标准状态码文案 (如 404→"Not Found"), 与 Rust/TS 端统一 JSON 格式对齐
                // (不用 $e->getMessage(), 否则 Laravel 会返回 "The route ... could not be found." 之类)
                $status = $e->getStatusCode();
                $message = Response::$statusTexts[$status] ?? 'Not Found';
                return response()->json(new Result(ResultCode::intToEnum($status), $message), $status);
            }
            return response()->json(new Result(ResultCode::Service, 'internal server error'), 500);
        });
    })->create();
