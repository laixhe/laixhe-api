<?php

use Illuminate\Support\Facades\Route;

use App\Http\Controllers\AuthController;
use App\Http\Controllers\UserController;
use App\Http\Middleware\AssignRequestId;
use App\Http\Middleware\AuthJwt;

Route::middleware(AssignRequestId::class)->prefix('v1')->group(function () {

    Route::prefix('auth')->group(function () {
        Route::post('register', [AuthController::class, 'register']);
        Route::post('login', [AuthController::class, 'login']);
        Route::post('refresh', [AuthController::class, 'refresh'])->middleware(AuthJwt::class);
    });

    Route::prefix('user')->middleware(AuthJwt::class)->group(function () {
        Route::get('info', [UserController::class, 'info']);
        Route::get('list', [UserController::class, 'list']);
        Route::post('update', [UserController::class, 'update']);
    });

});

