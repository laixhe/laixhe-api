//! 同步原语辅助 (提取重复的锁中毒恢复逻辑)
//!
//! 项目多处使用 `lock().unwrap_or_else(|poisoned| poisoned.into_inner())`
//! 处理锁中毒后继续使用数据 (限流计数/日志 guard 等场景数据一致性可接受),
//! 统一收敛到 [`lock_unpoison`] 避免重复。

use std::sync::{Mutex, MutexGuard};

/// 获取 Mutex 锁; 锁中毒时取回内部数据继续使用
///
/// 仅适用于"锁内数据一致性可接受"的场景 (计数、缓存等),
/// 不适用于对一致性要求严格的数据结构。
pub fn lock_unpoison<T>(m: &Mutex<T>) -> MutexGuard<'_, T> {
    m.lock().unwrap_or_else(|poisoned| poisoned.into_inner())
}
