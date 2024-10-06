import uvicorn
from fastapi import FastAPI

from app.controllers import RegisterRouterList
from app.exception import register_custom_error_handle
from config import appSettings

# 实例化
# docs_url=None: 代表关闭 SwaggerUi
# redoc_url=None: 代表关闭 redoc 文档
# app = FastAPI(docs_url=None, redoc_url=None)
app = FastAPI()

# 注册自定义错误处理器
register_custom_error_handle(app)

# 加载路由
for item in RegisterRouterList:
    app.include_router(item.router)


if __name__ == "__main__":
    print(appSettings)
    uvicorn.run(app="main:app", host=appSettings.app_host, port=appSettings.app_port, reload=appSettings.app_debug)
