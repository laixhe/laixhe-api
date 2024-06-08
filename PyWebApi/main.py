import uvicorn
from fastapi import FastAPI

from app.controllers import RegisterRouterList
from app.errors import register_custom_error_handle

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
    uvicorn.run(app, host="0.0.0.0", port=9090)
