// 应用配置
// 从环境变量读取，提供合理默认值，便于本地开发和容器化部署

// 环境变量校验
const missing: string[] = [];
if (!process.env.DATABASE_URL) missing.push("DATABASE_URL");
if (missing.length > 0) {
  console.error(`Missing required environment variables: ${missing.join(", ")}`);
  process.exit(1);
}

// 将环境变量解析为正整数; 非法值 (非数字/非正数) 回落默认值并给出提示,
// 避免 parseInt 静默产生 NaN 导致运行时晦涩错误 (如 HTTP_PORT=abc → NaN → 监听失败)
function parsePositiveInt(value: string | undefined, fallback: number, name: string): number {
  const n = parseInt(value ?? "", 10);
  if (Number.isNaN(n) || n <= 0) {
    console.warn(`[config] 非法的 ${name} "${value}", 已回落为默认值 ${fallback}`);
    return fallback;
  }
  return n;
}

const config = {
  // HTTP 服务配置
  http: {
    ip: process.env.HTTP_IP || "0.0.0.0",        // 监听地址（默认所有网卡）
    port: parsePositiveInt(process.env.HTTP_PORT, 6600, "HTTP_PORT"), // 监听端口（默认 6600）
    // 请求超时时间（秒），预留配置以对齐 Go/Rust 端 http.timeout: 30；
    // 当前框架未实现请求超时中间件，该值暂未生效（实现后可从此处读取）
    timeout: parsePositiveInt(process.env.HTTP_TIMEOUT, 30, "HTTP_TIMEOUT"),
  },
  // JWT Token 配置
  jwt: {
    // 签名密钥（需转为 Uint8Array 供 jose 使用）
    secretKey: new TextEncoder().encode(
      process.env.JWT_SECRET_KEY || "default-secret-key-change-me"
    ),
    expireTime: parsePositiveInt(process.env.JWT_EXPIRE_TIME, 2592000, "JWT_EXPIRE_TIME"), // 过期时间（秒），默认 30 天
  },
  // 日志配置
  log: {
    level: process.env.LOG_LEVEL || "debug", // 日志级别：debug | info | warn | error（非法值由 logger.ts 回落 debug）
  },
};

export default config;
