from typing import Any
from pydantic import BaseModel, Field


class HttpResponse(BaseModel):
    """ http统一响应 """
    code: int = Field(default=0)  # 响应码
    msg: str = Field(default='')  # 响应信息
    data: Any | None  # 具体数据


def response_success(resp: Any) -> HttpResponse:
    """ 成功响应 """
    return HttpResponse(data=resp)


def response_fail(msg: str, code: int = 1) -> HttpResponse:
    """ 响应失败 """
    return HttpResponse(
        code=code,
        msg=msg)
