"""通用响应结构"""
from pydantic import BaseModel, Field


class Error(BaseModel):
    """统一错误响应 {code, message}"""
    code: int = Field(description="错误码 (即 HTTP 状态码)")
    message: str = Field(description="错误信息")
