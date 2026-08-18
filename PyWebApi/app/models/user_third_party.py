"""用户第三方模型 (对齐 webapi.sql 的 user_third_party 表)"""
from sqlmodel import Field, SQLModel


class UserThirdParty(SQLModel, table=True):
    """用户第三方表"""
    __tablename__ = "user_third_party"

    id: int | None = Field(default=None, primary_key=True)
    uid: int = Field(default=0, unique=True, index=True, description="用户ID")
    wechat_unionid: str = Field(default="", max_length=200, description="微信unionid")
    wechat_openid: str = Field(default="", max_length=200, index=True, description="微信openid")
