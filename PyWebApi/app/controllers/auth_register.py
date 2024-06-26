from fastapi import APIRouter
# 导入pydantic对应的模型基类
from pydantic import BaseModel, Field, EmailStr

from app.result import result

router = APIRouter(
    prefix="/auth",
    tags=["鉴权相关"]
)


class AuthRegisterRequest(BaseModel):
    email: EmailStr = Field(default=..., min_length=4, description="邮箱")
    password: str = Field(default=..., min_length=6, max_length=20, description="密码")
    uname: str = Field(default=..., min_length=2, max_length=30, description="用户名")
    age: int = Field(default=..., ge=0, le=200, description="年龄")


@router.post("/register")
async def auth_register(req: AuthRegisterRequest) -> result.Result:
    """
    注册
    """

    return result.success(req)
