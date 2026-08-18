"""JWT 生成与校验 (Header.Payload.Signature 三段结构)"""
from datetime import datetime, timedelta, timezone

import jwt

from app.core.config import settings


def create_access_token(uid: int) -> str:
    now = datetime.now(timezone.utc)
    payload = {
        "uid": uid,
        "iat": now,
        "nbf": now,  # 生效时间 (与 Go 端 JwtClaims 对齐)
        "exp": now + timedelta(seconds=settings.jwt_expire_seconds),
    }
    return jwt.encode(payload, settings.jwt_secret_key, algorithm=settings.jwt_signing_algorithm)


def decode_access_token(token: str) -> dict:
    """解码并校验签名与过期时间; 无效时抛出 jwt.PyJWTError"""
    return jwt.decode(token, settings.jwt_secret_key, algorithms=[settings.jwt_signing_algorithm])
