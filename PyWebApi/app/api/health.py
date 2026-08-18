"""健康检查接口"""
import threading
import time
from datetime import datetime

from fastapi import APIRouter
from pydantic import BaseModel, Field
from sqlalchemy import text

from app.core.config import settings
from app.core.errors import service_unavailable
from app.db.database import engine
from app.schemas.common import Error

router = APIRouter(prefix="/api/v1", tags=["Health"])

DATETIME_FORMAT = "%Y-%m-%d %H:%M:%S"

# 服务启动时间 (进程级: 启动时赋值一次, 之后只读, 服务器本地时区)
_started_at = datetime.now().strftime(DATETIME_FORMAT)

# 数据库探测结果缓存时长(秒), 避免频繁探活压垮数据库 (与 Go healthPingInterval / Java PING_INTERVAL_MS 对齐)
_HEALTH_PING_INTERVAL = 5.0
_ping_lock = threading.Lock()
_last_ping = 0.0
_last_healthy = True


class HealthResponse(BaseModel):
    """健康检查响应体"""
    status: str = Field(description='服务状态 (固定 "ok")')
    database: str = Field(description='数据库状态 (正常时为 "up"; 数据库不可用时直接返回 503 错误体, 不返回本字段)')
    version: str = Field(description="服务版本")
    started_at: str = Field(description="服务启动时间 (服务器本地时区)")
    now: str = Field(description="当前时间 (服务器本地时区)")


def _db_healthy() -> bool:
    """探测数据库连接, 结果缓存 _HEALTH_PING_INTERVAL 秒。

    读路径 (缓存有效) 无锁并发读, 互不阻塞; 缓存过期时由首个请求加锁探测 (double-check 防并发重复)。
    """
    global _last_ping, _last_healthy
    now = time.monotonic()
    if now - _last_ping < _HEALTH_PING_INTERVAL:
        return _last_healthy
    with _ping_lock:
        # double-check: 并发第一个抢到锁的请求已刷新缓存, 后续请求直接复用
        if now - _last_ping < _HEALTH_PING_INTERVAL:
            return _last_healthy
        ok = True
        try:
            with engine.connect() as conn:
                conn.execute(text("SELECT 1"))
        except Exception:
            ok = False
        _last_ping = time.monotonic()
        _last_healthy = ok
        return ok


@router.get(
    "/health",
    response_model=HealthResponse,
    summary="健康检查",
    responses={503: {"model": Error}},
)
def health() -> HealthResponse:
    # 探测数据库连接; 不可用时返回 503 统一错误体, 便于负载均衡探活
    if not _db_healthy():
        raise service_unavailable()
    return HealthResponse(
        status="ok",
        database="up",
        version=settings.app_version,
        started_at=_started_at,
        now=datetime.now().strftime(DATETIME_FORMAT),
    )
