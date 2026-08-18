"""密码哈希: bcrypt 单向不可逆 (与 Go 端 bcrypt cost=10 对齐)"""
import bcrypt


def hash_password(password: str) -> str:
    # cost=10 与 Go 端一致 (bcrypt 5.x 默认 12, 显式指定保持两端哈希成本相同)
    return bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt(rounds=10)).decode("utf-8")


def verify_password(password: str, hashed: str) -> bool:
    try:
        return bcrypt.checkpw(password.encode("utf-8"), hashed.encode("utf-8"))
    except ValueError:
        # 存量哈希格式异常等场景视为校验失败
        return False
