"""用户业务逻辑"""
from datetime import datetime

from sqlmodel import Session, func, select

from app.core.errors import param_error, unauthorized
from app.models.user import User, UserState
from app.schemas.user import User as UserSchema
from app.schemas.user import UserListResponse, UserUpdateRequest


class UserService:
    def __init__(self, session: Session) -> None:
        self.session = session

    def info(self, uid: int) -> UserSchema:
        user = self.session.get(User, uid)
        if user is None:
            raise param_error("用户不存在")
        return UserSchema.from_model(user)

    def update(self, uid: int, req: UserUpdateRequest) -> UserSchema:
        user = self.session.get(User, uid)
        if user is None:
            raise param_error("用户不存在")
        if user.states != UserState.NORMAL:
            raise unauthorized()
        user.nickname = req.nickname
        # 与 Go 端 UpdateUser 非零字段更新语义一致: avatar_url 为空时不覆盖原值
        if req.avatar_url:
            user.avatar_url = req.avatar_url
        user.updated_at = datetime.now()
        self.session.add(user)
        self.session.commit()
        self.session.refresh(user)
        return UserSchema.from_model(user)

    def list(self, page: int, page_size: int) -> UserListResponse:
        total = self.session.exec(select(func.count(User.id))).one()
        users = self.session.exec(
            select(User).order_by(User.id.desc()).offset((page - 1) * page_size).limit(page_size)
        ).all()
        return UserListResponse(
            total=total,
            page=page,
            page_size=page_size,
            list=[UserSchema.from_model(u) for u in users],
        )
