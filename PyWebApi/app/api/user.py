"""用户接口"""
from fastapi import APIRouter, Depends, Query
from sqlmodel import Session

from app.api.deps import get_current_uid
from app.core.errors import param_error
from app.db.database import get_session
from app.schemas.common import Error
from app.schemas.user import User as UserSchema
from app.schemas.user import UserListResponse, UserUpdateRequest
from app.services.user import UserService
from app.utils.validators import validate_avatar_url, validate_nickname

router = APIRouter(prefix="/api/v1/user", tags=["User"])

error_responses = {400: {"model": Error}, 422: {"model": Error}, 500: {"model": Error}}


def normalize_pagination(page: int, page_size: int) -> tuple[int, int]:
    """归一化分页参数 (与各端钳制逻辑保持一致): page<=0→1, page_size<=0→12, page_size>100→100"""
    if page <= 0:
        page = 1
    if page_size <= 0:
        page_size = 12
    if page_size > 100:
        page_size = 100
    return page, page_size


@router.get(
    "/info",
    response_model=UserSchema,
    summary="获取用户信息",
    responses=error_responses,
)
def info(uid: int = Query(description="用户id"), session: Session = Depends(get_session)) -> UserSchema:
    if uid <= 0:
        raise param_error("无效的用户ID")
    return UserService(session).info(uid)


@router.get(
    "/list",
    response_model=UserListResponse,
    summary="获取用户列表",
    responses=error_responses,
)
def list(
    page: int = Query(default=1, description="分页-当前页(默认 1)"),
    page_size: int = Query(default=12, description="分页-每页数量(默认 12)"),
    session: Session = Depends(get_session),
) -> UserListResponse:
    page, page_size = normalize_pagination(page, page_size)
    return UserService(session).list(page, page_size)


@router.post(
    "/update",
    response_model=UserSchema,
    summary="更新用户信息",
    responses={400: {"model": Error}, 401: {"model": Error}, 422: {"model": Error}, 500: {"model": Error}},
)
def update(
    req: UserUpdateRequest,
    uid: int = Depends(get_current_uid),
    session: Session = Depends(get_session),
) -> UserSchema:
    validate_nickname(req.nickname)
    validate_avatar_url(req.avatar_url)
    return UserService(session).update(uid, req)
