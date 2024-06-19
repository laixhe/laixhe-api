from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from contextlib import contextmanager

from config import appSettings


# 创建引擎
engine = create_engine(
    appSettings.db_dsn,
    echo=appSettings.db_echo_sql,
    pool_size=appSettings.db_pool_size,
    max_overflow=appSettings.db_max_overflow,
)

# 封装获取会话
Session = sessionmaker(bind=engine, expire_on_commit=False)


@contextmanager
def getDatabaseSession(autoCommitByExit=True):
    """使用上下文管理资源关闭"""
    _session = Session()
    try:
        yield _session
        # 退出时，是否自动提交
        if autoCommitByExit:
            _session.commit()
    except Exception as e:
        _session.rollback()
        raise e
