<?php

return [

    /*
    |--------------------------------------------------------------------------
    | Snowflake ID 生成 (godruoyi/php-snowflake)
    |--------------------------------------------------------------------------
    |
    | datacenter / worker 用于在多实例部署时划分 ID 空间 (范围 0-31, 或 -1 随机)。
    | 单实例部署时保持 -1 (构造时随机分配) 即可;
    | 多实例部署时建议为各实例分配不同的 (datacenter, worker) 组合——这是"双保险":
    | 即使不区分, 序列解析器 (LaravelSequenceResolver) 在共享原子缓存下也能保证同毫秒唯一。
    |
    */

    // 数据中心 ID (0-31), -1 表示随机分配
    'datacenter' => (int) env('SNOWFLAKE_DATACENTER', -1),

    // 工作节点 ID (0-31), -1 表示随机分配
    'worker' => (int) env('SNOWFLAKE_WORKER', -1),

];
