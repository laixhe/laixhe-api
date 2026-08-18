"""应用入口"""
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException as StarletteHTTPException

from app.api.router import api_router
from app.core.config import settings
from app.core.errors import (
    APIError,
    api_error_handler,
    http_error_handler,
    unhandled_error_handler,
    validation_error_handler,
)
from app.core.middlewares import access_log_middleware, request_id_middleware, timeout_middleware
from app.core.rate_limit import rate_limit_middleware
from app.db.database import init_db
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_db()
    yield


app = FastAPI(
    title="API接口",
    description="用户认证与用户管理 API 服务",
    version="1.0.0",
    lifespan=lifespan,
)

# 中间件执行顺序 (与 Go 端一致): requestId → 访问日志 → CORS → gzip → 请求超时(408) → IP 限流(429) → 业务路由
# Starlette 后注册的中间件在最外层, 因此按相反顺序注册
app.middleware("http")(rate_limit_middleware)  # 最内层: IP 限流 (健康检查豁免)
app.middleware("http")(timeout_middleware)  # 请求超时 (408)
app.add_middleware(GZipMiddleware, minimum_size=1000)  # gzip 压缩
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)  # CORS
app.middleware("http")(access_log_middleware)  # 访问日志
app.middleware("http")(request_id_middleware)  # 最外层: requestId

app.add_exception_handler(APIError, api_error_handler)
app.add_exception_handler(RequestValidationError, validation_error_handler)
app.add_exception_handler(StarletteHTTPException, http_error_handler)
app.add_exception_handler(Exception, unhandled_error_handler)

app.include_router(api_router)
