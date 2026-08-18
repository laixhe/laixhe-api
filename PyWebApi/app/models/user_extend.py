"""用户扩展模型 (对齐 webapi.sql 的 user_extend 表)"""
from sqlmodel import Field, SQLModel


class UserExtend(SQLModel, table=True):
    """用户扩展表"""
    __tablename__ = "user_extend"

    id: int | None = Field(default=None, primary_key=True)
    uid: int = Field(default=0, unique=True, index=True, description="用户ID")
    birthday: int = Field(default=0, description="生日(年月日)")
    height: int = Field(default=0, description="身高(cm)")
    weight: int = Field(default=0, description="体重(kg)")
