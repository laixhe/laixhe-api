from typing import Any

from fastapi.exceptions import RequestValidationError
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse
from fastapi import Request, status
from starlette.exceptions import HTTPException


class HttpResponse:
    code: int
    msg: str
    data: Any

    def __init__(self, code: int, msg: str):
        self.code = code
        self.msg = msg
        self.data = None


def response_fail(msg: str, code: int = 1) -> HttpResponse:
    """ 响应失败 """
    return HttpResponse(
        code=code,
        msg=msg)


async def validation_exception_handler(request: Request, exc: RequestValidationError):
    """ 自定义参数验证异常错误 """

    msg = ""
    for error in exc.errors():
        msg += ".".join(error.get("loc")) + ":" + error.get("msg") + ";"

    return JSONResponse(status_code=status.HTTP_200_OK, content=jsonable_encoder(response_fail(msg)))


async def http_exception_handler(request, exc: HTTPException) -> JSONResponse:
    """ 自定义处理HTTPException """
    print("request:", request)
    print("status_code:", exc.status_code)
    if exc.status_code == status.HTTP_404_NOT_FOUND:
        # 处理404错误
        return JSONResponse(
            content=jsonable_encoder(response_fail("接口路由不存在")),
            status_code=status.HTTP_200_OK,
        )
    elif exc.status_code == status.HTTP_405_METHOD_NOT_ALLOWED:
        # 处理405错误
        return JSONResponse(
            content=jsonable_encoder(response_fail("请求方式错误")),
            status_code=status.HTTP_200_OK,
        )
    else:
        return JSONResponse(
            content=jsonable_encoder(response_fail(str(exc))),
            status_code=status.HTTP_200_OK,
        )


async def app_exception_handler(request: Request, exc: Exception):
    """自定义全局系统错误"""
    return JSONResponse(
        content=jsonable_encoder(response_fail("系统运行异常,稍后重试~")),
        status_code=status.HTTP_200_OK,
    )
