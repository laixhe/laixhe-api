<?php

namespace App\Providers;

use Godruoyi\Snowflake\LaravelSequenceResolver;
use Godruoyi\Snowflake\Snowflake;
use Illuminate\Foundation\Application;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        // Snowflake 全局单例: id() 内部的序列状态 (同一毫秒内递增) 必须跨请求保留,
        // 若每次 new 都会丢失状态, 同一毫秒内的多次调用可能生成重复 ID。
        // 序列解析器基于缓存原子自增 (LaravelSequenceResolver), 跨进程/实例也能保证同毫秒唯一。
        // 注意: 原子自增依赖缓存 store 支持原子 add/increment (redis/memcached/database 支持;
        // file 缓存非原子, 单进程 FPM 下可用, 但会产生大量永不清理的小文件, 生产建议 CACHE_STORE=redis)。
        // datacenter/worker 由 config/snowflake.php 配置 (默认 -1 构造时随机)。多实例部署时为各实例
        // 分配不同 (datacenter, worker) 组合属于"双保险", 可降低对缓存原子性的依赖, 并非必须
        // (共享原子缓存下相同组合也不会在同毫秒内碰撞)。
        $this->app->singleton(Snowflake::class, function (Application $app) {
            return (new Snowflake(
                (int) config('snowflake.datacenter'),
                (int) config('snowflake.worker'),
            ))->setSequenceResolver(
                // 说明: $app->make('cache')->store() 与 app(\Illuminate\Contracts\Cache\Repository::class)
                // 完全等价。缓存接口类名由 Application::registerCoreContainerAliases() 注册的
                // 容器类名别名解析 ('cache.store' => [Repository::class, Contracts\Cache\Repository::class, ...]),
                // 并非注册于 CacheServiceProvider (Laravel 7.x~13.x 均如此)。详见 README「技术要点」。
                (new LaravelSequenceResolver($app->make('cache')->store()))
                    ->setCachePrefix('snowflake:')
            );
        });
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        //
    }
}
