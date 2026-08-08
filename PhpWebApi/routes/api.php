<?php

use Illuminate\Support\Facades\Route;

use App\Http\Controllers\AuthController;
use App\Http\Controllers\HealthController;
use App\Http\Controllers\SwaggerController;
use App\Http\Controllers\UserController;
use App\Http\Middleware\AssignRequestId;
use App\Http\Middleware\AuthJwt;

Route::middleware(AssignRequestId::class)->prefix('v1')->group(function () {

    // 健康检查 (含数据库探测, 限流已豁免该路径)
    Route::get('health', [HealthController::class, 'index']);

    // Swagger 文档端点
    Route::get('swagger.json', [SwaggerController::class, 'json']);
    Route::get('swagger.yaml', [SwaggerController::class, 'yaml']);

    // 鉴权相关: 注册/登录为公开接口, 仅刷新需要 JWT
    Route::prefix('auth')->group(function () {
        Route::post('register', [AuthController::class, 'register']);
        Route::post('login', [AuthController::class, 'login']);
        Route::post('refresh', [AuthController::class, 'refresh'])->middleware(AuthJwt::class);
    });

    // 用户相关: 获取信息/列表为公开接口, 仅更新需要 JWT
    Route::prefix('user')->group(function () {
        Route::get('info', [UserController::class, 'info']);
        Route::get('list', [UserController::class, 'list']);
        Route::post('update', [UserController::class, 'update'])->middleware(AuthJwt::class);
    });

});
