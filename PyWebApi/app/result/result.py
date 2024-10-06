from typing import Any
from pydantic import BaseModel, Field


class Result(BaseModel):
    """ http统一响应 """
    code: int = Field(default=0)  # 响应码
    msg: str = Field(default='')  # 响应信息
    data: Any | None  # 具体数据


def success(resp: Any) -> Result:
    """ 成功响应 """
    return Result(data=resp)


def fail(msg: str, code: int = 1) -> Result:
    """ 响应失败 """
    return Result(
        code=code,
        msg=msg)
