"""统一业务错误: 所有 API 错误以 {code, message} 结构返回 (与各端 core.Error 对齐)"""
from fastapi import Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from starlette.exceptions import HTTPException as StarletteHTTPException

from app.core.logger import logger


class APIError(Exception):
    """业务错误: code 即 HTTP 状态码"""

    def __init__(self, status_code: int, message: str) -> None:
        self.status_code = status_code
        self.message = message
        super().__init__(message)


def bad_request(message: str) -> APIError:
    return APIError(400, message)


def param_error(message: str) -> APIError:
    """参数错误 → 422 (与 Go 端 xfiber.ParamError / Java 端 ApiException.paramError 对齐)"""
    return APIError(422, message)


def unauthorized(message: str = "Unauthorized") -> APIError:
    return APIError(401, message)


def service_unavailable(message: str = "database unavailable") -> APIError:
    return APIError(503, message)


async def api_error_handler(request: Request, exc: APIError) -> JSONResponse:
    return JSONResponse(status_code=exc.status_code, content={"code": exc.status_code, "message": exc.message})


async def validation_error_handler(request: Request, exc: RequestValidationError) -> JSONResponse:
    # 参数校验失败统一 422 (与 Go 端 bind 错误状态码对齐)
    first = exc.errors()[0] if exc.errors() else {}
    field = ".".join(str(x) for x in first.get("loc", []))
    message = first.get("msg", "validation error")
    text = f"validation error: {field}: {message}" if field else message
    return JSONResponse(status_code=422, content={"code": 422, "message": text})


async def http_error_handler(request: Request, exc: StarletteHTTPException) -> JSONResponse:
    return JSONResponse(status_code=exc.status_code, content={"code": exc.status_code, "message": str(exc.detail)})


async def unhandled_error_handler(request: Request, exc: Exception) -> JSONResponse:
    # 未知错误记录服务端日志 (含 requestId), 客户端统一返回固定 500 文案,
    # 避免将内部实现细节泄露 (与 Go 端 ErrorHandler 对齐)
    logger.error(
        "unhandled error request_id=%s path=%s error=%s",
        getattr(request.state, "request_id", ""),
        request.url.path,
        exc,
    )
    return JSONResponse(status_code=500, content={"code": 500, "message": "internal server error"})
