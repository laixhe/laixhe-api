mod controllers;
mod routers;

#[tokio::main]
async fn main() {
    // 初始化跟踪日志
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .with_span_events(tracing_subscriber::fmt::format::FmtSpan::CLOSE)
        .init();

    // 初始化路由
    let app = routers::init()
        .await.
        into_make_service();

    // 要绑定的 IP 和 端口
    let addr = "0.0.0.0:5050";
    let listener = tokio::net::TcpListener::bind(addr)
        .await
        .unwrap();

    tracing::info!("开始监听 http 端口: {}", addr);

    // 启动 http
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

async fn shutdown_signal() {
    match tokio::signal::ctrl_c().await {
        Ok(_val) => tracing::event!(tracing::Level::INFO, "shutting down systems"),
        Err(error) => tracing::event!(tracing::Level::DEBUG, "{error:?}"),
    };
}
