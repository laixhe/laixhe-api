"""IP 限流中间件 (与 Go 端 RateLimit 中间件行为对齐): 滑动窗口, 默认 1000 次/60s/IP"""
import threading
import time
from collections import defaultdict, deque

from fastapi import Request, Response
from fastapi.responses import JSONResponse

from app.core.config import settings

# 健康检查路径限流豁免 (负载均衡/容器编排探活不应被 429 拦截)
RATE_LIMIT_HEALTH_PATH = "/api/v1/health"


class RateLimiter:
    """滑动窗口限流器: key(IP) → 窗口内请求时间戳队列 (FIFO)"""

    def __init__(self, max_requests: int, window: float, max_keys: int = 100_000) -> None:
        self.max_requests = max_requests
        self.window = window
        self.max_keys = max_keys
        self._records: dict[str, deque[float]] = defaultdict(deque)
        self._lock = threading.Lock()

    def allow(self, key: str) -> bool:
        """是否允许本次请求 (窗口外过期记录队首清理, 平摊 O(1)); 返回 False 表示超限"""
        now = time.monotonic()
        with self._lock:
            records = self._records[key]
            while records and now - records[0] >= self.window:
                records.popleft()
            if len(records) >= self.max_requests:
                return False
            records.append(now)
            # 内存保护: key 数超限时清理已无活动窗口的 key (过期 key 的回收见下)
            if len(self._records) > self.max_keys:
                for k, q in list(self._records.items()):
                    if not q or now - q[-1] >= self.window:
                        del self._records[k]
            return True


_rate_limiter = RateLimiter(settings.limit_max, settings.limit_window)


def resolve_client_ip(request: Request) -> str:
    """解析客户端 IP: 优先 X-Forwarded-For 第一个, 其次真实连接地址, 兜底 "unknown" (与 Go 端一致)"""
    xff = request.headers.get("x-forwarded-for", "").strip()
    if xff:
        first = xff.split(",")[0].strip()
        if first:
            return first
    if request.client and request.client.host:
        return request.client.host
    return "unknown"


async def rate_limit_middleware(request: Request, call_next) -> Response:
    """IP 限流中间件: 超限返回 429 统一 JSON; 健康检查路径豁免"""
    if not settings.limit_enable or request.url.path == RATE_LIMIT_HEALTH_PATH:
        return await call_next(request)
    if not _rate_limiter.allow(resolve_client_ip(request)):
        return JSONResponse(
            status_code=429,
            content={"code": 429, "message": "请求过于频繁，请稍后再试"},
        )
    return await call_next(request)
