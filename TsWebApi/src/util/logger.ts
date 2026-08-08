// 简易日志工具，按级别输出带时间戳的结构化日志

type LogLevel = "debug" | "info" | "warn" | "error";

const levelPriority: Record<LogLevel, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
};

const currentLevel: LogLevel =
  (process.env.LOG_LEVEL as LogLevel) || "debug";

const currentLevelPriority = levelPriority[currentLevel];

// 预计算各级别是否启用，避免每次日志调用都执行比较
const debugEnabled = levelPriority["debug"] >= currentLevelPriority;
const infoEnabled = levelPriority["info"] >= currentLevelPriority;
const warnEnabled = levelPriority["warn"] >= currentLevelPriority;
const errorEnabled = levelPriority["error"] >= currentLevelPriority;

function timestamp(): string {
  return new Date().toISOString();
}

// 安全序列化日志数据：循环引用等异常数据降级为 String 输出，避免日志本身抛错
function toJson(data: unknown): string {
  try {
    return JSON.stringify(data);
  } catch {
    return String(data);
  }
}

export function debug(module: string, message: string, data?: unknown) {
  if (debugEnabled) {
    const extra = data !== undefined ? ` | ${toJson(data)}` : "";
    console.debug(`[${timestamp()}] [DEBUG] [${module}] ${message}${extra}`);
  }
}

export function info(module: string, message: string, data?: unknown) {
  if (infoEnabled) {
    const extra = data !== undefined ? ` | ${toJson(data)}` : "";
    console.info(`[${timestamp()}] [INFO] [${module}] ${message}${extra}`);
  }
}

export function warn(module: string, message: string, data?: unknown) {
  if (warnEnabled) {
    const extra = data !== undefined ? ` | ${toJson(data)}` : "";
    console.warn(`[${timestamp()}] [WARN] [${module}] ${message}${extra}`);
  }
}

export function error(module: string, message: string, data?: unknown) {
  if (errorEnabled) {
    const extra = data !== undefined ? ` | ${toJson(data)}` : "";
    console.error(`[${timestamp()}] [ERROR] [${module}] ${message}${extra}`);
  }
}
