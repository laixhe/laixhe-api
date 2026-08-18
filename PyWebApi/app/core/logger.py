"""日志配置 (console 模式, 与 Go 端 xlog console 对齐)"""
import logging
import sys

_FORMAT = "%(asctime)s [%(levelname)s] %(message)s"


def setup_logger(name: str = "pywebapi", level: int = logging.INFO) -> logging.Logger:
    logger = logging.getLogger(name)
    if not logger.handlers:
        handler = logging.StreamHandler(sys.stdout)
        handler.setFormatter(logging.Formatter(_FORMAT))
        logger.addHandler(handler)
        logger.setLevel(level)
        logger.propagate = False
    return logger


logger = setup_logger()
