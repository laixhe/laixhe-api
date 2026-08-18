//! 服务与应用状态 (对应 Go 项目 core/server.go + core/orm.go)
//!
//! 持有 Config、sea-orm 数据库连接和运行时通用配置。

use std::sync::{Arc, Mutex};
use std::time::Duration;

use regex::Regex;
use sea_orm::{ConnectOptions, Database, DatabaseConnection, DbErr};
use tracing_subscriber::filter::LevelFilter;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;
use tracing_subscriber::EnvFilter;

use crate::config::{CommonConfig, Config, LogConfig, OrmConfig};
use crate::log_elapsed;
use crate::logger::Timer;
use crate::middleware::rate_limit::RateLimiter;

/// 日志文件写入器 guard: 静态持有到进程退出, 退出时自动 drop 并 flush 尾部缓冲日志
/// (替代 mem::forget, 避免优雅停机时丢失 file 模式的尾部日志)
static LOG_GUARD: Mutex<Option<tracing_appender::non_blocking::WorkerGuard>> = Mutex::new(None);

/// 应用状态
#[derive(Clone)]
pub struct AppState {
    pub config: Arc<Config>,
    /// sea-orm 数据库连接 (对应 s.orm["default"])
    pub db: DatabaseConnection,
    /// 运行时通用配置，从数据库 config_common 表动态加载
    pub common: Arc<Mutex<CommonConfig>>,
    /// IP 限流器
    pub limiter: Arc<RateLimiter>,
}

impl AppState {
    /// 创建服务: 加载配置 → 初始化日志 → 初始化 ORM
    pub async fn new(config_file: &str) -> AppState {
        let start = Timer::new();
        let config = Arc::new(
            Config::load(config_file).unwrap_or_else(|e| panic!("load config failed: {e}")),
        );
        // 初始化日志
        init_log(&config.log);
        // 初始化 ORM 数据库连接池
        let step = Timer::new();
        let db = init_orm(&config.orm)
            .await
            .unwrap_or_else(|e| panic!("init orm failed: {e}"));
        log_elapsed!(step, elapsed_ms, info, "数据库连接池建立完成");
        log_elapsed!(start, total_ms, info, "服务初始化完成");
        // 限流器: 参数取 config.limit (整段缺省时 Config::load 回填 {enable, 1000, 60}, 见 config.rs)
        // 在 config Arc 移入 AppState 之前取值
        let limiter = Arc::new(RateLimiter::new(
            config.limit.max,
            Duration::from_secs(config.limit.window),
        ));
        AppState {
            config,
            db,
            common: Arc::new(Mutex::new(CommonConfig::default())),
            limiter,
        }
    }
}

/// 初始化日志 (对应 xlog.InitZap)
///
/// - `run`:
///   - `console` 模式: 输出到控制台
///   - `file` 模式: 按大小轮转写入文件 (对应 lumberjack 的 max_size / max_backups)
/// - `format`:
///   - `console`: 人类可读文本 (默认)
///   - `json`: 结构化 JSON 日志，便于 ELK / Loki 等采集
pub fn init_log(log_config: &LogConfig) {
    let level = match log_config.level.as_str() {
        "debug" => LevelFilter::DEBUG,
        "warn" => LevelFilter::WARN,
        "error" => LevelFilter::ERROR,
        _ => LevelFilter::INFO,
    };
    let filter = EnvFilter::new(level.to_string());
    let is_json = log_config.format == "json";

    // console 模式: 输出到控制台 (同样接 non_blocking, 避免高 QPS 下同步写 stdout 阻塞请求线程)
    // 注意: guard 必须静态持有! 若直接丢弃, WorkerGuard::drop 会关闭 worker 线程,
    // 导致之后的所有日志静默丢失 (tracing-appender lossy 模式直接丢弃)
    if log_config.run != "file" {
        let (non_blocking, guard) = tracing_appender::non_blocking(std::io::stdout());
        *crate::sync::lock_unpoison(&LOG_GUARD) = Some(guard);
        if is_json {
            let layer = tracing_subscriber::fmt::layer()
                .with_writer(non_blocking)
                .json()
                .with_ansi(false);
            let _ = tracing_subscriber::registry()
                .with(filter)
                .with(layer)
                .try_init();
        } else {
            let layer = tracing_subscriber::fmt::layer()
                .with_writer(non_blocking)
                .with_target(false)
                .with_ansi(false);
            let _ = tracing_subscriber::registry()
                .with(filter)
                .with(layer)
                .try_init();
        }
        return;
    }

    // file 模式: 按大小轮转写入文件 (max_size(MB) → 字节, max_backups 保留备份数)
    let max_size_bytes = (log_config.max_size * 1024 * 1024) as u64;
    let condition = rolling_file::RollingConditionBasic::new().max_size(max_size_bytes);
    let appender = rolling_file::RollingFileAppender::new(
        &log_config.path,
        condition,
        log_config.max_backups.max(1) as usize,
    )
    .unwrap_or_else(|e| panic!("create rolling file appender failed: {e}"));
    let (non_blocking, guard) = tracing_appender::non_blocking(appender);
    // 静态持有 guard: 进程退出时自动 drop 并 flush 尾部缓冲日志 (替代 mem::forget, 避免丢尾部日志);
    // 重复初始化(测试场景)时覆盖旧 guard, 旧 guard drop 即 flush
    *crate::sync::lock_unpoison(&LOG_GUARD) = Some(guard);
    // try_init: 重复初始化(如测试场景)时不 panic
    if is_json {
        let layer = tracing_subscriber::fmt::layer()
            .with_writer(non_blocking)
            .json()
            .with_ansi(false);
        let _ = tracing_subscriber::registry()
            .with(filter)
            .with(layer)
            .try_init();
    } else {
        let layer = tracing_subscriber::fmt::layer()
            .with_writer(non_blocking)
            .with_ansi(false);
        let _ = tracing_subscriber::registry()
            .with(filter)
            .with(layer)
            .try_init();
    }
}

/// 初始化 ORM 数据库连接 (对应 s.initOrm / mysql.Init)
async fn init_orm(orm: &OrmConfig) -> Result<DatabaseConnection, DbErr> {
    let start = Timer::new();
    let dsn = dsn_to_url(&orm.driver, &orm.dsn);
    // 日志只输出脱敏后的 dsn, 避免数据库账号密码落盘 (config.yaml 的 dsn 含明文凭据)
    tracing::info!(
        driver = %orm.driver,
        dsn = %mask_dsn(&dsn),
        max_connections = orm.max_open_count,
        min_connections = orm.max_idle_count,
        max_life_time_secs = orm.max_life_time,
        "开始初始化数据库连接池"
    );
    let mut opt = ConnectOptions::new(dsn);
    opt.max_connections(orm.max_open_count)
        .min_connections(orm.max_idle_count)
        // max_lifetime: 连接总寿命 (对齐 Go gorm SetConnMaxLifetime)。
        // 注意: 不用 sea-orm 的 idle_timeout 选项: 它只回收"空闲"连接, 语义不同,
        // 活跃连接永不回收, 存在被 MySQL wait_timeout 提前关闭成脏连接的风险)
        .max_lifetime(Duration::from_secs(orm.max_life_time))
        .acquire_timeout(Duration::from_secs(10))
        .connect_timeout(Duration::from_secs(10))
        // SQL 日志级别: warn (仅 SQL 错误才输出; sea-orm 的慢查询日志默认关闭,
        // 需要排查 SQL 时可临时改回 Debug 并开启 slow_statements)
        // (对应 Go 侧 gorm 的 log_level 可配, 两侧输出量差异属框架差异)
        .sqlx_logging_level(log::LevelFilter::Warn);
    let db = Database::connect(opt).await?;
    log_elapsed!(start, elapsed_ms, info, driver = %orm.driver, "数据库连接池初始化完成");
    Ok(db)
}

/// Go DSN → sqlx URL 转换
///
/// Go 格式: root:123456@tcp(127.0.0.1:3306)/webapi?charset=utf8mb4&parseTime=True&loc=Local
/// sqlx 格式: mysql://root:123456@127.0.0.1:3306/webapi
fn dsn_to_url(driver: &str, dsn: &str) -> String {
    if dsn.contains("://") {
        return dsn.to_string();
    }
    let re = Regex::new(r"^([^@]+)@tcp\(([^)]+)\)/([^?]+)").unwrap();
    if let Some(caps) = re.captures(dsn) {
        let creds = &caps[1];
        let addr = &caps[2];
        let db = &caps[3];
        return match driver {
            "mysql" => format!("mysql://{creds}@{addr}/{db}"),
            "postgresql" | "postgres" => format!("postgres://{creds}@{addr}/{db}"),
            "sqlite" => format!("sqlite://{db}"),
            _ => dsn.to_string(),
        };
    }
    dsn.to_string()
}

/// 脱敏 DSN: 将 URL 中的账号密码替换为 ***, 用于启动日志, 避免数据库凭据落盘。
///
/// mysql://root:123456@127.0.0.1:3306/webapi → mysql://***:***@127.0.0.1:3306/webapi
fn mask_dsn(dsn: &str) -> String {
    let re = Regex::new(r"://[^@/]+@").unwrap();
    re.replace(dsn, "://***:***@").to_string()
}

#[cfg(test)]
mod tests {
    use super::mask_dsn;

    #[test]
    fn mask_dsn_hides_credentials() {
        // 含账号密码的 sqlx URL → 账号密码均脱敏, 主机/库名保留
        assert_eq!(
            mask_dsn("mysql://root:123456@127.0.0.1:3306/webapi"),
            "mysql://***:***@127.0.0.1:3306/webapi"
        );
        // 带查询参数同样只脱敏凭据部分
        assert_eq!(
            mask_dsn("mysql://root:123456@127.0.0.1:3306/webapi?charset=utf8mb4"),
            "mysql://***:***@127.0.0.1:3306/webapi?charset=utf8mb4"
        );
        // 无账号密码的 DSN (如 sqlite) 保持原样
        assert_eq!(mask_dsn("sqlite://webapi.db"), "sqlite://webapi.db");
    }
}
