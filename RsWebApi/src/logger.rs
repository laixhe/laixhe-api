//! 统一日志模块
//!
//! 集中封装耗时统计日志，统一结构化字段格式（`elapsed_ms` / `total_ms`），
//! 供 controller / service / middleware / state 等模块复用。
//!
//! 使用方式：
//! ```ignore
//! use crate::logger::Timer;
//! use crate::log_elapsed;
//!
//! let step = Timer::new();
//! // ... 业务操作 ...
//! log_elapsed!(step, elapsed_ms, info, email = %req.email, "检查通过");
//! log_elapsed!(start, total_ms, info, email = %req.email, uid, "注册完成");
//! ```

use std::time::Instant;

/// 步骤计时器：记录一段操作的耗时
#[derive(Debug)]
pub struct Timer {
    start: Instant,
}

impl Timer {
    /// 创建并开始计时
    pub fn new() -> Self {
        Self {
            start: Instant::now(),
        }
    }

    /// 当前累计耗时（毫秒）
    pub fn elapsed_ms(&self) -> u64 {
        self.start.elapsed().as_millis() as u64
    }
}

impl Default for Timer {
    fn default() -> Self {
        Self::new()
    }
}

/// 输出耗时日志：在结构化字段前自动追加耗时字段
///
/// 参数说明：
/// - `$timer`：[`Timer`] 计时器
/// - `$field`：耗时字段名（`elapsed_ms` / `total_ms`）
/// - `$level`：日志级别宏名（`info` / `warn` / `debug` / `error`）
/// - `$rest`：与 `tracing::info!` 相同的结构化字段与消息
///
/// 示例：
/// ```ignore
/// let step = logger::Timer::new();
/// log_elapsed!(step, elapsed_ms, info, email = %req.email, "邮箱检查通过");
/// log_elapsed!(start, total_ms, info, email = %req.email, uid, "注册完成");
/// log_elapsed!(start, elapsed_ms, warn, method = %method, path = %path, "缺少 Token");
/// ```
#[macro_export]
macro_rules! log_elapsed {
    ($timer:expr, $field:ident, $level:ident, $($rest:tt)+) => {
        ::tracing::$level!($field = $timer.elapsed_ms(), $($rest)+);
    };
}
