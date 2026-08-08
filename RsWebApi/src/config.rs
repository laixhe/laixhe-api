//! 配置 (对应 Go 项目 core/config.go)
//!
//! 使用 serde_yaml 加载 config.yaml，并支持 `${ENV}` / `$ENV` 环境变量展开。

use serde::Deserialize;
use std::sync::Arc;

/// HTTP 服务监听地址
#[derive(Debug, Clone, Deserialize)]
pub struct AddrConfig {
    pub ip: String,
    pub port: i32,
    /// 请求超时时间(单位秒)，缺省 30 秒，用于请求超时中间件
    #[serde(default)]
    pub timeout: Option<i32>,
}

impl AddrConfig {
    /// 返回 "ip:port" 格式的监听地址
    pub fn addr(&self) -> String {
        format!("{}:{}", self.ip, self.port)
    }
}

/// 日志配置
#[derive(Debug, Clone, Deserialize)]
pub struct LogConfig {
    /// 日志模式 console / file
    pub run: String,
    /// 日志输出格式 console(文本) / json(结构化, 便于采集)
    #[serde(default = "default_log_format")]
    pub format: String,
    /// 日志文件路径
    pub path: String,
    /// 日志级别 debug info warn error
    pub level: String,
    /// 每个日志文件保存大小(单位 MB); config.yaml 示例 20, 缺省 3 (对齐 Go xlog 缺省值)
    #[serde(default = "default_max_size")]
    pub max_size: i64,
    /// 保留 N 个备份 (缺省 3, 对齐 Go xlog 缺省值; 对应 rolling-file 的 max_files)
    #[serde(default = "default_max_backups")]
    pub max_backups: i64,
    /// 保留 N 天 (缺省 3, 对齐 Go xlog 缺省值;
    /// 注意: rolling-file 无按天清理能力, 该字段仅保留以对齐配置项, 实际清理由 max_backups 控制)
    #[allow(dead_code)]
    #[serde(default = "default_max_age")]
    pub max_age: i64,
}

fn default_log_format() -> String {
    "console".to_string()
}

fn default_max_size() -> i64 {
    3
}
fn default_max_backups() -> i64 {
    3
}
fn default_max_age() -> i64 {
    3
}

impl Default for LogConfig {
    fn default() -> Self {
        LogConfig {
            run: "console".to_string(),
            format: default_log_format(),
            path: "logs.log".to_string(),
            level: "debug".to_string(),
            max_size: default_max_size(),
            max_backups: default_max_backups(),
            max_age: default_max_age(),
        }
    }
}

/// ORM 配置
#[derive(Debug, Clone, Deserialize)]
pub struct OrmConfig {
    /// 驱动名称: mysql (当前仅编译 MySQL 后端; postgresql / sqlite 需在 Cargo.toml 追加对应 sqlx-* feature)
    pub driver: String,
    /// 连接地址
    pub dsn: String,
    /// 设置空闲连接池中连接的最大数量
    pub max_idle_count: u32,
    /// 设置打开数据库连接的最大数量
    pub max_open_count: u32,
    /// 设置了连接可复用的最大时间(单位秒)
    pub max_life_time: u64,
}

/// JWT 配置
#[derive(Debug, Clone, Deserialize)]
pub struct JwtConfig {
    /// 密钥
    pub secret_key: String,
    /// 过期时长(单位秒)
    pub expire_time: i64,
}

/// 接口限流配置
///
/// 字段级缺省为零值 (enable=false 即限流关闭), 对齐 Go mapstructure 语义;
/// 整段缺失时由 Config::load 填充默认 {true, 1000, 60} (对齐 Go config.Check)
#[derive(Debug, Clone, Deserialize, Default)]
#[serde(default)]
pub struct LimitConfig {
    /// 是否启用限流
    pub enable: bool,
    /// 单个 IP 在窗口内允许的最大请求数
    pub max: usize,
    /// 滑动窗口时长 (单位秒)
    pub window: u64,
}

/// 运行时通用配置，从数据库 config_common 表动态加载
#[derive(Debug, Clone, Default)]
pub struct CommonConfig {
    pub env: String,
}

/// 总配置
#[derive(Debug, Clone)]
pub struct Config {
    pub http: Arc<AddrConfig>,
    pub log: LogConfig,
    pub orm: OrmConfig,
    pub jwt: JwtConfig,
    pub limit: LimitConfig,
}

/// 磁盘上的原始配置结构
#[derive(Debug, Deserialize)]
struct RawConfig {
    http: Option<AddrConfig>,
    #[serde(default)]
    log: Option<LogConfig>,
    orm: Option<OrmConfig>,
    jwt: Option<JwtConfig>,
    #[serde(default)]
    limit: Option<LimitConfig>,
}

impl Config {
    /// 加载配置文件并校验
    pub fn load(config_file: &str) -> Result<Config, String> {
        let content = std::fs::read_to_string(config_file)
            .map_err(|e| format!("read config file failed: {e}"))?;
        // 展开环境变量 ${VAR} / $VAR
        let content = expand_env(&content);
        let raw: RawConfig =
            serde_yaml::from_str(&content).map_err(|e| format!("parse config file failed: {e}"))?;
        let config = Config {
            http: raw
                .http
                .map(Arc::new)
                .ok_or_else(|| "http config is nil".to_string())?,
            log: raw.log.unwrap_or_default(),
            orm: raw.orm.ok_or_else(|| "orm config is nil".to_string())?,
            jwt: raw.jwt.ok_or_else(|| "jwt config is nil".to_string())?,
            // limit 整段缺失时补默认 {enable, 1000, 60} (对齐 Go config.Check)
            limit: raw.limit.unwrap_or(LimitConfig {
                enable: true,
                max: 1000,
                window: 60,
            }),
        };
        config.check()?;
        Ok(config)
    }

    /// 校验配置有效性
    pub fn check(&self) -> Result<(), String> {
        if self.http.port <= 0 {
            return Err("http port is invalid".to_string());
        }
        // 驱动白名单校验 (对齐 Go gonet/orm 的 driver 校验)
        if !["mysql", "postgresql", "sqlite"].contains(&self.orm.driver.as_str()) {
            return Err(format!("orm driver is invalid: {}", self.orm.driver));
        }
        if self.orm.dsn.is_empty() {
            return Err("orm dsn is empty".to_string());
        }
        if self.jwt.secret_key.is_empty() {
            return Err("jwt secret_key is empty".to_string());
        }
        // 过期时长必须为正数, 否则会签发立即过期的 token (对齐 Go gonet/jwt 对 expire_time <= 0 的校验)
        if self.jwt.expire_time <= 0 {
            return Err("jwt expire_time is invalid".to_string());
        }
        Ok(())
    }
}

/// 展开配置字符串中的环境变量 ${VAR} / $VAR
fn expand_env(content: &str) -> String {
    let bytes = content.as_bytes();
    let mut result = String::with_capacity(content.len());
    let mut i = 0;
    while i < bytes.len() {
        if bytes[i] == b'$' && i + 1 < bytes.len() {
            if bytes[i + 1] == b'{' {
                // ${VAR} 形式
                if let Some(end) = content[i + 2..].find('}') {
                    let key = &content[i + 2..i + 2 + end];
                    let value = std::env::var(key).unwrap_or_default();
                    result.push_str(&value);
                    i += 2 + end + 1;
                    continue;
                }
            } else if bytes[i + 1].is_ascii_alphabetic() || bytes[i + 1] == b'_' {
                // $VAR 形式
                let mut end = i + 1;
                while end < bytes.len()
                    && (bytes[end].is_ascii_alphanumeric() || bytes[end] == b'_')
                {
                    end += 1;
                }
                let key = &content[i + 1..end];
                let value = std::env::var(key).unwrap_or_default();
                result.push_str(&value);
                i = end;
                continue;
            }
        }
        // 普通字符
        let ch = content[i..].chars().next().unwrap();
        result.push(ch);
        i += ch.len_utf8();
    }
    result
}
