"""鉴权相关请求/响应"""
from pydantic import BaseModel, Field

from app.schemas.user import User


class AuthRegisterRequest(BaseModel):
    """请求-注册"""
    email: str = Field(description="邮箱")
    nickname: str = Field(description="昵称")
    password: str = Field(description="密码")


class AuthLoginRequest(BaseModel):
    """请求-登录"""
    email: str = Field(description="邮箱")
    password: str = Field(description="密码")


class AuthResponse(BaseModel):
    """响应-登录/注册/刷新"""
    token: str = Field(description="jwt token")
    user: User
