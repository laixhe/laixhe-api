package com.laixhe.webapi.controller;

import io.swagger.v3.oas.annotations.Operation;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.http.CacheControl;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.concurrent.TimeUnit;

/**
 * Swagger 文档端点 (与 Go 版 /api/v1/swagger.* 对齐)
 * 文档内容由 springdoc 从代码注解动态生成 (原始端点 /v3/api-docs 与 /v3/api-docs.yaml),
 * 本控制器仅做路径转发, 并承载 Swagger UI 页面。
 */
@RestController
@RequestMapping("/api/v1")
public class SwaggerDocController {

    /** 隐藏: 文档端点本身不进入 OpenAPI 文档 */
    @Operation(hidden = true)
    @GetMapping(value = "/swagger.yaml")
    public void swaggerYaml(HttpServletRequest request, HttpServletResponse response) throws Exception {
        request.getRequestDispatcher("/v3/api-docs.yaml").forward(request, response);
    }

    @Operation(hidden = true)
    @GetMapping(value = "/swagger.json")
    public void swaggerJson(HttpServletRequest request, HttpServletResponse response) throws Exception {
        request.getRequestDispatcher("/v3/api-docs").forward(request, response);
    }

    @Operation(hidden = true)
    @GetMapping(value = "/swagger", produces = "text/html;charset=UTF-8")
    public ResponseEntity<String> swaggerUi() {
        return ResponseEntity.ok()
                .cacheControl(CacheControl.maxAge(300, TimeUnit.SECONDS).cachePublic())
                .body(SWAGGER_UI_HTML);
    }

    private static final String SWAGGER_UI_HTML = """
            <!DOCTYPE html>
            <html lang="zh-CN">
            <head>
              <meta charset="utf-8">
              <title>API接口</title>
              <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
            </head>
            <body>
              <div id="swagger-ui"></div>
              <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
              <script>
                window.onload = function () {
                  window.ui = SwaggerUIBundle({
                    url: "/api/v1/swagger.yaml",
                    dom_id: "#swagger-ui",
                    deepLinking: true,
                    presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
                    layout: "BaseLayout"
                  });
                };
              </script>
            </body>
            </html>
            """;
}
