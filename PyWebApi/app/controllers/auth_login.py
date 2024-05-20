from fastapi import APIRouter
# 导入pydantic对应的模型基类
from pydantic import BaseModel, Field, EmailStr

from app.types import response

router = APIRouter(
    prefix="/auth",
    tags=["鉴权相关"]
)


class AuthLoginRequest(BaseModel):
    email: EmailStr = Field(default='', min_length=4, description="邮箱")
    password: str = Field(default='', min_length=6, description="密码")
    # phone: str = Field(default='', regex=r'^1\d{10}$', description="手机号")


@router.post("/login")
async def auth_login(req: AuthLoginRequest) -> response.HttpResponse:
    """
    登录
    """

    return response.response_success(req)
