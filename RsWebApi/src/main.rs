mod controllers;
mod routers;

// use sqlx::mysql::MySqlPoolOptions;

#[tokio::main]
async fn main() {
    // 初始化跟踪日志
    tracing_subscriber::fmt()
        .with_span_events(tracing_subscriber::fmt::format::FmtSpan::CLOSE)
        .init();

    // 初始化数据库
    // let mysql_pool = MySqlPoolOptions::new()
    //     .max_connections(5)
    //     .connect("mysql://root:123456@127.0.0.1:3306/webapi")
    //     .await
    //     .expect("数据库连接失败");

    // 初始化路由
    let app = routers::init().await.into_make_service();

    // 要绑定的 IP 和 端口
    let addr = "0.0.0.0:3000";
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    tracing::info!("开始监听 http 端口: {}", addr);

    // 启动 http
    axum::serve(listener, app).await.unwrap();
}
