# 导入pydantic对应的模型基类
from pydantic import BaseModel, Field, EmailStr


class AuthLoginRequest(BaseModel):
    email: EmailStr = Field(default=..., min_length=4, description="邮箱")
    password: str = Field(default=..., min_length=6, max_length=20, description="密码")
    # phone: str = Field(default=..., regex=r'^1\d{10}$', description="手机号")
