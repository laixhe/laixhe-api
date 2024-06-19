from functools import lru_cache

from dotenv import load_dotenv
from .env import *


@lru_cache
def getAppConfig() -> AppConfigSettings:
    # 加载 .env 文件，dotenv_path 变量默认是 .env
    load_dotenv()
    # 实例化配置模型
    return AppConfigSettings()


# @lru_cache 装饰器，相当于是缓存了函数结果

# 获取配置
appSettings = getAppConfig()
