"""应用配置: 支持环境变量与 .env 文件覆盖 (pydantic-settings)"""
from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8", extra="ignore")

    # 应用信息
    app_name: str = "PyWebApi"
    app_version: str = "1.0.0"

    # 服务监听 (与 GoWebApi 默认端口保持一致)
    host: str = "0.0.0.0"
    port: int = 6600
    http_timeout: int = 30  # 请求超时时间(秒), 超时返回 408 (与 Go 端 http.timeout 对齐)

    # 数据库连接串 (默认 mysql, 可切换 sqlite, 如 sqlite:///./pywebapi.db)
    database_url: str = "mysql+pymysql://root:123456@127.0.0.1:3306/webapi?charset=utf8mb4"

    # JWT 配置 (与 GoWebApi/config.yaml 的 jwt 段对齐)
    jwt_secret_key: str = "6Kbj0VFeXYMp60lEyiFoVq4UzqX8Z0GSSfnvTh2VuAQn0oHgQNYexU6yYVTk4xf9"
    jwt_signing_algorithm: str = "HS256"
    jwt_expire_seconds: int = 2592000  # 过期时长(秒) = 30 天

    # 接口限流 (与 GoWebApi/config.yaml 的 limit 段对齐)
    limit_enable: bool = True
    limit_max: int = 1000  # 单个 IP 在窗口内允许的最大请求数
    limit_window: int = 60  # 滑动窗口时长(秒)


@lru_cache
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
