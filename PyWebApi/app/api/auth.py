"""鉴权接口"""
from fastapi import APIRouter, Depends
from sqlmodel import Session

from app.api.deps import get_current_uid
from app.db.database import get_session
from app.schemas.auth import AuthLoginRequest, AuthRegisterRequest, AuthResponse
from app.schemas.common import Error
from app.services.auth import AuthService
from app.utils.validators import validate_email_and_password, validate_nickname

router = APIRouter(prefix="/api/v1/auth", tags=["Auth"])

error_responses = {400: {"model": Error}, 422: {"model": Error}, 500: {"model": Error}}


@router.post(
    "/register",
    response_model=AuthResponse,
    summary="注册",
    responses=error_responses,
)
def register(req: AuthRegisterRequest, session: Session = Depends(get_session)) -> AuthResponse:
    validate_nickname(req.nickname)
    validate_email_and_password(req.email, req.password)
    return AuthService(session).register(req)


@router.post(
    "/login",
    response_model=AuthResponse,
    summary="登录",
    responses=error_responses,
)
def login(req: AuthLoginRequest, session: Session = Depends(get_session)) -> AuthResponse:
    validate_email_and_password(req.email, req.password)
    return AuthService(session).login(req)


@router.post(
    "/refresh",
    response_model=AuthResponse,
    summary="刷新Jwt",
    responses={400: {"model": Error}, 401: {"model": Error}, 500: {"model": Error}},
)
def refresh(uid: int = Depends(get_current_uid), session: Session = Depends(get_session)) -> AuthResponse:
    return AuthService(session).refresh(uid)
