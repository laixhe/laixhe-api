import os

# =========== 数据库配置 ===========

# 数据库驱动
DB_DRIVER = os.getenv('DB_DRIVER', 'mysql')
# 数据库地址
DB_HOST = os.getenv('DB_HOST', '127.0.0.1')
# 数据库端口
DB_PORT = os.getenv('DB_PORT', 3306)
# 数据库名称
DB_DATABASE = os.getenv('DB_DATABASE', 'webapi')
# 数据库账号
DB_USERNAME = os.getenv('DB_USERNAME', 'root')
# 数据库密码
DB_PASSWORD = os.getenv('DB_PASSWORD', '')
# 数据表前缀
DB_PREFIX = os.getenv('DB_PREFIX', '')
# 是否开启调试模式：是-True,否-False
DB_DEBUG = (os.getenv('DB_DEBUG', 'True') == 'True')
# MySQL数据库链接(当前使用的数据库)
SQLALCHEMY_MYSQL_URL = 'mysql+pymysql://' + DB_USERNAME + ':' + DB_PASSWORD + '@' + DB_HOST + ':' + str(
    DB_PORT) + '/' + DB_DATABASE + '?charset=utf8mb4'
