"""路由汇总"""
from fastapi import APIRouter

from app.api import auth, health, swagger, user

api_router = APIRouter()
api_router.include_router(auth.router)
api_router.include_router(user.router)
api_router.include_router(health.router)
api_router.include_router(swagger.router)
