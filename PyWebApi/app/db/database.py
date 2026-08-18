"""数据库引擎与会话"""
from sqlmodel import Session, SQLModel, create_engine

from app.core.config import settings

# sqlite 多线程访问需关闭同线程检查; 其余驱动走连接池
connect_args = {"check_same_thread": False} if settings.database_url.startswith("sqlite") else {}
engine = create_engine(settings.database_url, connect_args=connect_args, pool_pre_ping=True)


def init_db() -> None:
    """启动时建表 (教学规模直接 create_all; 生产建议引入 alembic 迁移)"""
    from app.models import config_common, user, user_extend, user_third_party  # noqa: F401  确保模型注册到 metadata

    SQLModel.metadata.create_all(engine)


def get_session():
    with Session(engine) as session:
        yield session
