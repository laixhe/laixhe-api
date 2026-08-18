"""Swagger 文档端点 (与 Go 端 /api/v1/swagger.json|yaml|swagger 对齐)

文档内容由 FastAPI 根据代码自动生成 (OpenAPI 3.1), 保持接口定义单一来源。
"""
import json

import yaml
from fastapi import APIRouter, Request
from fastapi.responses import HTMLResponse, Response

router = APIRouter(prefix="/api/v1", tags=["Swagger"])

# 文档响应缓存 (与 Go 端 Cache-Control 头一致)
CACHE_HEADERS = {"Cache-Control": "public, max-age=300"}

SWAGGER_UI_HTML = """<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>API接口 - Swagger UI</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
<style>body { margin: 0; }</style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
window.onload = function () {
  window.ui = SwaggerUIBundle({ url: "/api/v1/swagger.json", dom_id: "#swagger-ui" });
};
</script>
</body>
</html>
"""


@router.get("/swagger.json", summary="Swagger Json 文档")
def swagger_json(request: Request) -> Response:
    return Response(
        content=json.dumps(request.app.openapi(), ensure_ascii=False, indent=2),
        media_type="application/json",
        headers=CACHE_HEADERS,
    )


@router.get("/swagger.yaml", summary="Swagger Yaml 文档")
def swagger_yaml(request: Request) -> Response:
    content = yaml.safe_dump(request.app.openapi(), allow_unicode=True, sort_keys=False)
    return Response(
        content=content,
        media_type="application/x-yaml",
        headers=CACHE_HEADERS,
    )


@router.get("/swagger", summary="Swagger UI 页面")
def swagger_ui() -> HTMLResponse:
    return HTMLResponse(content=SWAGGER_UI_HTML)
