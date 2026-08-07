// 应用配置
// 从环境变量读取，提供合理默认值，便于本地开发和容器化部署

// 环境变量校验
const missing: string[] = [];
if (!process.env.DATABASE_URL) missing.push("DATABASE_URL");
if (missing.length > 0) {
  console.error(`Missing required environment variables: ${missing.join(", ")}`);
  process.exit(1);
}

const config = {
  // HTTP 服务配置
  http: {
    ip: process.env.HTTP_IP || "0.0.0.0",        // 监听地址（默认所有网卡）
    port: parseInt(process.env.HTTP_PORT || "6600", 10), // 监听端口（默认 6600）
  },
  // JWT Token 配置
  jwt: {
    // 签名密钥（需转为 Uint8Array 供 jose 使用）
    secretKey: new TextEncoder().encode(
      process.env.JWT_SECRET_KEY || "default-secret-key-change-me"
    ),
    expireTime: parseInt(process.env.JWT_EXPIRE_TIME || "2592000", 10), // 过期时间（秒），默认 30 天
  },
  // 日志配置
  log: {
    level: process.env.LOG_LEVEL || "debug", // 日志级别：debug | info | warn | error
  },
};

export default config;
