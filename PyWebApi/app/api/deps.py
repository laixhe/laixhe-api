"""路由依赖: JWT 鉴权"""
import jwt
from fastapi import Header

from app.core.errors import unauthorized
from app.security.jwt import decode_access_token


def get_current_uid(authorization: str = Header(description="Bearer XXX令牌")) -> int:
    """从 Authorization: Bearer <token> 中解析出用户 uid, 无效令牌返回 401"""
    if not authorization.startswith("Bearer "):
        raise unauthorized()
    try:
        payload = decode_access_token(authorization.removeprefix("Bearer "))
    except jwt.PyJWTError:
        raise unauthorized()
    uid = payload.get("uid")
    if not uid:
        raise unauthorized()
    return uid
