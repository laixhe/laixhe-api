from fastapi import APIRouter

from app.request.auth_login import AuthLoginRequest
from app.services.auth_login import service_auth_login
from app.types import response

router = APIRouter(
    prefix="/api/auth",
    tags=["鉴权相关"]
)


@router.post("/login", summary='登录')
async def auth_login(req: AuthLoginRequest) -> response.HttpResponse:
    """
    登录
    """

    await service_auth_login(req)

    return response.success(req)
