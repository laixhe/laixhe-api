"""通用配置模型 (对齐 webapi.sql 的 config_common 表)"""
from sqlmodel import Field, SQLModel


class ConfigCommon(SQLModel, table=True):
    """通用配置表"""
    __tablename__ = "config_common"

    id: int | None = Field(default=None, primary_key=True)
    key: str = Field(default="", max_length=255, index=True, description="配置键")
    value: str = Field(default="", max_length=512, description="配置值")
    describe: str = Field(default="", max_length=255, description="描述")
