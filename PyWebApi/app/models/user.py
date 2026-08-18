"""用户模型 (对齐 webapi.sql 的 user 表)"""
from datetime import datetime

from sqlmodel import Field, SQLModel


class UserType:
    """用户类型"""
    ORDINARY = 1  # 普通用户


class UserSex:
    """性别"""
    UNKNOWN = 0  # 未填写
    MALE = 1  # 男
    FEMALE = 2  # 女


class UserState:
    """账号状态"""
    DISABLED = 0  # 禁用
    NORMAL = 1  # 正常


class User(SQLModel, table=True):
    """用户表"""
    __tablename__ = "user"

    id: int | None = Field(default=None, primary_key=True, description="用户id")
    type_id: int = Field(default=0, description="类型 1普通 (对齐 webapi.sql default 0)")
    account: str = Field(default="", max_length=120, unique=True, index=True, description="账号")
    mobile: str = Field(default="", max_length=120, index=True, description="手机号")
    email: str = Field(default="", max_length=120, unique=True, index=True, description="邮箱")
    password: str = Field(default="", max_length=120, description="密码(bcrypt)")
    nickname: str = Field(default="", max_length=120, description="昵称")
    avatar_url: str = Field(default="", max_length=255, description="头像地址")
    sex: int = Field(default=UserSex.UNKNOWN, description="性别 0未填写 1男 2女")
    states: int = Field(default=0, description="状态 0封禁 1正常 (对齐 webapi.sql default 0)")
    created_at: datetime = Field(default_factory=datetime.now, description="创建时间")
    updated_at: datetime = Field(default_factory=datetime.now, description="更新时间")
