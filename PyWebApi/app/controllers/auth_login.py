from fastapi import APIRouter

from app.request.auth_login import AuthLoginRequest
from app.services.auth_login import service_auth_login
from app.result import result

router = APIRouter(
    prefix="/api/auth",
    tags=["鉴权相关"]
)


@router.post("/login", summary='登录')
async def auth_login(req: AuthLoginRequest) -> result.Result:
    """
    登录
    """

    await service_auth_login(req)

    return result.success(req)
