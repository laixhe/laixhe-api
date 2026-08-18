<?php

// JWT 配置 (与 Go 端 config.yaml 的 jwt 节点对齐)
//
// 注意: 业务代码必须通过 config('jwt.*') 读取, 不能直接调 env():
// 生产环境执行 php artisan config:cache 后, config 之外的 env() 一律返回 null,
// 会导致 JWT_SECRET 为空、所有签发/校验 JWT 的接口崩溃。
return [
    // 签名密钥 (生产环境务必通过环境变量注入强随机值)
    'secret' => env('JWT_SECRET', ''),
    // 过期时长 (单位秒), 默认 30 天 (与 .env.example / README 及其它三端配置一致)
    'expire_time' => (int) env('JWT_EXPIRE_TIME', 2592000),
];
