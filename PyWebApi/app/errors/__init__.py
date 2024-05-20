from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException as StarletteHTTPException

from .exception_handler import *


def register_custom_error_handle(server: FastAPI):
    """ 统一注册自定义错误处理器 """
    # 注册参数验证错误,并覆盖模式 RequestValidationError
    server.add_exception_handler(RequestValidationError, validation_exception_handler)
    # 错误处理StarletteHTTPException
    server.add_exception_handler(StarletteHTTPException, http_exception_handler)
    # 自定义全局系统错误
    server.add_exception_handler(Exception, app_exception_handler)
