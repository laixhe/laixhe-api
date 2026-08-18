"""参数格式校验 (与 Go 端 controllers 层手工校验对齐, 失败返回 422)"""
import re

from app.core.errors import param_error

# 邮箱格式 (轻量正则, 仅做基本格式校验)
EMAIL_RE = re.compile(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$")
# 密码规则: 长度 6~64 位, 仅含字母 数字 _ @ $ (上限 64 保证 bcrypt 72 字节内不会截断)
PASSWORD_RE = re.compile(r"^[a-zA-Z0-9_@$]{6,64}$")


def validate_nickname(nickname: str) -> None:
    """昵称 2~20 位 (按字符计数, 中文字符按 1 位)"""
    if len(nickname) < 2:
        raise param_error("昵称长度不能小于2位")
    if len(nickname) > 20:
        raise param_error("昵称长度不能超过20位")


def validate_email_and_password(email: str, password: str) -> None:
    if not EMAIL_RE.match(email):
        raise param_error("邮箱格式错误")
    if not PASSWORD_RE.match(password):
        raise param_error("密码格式错误，需 6~64 位，只能包含字母 数字 _ @ $")


def validate_avatar_url(avatar_url: str) -> None:
    if len(avatar_url) > 255:
        raise param_error("头像地址长度不能超过255位")
    if avatar_url and not (avatar_url.startswith("http://") or avatar_url.startswith("https://")):
        raise param_error("头像地址必须以http或https开头")
