//! WebApi 服务入口 (对应 Go 项目 main.go)
//!
//! # 模块结构
//!
//! | 模块 | 说明 |
//! | ---- | ---- |
//! | `config` | 配置加载 (YAML + `${ENV}` 展开) 与校验 |
//! | `state` | `AppState`：日志 / 数据库连接池 / 限流器的初始化与持有 |
//! | `routes` | 路由注册与全局中间件编排 |
//! | `middleware` | 请求日志 (X-Request-ID) / JWT 鉴权 / IP 限流 |
//! | `app::controllers` | 控制器层：参数校验 + 直接返回业务实体 (`Result<Json<T>, ApiError>`) |
//! | `app::services` | 业务逻辑层（含耗时日志） |
//! | `app::models` | sea-orm 数据库实体 |
//! | `logger` | 统一耗时日志模块 (`log_elapsed!` 宏) |
//! | `error` | 统一错误格式 (`{"code": <int>, "message": "<string>"}`, 与 Go fiber.Error 一致) |
//!
//! # 启动流程
//!
//! 1. 解析命令行参数 `--config=<file>` (默认 `./config.yaml`)
//! 2. 加载配置 → 初始化日志 → 初始化数据库连接池
//! 3. 从数据库 `config_common` 表加载运行时配置
//! 4. 构建路由，监听 HTTP，等待退出信号触发优雅停机
//!
//! # 使用
//!
//! ```bash
//! webapi --config=./config.yaml
//! ```

mod app;
mod config;
mod error;
mod logger;
mod middleware;
mod routes;
mod state;

#[cfg(test)]
mod tests;

use state::AppState;

#[tokio::main]
async fn main() {
    let config_file = parse_args();
    // 指定版本号 (编译时可通过 env 注入: GIT_VERSION=xxx, 对应 Go 的 -ldflags "-X main.GitVersion=xxx")
    let git_version = option_env!("GIT_VERSION").unwrap_or("");

    // hostname: 优先 COMPUTERNAME (Windows), 回退 HOSTNAME (Linux/macOS), 对齐 Go os.Hostname() 的跨平台行为
    let hostname = std::env::var("COMPUTERNAME")
        .or_else(|_| std::env::var("HOSTNAME"))
        .unwrap_or_else(|_| "unknown".to_string());
    println!(
        "[rust version: {}] [git: {}] [config file: {}] [hostname: {}]",
        env!("CARGO_PKG_VERSION"),
        git_version,
        config_file,
        hostname
    );

    // 创建服务: 加载配置 → 初始化日志 → 初始化 ORM
    let state = AppState::new(&config_file).await;
    // 从数据库 config_common 表加载运行时配置
    app::services::init_config_common(&state).await;

    // 构建路由并启动 HTTP 服务 (支持优雅停机)
    let app = routes::build(state.clone());
    let addr = state.config.http.addr();
    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .unwrap_or_else(|e| panic!("listen {addr} failed: {e}"));
    tracing::info!("server listening on http://{addr}");
    // into_make_service_with_connect_info: 注入客户端 SocketAddr, 供限流中间件获取真实 IP
    if let Err(e) = axum::serve(
        listener,
        app.into_make_service_with_connect_info::<std::net::SocketAddr>(),
    )
    .with_graceful_shutdown(shutdown_signal())
    .await
    {
        panic!("server error: {e}");
    }
}

/// 等待退出信号: Ctrl+C / SIGTERM, 触发优雅停机
async fn shutdown_signal() {
    #[cfg(unix)]
    {
        let mut term = tokio::signal::unix::signal(tokio::signal::unix::SignalKind::terminate())
            .expect("install SIGTERM handler failed");
        tokio::select! {
            _ = tokio::signal::ctrl_c() => tracing::info!("收到 Ctrl+C, 开始优雅停机"),
            _ = term.recv() => tracing::info!("收到 SIGTERM, 开始优雅停机"),
        }
    }
    #[cfg(not(unix))]
    {
        tokio::signal::ctrl_c()
            .await
            .expect("install ctrl_c handler failed");
        tracing::info!("收到 Ctrl+C, 开始优雅停机");
    }
}

/// 解析命令行参数，返回配置文件路径
///
/// 优先级：
/// 1. `--config=xxx.yaml` 显式指定
/// 2. `--env=dev` 时加载 `config.dev.yaml`（存在才使用，否则回退默认）
/// 3. 默认 `./config.yaml`
fn parse_args() -> String {
    let mut config_file = "./config.yaml".to_string();
    let mut env_name: Option<String> = None;
    for arg in std::env::args().skip(1) {
        if let Some(v) = arg.strip_prefix("--config=") {
            config_file = v.to_string();
        } else if let Some(v) = arg.strip_prefix("--env=") {
            env_name = Some(v.to_string());
        }
    }
    // 未显式指定 config 且存在 config.{env}.yaml 时按环境加载
    if env_name.is_some() && config_file == "./config.yaml" {
        let env_file = format!("./config.{}.yaml", env_name.unwrap());
        if std::path::Path::new(&env_file).exists() {
            config_file = env_file;
        }
    }
    config_file
}
