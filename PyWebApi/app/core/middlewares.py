"""通用 HTTP 中间件: requestId / 访问日志 / 请求超时 (与 Go 端 xfiber 中间件对齐)"""
import asyncio
import time
import uuid

from fastapi import Request
from fastapi.responses import JSONResponse

from app.core.config import settings
from app.core.logger import logger
from app.core.rate_limit import resolve_client_ip

# 请求唯一标识响应头 (与 Go 端 fiber.HeaderXRequestID 一致)
REQUEST_ID_HEADER = "X-Request-Id"


async def request_id_middleware(request: Request, call_next):
    """请求唯一标识: 无则生成, 响应头透出 X-Request-Id (与 Go 端 xfiber 内置 requestId 对齐)"""
    request_id = request.headers.get(REQUEST_ID_HEADER) or uuid.uuid4().hex
    request.state.request_id = request_id
    response = await call_next(request)
    response.headers[REQUEST_ID_HEADER] = request_id
    return response


async def access_log_middleware(request: Request, call_next):
    """访问日志: 记录 method/path/状态码/耗时/客户端IP/requestId (与 Go 端 UseLog 对齐)"""
    start = time.perf_counter()
    response = await call_next(request)
    cost_ms = (time.perf_counter() - start) * 1000
    logger.info(
        "access method=%s path=%s status=%d cost=%.1fms ip=%s request_id=%s",
        request.method,
        request.url.path,
        response.status_code,
        cost_ms,
        resolve_client_ip(request),
        getattr(request.state, "request_id", ""),
    )
    return response


async def timeout_middleware(request: Request, call_next):
    """请求超时: 超过 http.timeout 秒返回 408 统一 JSON (与 Go 端 timeout 中间件对齐)"""
    try:
        return await asyncio.wait_for(call_next(request), timeout=settings.http_timeout)
    except asyncio.TimeoutError:
        return JSONResponse(status_code=408, content={"code": 408, "message": "Request Timeout"})
