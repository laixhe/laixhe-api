<?php

return [

    /*
    |--------------------------------------------------------------------------
    | IP 接口限流 (与 Go 端 config.yaml limit 对齐)
    |--------------------------------------------------------------------------
    |
    | 单个 IP 在滑动窗口内允许的最大请求数, 超过阈值返回 429。
    | 健康检查路径 (/api/v1/health) 豁免限流, 避免负载均衡探活被误伤。
    |
    */

    // 是否启用接口限流
    'enable' => filter_var(env('RATE_LIMIT_ENABLE', true), FILTER_VALIDATE_BOOLEAN),

    // 单个 IP 在窗口内允许的最大请求数
    'max' => (int) env('RATE_LIMIT_MAX', 1000),

    // 滑动窗口时长(单位秒)
    'window' => (int) env('RATE_LIMIT_WINDOW', 60),

];
