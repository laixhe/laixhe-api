import Redis from "ioredis";

const redisConfig = {
    port: parseInt(process.env.REDIS_PORT || "6379"),
    host: process.env.REDIS_HOST || "localhost",
    username: process.env.REDIS_USERNAME || undefined,
    password: process.env.REDIS_PASSWORD || undefined,
    db: parseInt(process.env.REDIS_DB || "0"),
};

const cache = new Redis(redisConfig.port, redisConfig.host, {
    username: redisConfig.username,
    password: redisConfig.password,
    db: redisConfig.db,
});

export default cache;
