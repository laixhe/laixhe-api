from pydantic_settings import BaseSettings


class AppConfigSettings(BaseSettings):

    """应用配置"""

    app_name: str = ""
    app_host: str = "127.0.0.1"
    app_port: int = 8080
    app_debug: bool = False

    """数据库配置"""

    # 数据库连接
    db_dsn: str = ""
    # 使用打印SQL日志信息
    db_echo_sql: bool = True
    # 连接池中的初始连接数，默认为 5
    db_pool_size: int = 5
    # 连接池中允许的最大超出连接数，默认为 10
    db_max_overflow: int = 10

    """jwt配置"""

    # 秘钥
    jwt_secret: str = ""
    # 过期时间(秒)
    jwt_expire: int = 604800

