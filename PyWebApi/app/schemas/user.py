"""用户相关请求/响应"""
from typing import List

from pydantic import BaseModel, Field

from app.models.user import User as UserModel

# 时间格式: 与 Go 端 time.DateTime (2006-01-02 15:04:05) 保持一致
DATETIME_FORMAT = "%Y-%m-%d %H:%M:%S"


class User(BaseModel):
    """用户信息"""
    uid: int = Field(description="用户id")
    type_id: int = Field(description="类型 (UserType: * 1 - 普通用户)")
    account: str = Field(description="账号")
    mobile: str = Field(description="手机号")
    email: str = Field(description="邮箱")
    nickname: str = Field(description="昵称")
    avatar_url: str = Field(description="头像地址")
    sex: int = Field(description="性别 (UserSex: * 0 - 未填写 * 1 - 男 * 2 - 女)")
    states: int = Field(description="状态 (UserState: * 0 - 禁用 * 1 - 正常)")
    created_at: str = Field(description="创建时间")

    @classmethod
    def from_model(cls, m: UserModel, override_nickname: str = "", override_avatar_url: str = "") -> "User":
        return cls(
            uid=m.id,
            type_id=m.type_id,
            account=m.account,
            mobile=m.mobile,
            email=m.email,
            nickname=override_nickname or m.nickname,
            avatar_url=override_avatar_url or m.avatar_url,
            sex=m.sex,
            states=m.states,
            created_at=m.created_at.strftime(DATETIME_FORMAT),
        )


class UserUpdateRequest(BaseModel):
    """请求-更新用户信息 (Uid 由 JWT 提供)"""
    nickname: str = Field(description="昵称")
    avatar_url: str = Field(default="", description="头像地址")


class UserListResponse(BaseModel):
    """响应-获取用户列表"""
    total: int = Field(description="总数")
    page: int = Field(description="分页-当前页")
    page_size: int = Field(description="分页-每页数量")
    list: List[User] = Field(description="列表")
