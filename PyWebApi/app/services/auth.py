"""鉴权业务逻辑"""
import uuid

from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, select

from app.core.errors import param_error, unauthorized
from app.models.user import User, UserSex, UserState, UserType
from app.models.user_extend import UserExtend
from app.models.user_third_party import UserThirdParty
from app.schemas.auth import AuthLoginRequest, AuthRegisterRequest, AuthResponse
from app.schemas.user import User as UserSchema
from app.security.jwt import create_access_token
from app.security.password import hash_password, verify_password


class AuthService:
    def __init__(self, session: Session) -> None:
        self.session = session

    def register(self, req: AuthRegisterRequest) -> AuthResponse:
        # 先查邮箱是否已注册 (email 为唯一索引, 先查后插 + 唯一约束双重防重)
        if self.session.exec(select(User.id).where(User.email == req.email)).first() is not None:
            raise param_error("邮箱已存在")

        user = User(
            type_id=UserType.ORDINARY,
            account=uuid.uuid4().hex,
            mobile="",
            email=req.email,
            password=hash_password(req.password),
            nickname=req.nickname,
            avatar_url="",
            sex=UserSex.UNKNOWN,
            states=UserState.NORMAL,
        )
        self.session.add(user)
        try:
            # 与 Go 端 CreateUser 事务语义一致: 同时创建扩展信息与第三方关联记录
            self.session.flush()  # 先落库获取自增 uid
            self.session.add(UserExtend(uid=user.id))
            self.session.add(UserThirdParty(uid=user.id))
            self.session.commit()
        except IntegrityError:
            # 并发注册同邮箱等极端情况触发唯一约束冲突
            self.session.rollback()
            raise param_error("注册失败，请稍后再试")
        self.session.refresh(user)
        return AuthResponse(token=create_access_token(user.id), user=UserSchema.from_model(user))

    def login(self, req: AuthLoginRequest) -> AuthResponse:
        user = self.session.exec(select(User).where(User.email == req.email)).first()
        # 账号不存在/封禁/密码错误统一提示, 避免暴露账号状态
        if user is None or user.states != UserState.NORMAL or not verify_password(req.password, user.password):
            raise param_error("邮箱或密码错误")
        return AuthResponse(token=create_access_token(user.id), user=UserSchema.from_model(user))

    def refresh(self, uid: int) -> AuthResponse:
        user = self.session.get(User, uid)
        if user is None or user.states != UserState.NORMAL:
            raise unauthorized()
        return AuthResponse(token=create_access_token(user.id), user=UserSchema.from_model(user))
