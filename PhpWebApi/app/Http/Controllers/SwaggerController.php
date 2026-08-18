<?php

namespace App\Http\Controllers;

use Illuminate\Http\Response;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;

/**
 * Swagger 文档端点 (与 Go 端 /api/v1/swagger.json、/api/v1/swagger.yaml 对齐)
 *
 * 文档文件存放于 public/swagger/ 目录, 与 Go 端 docs/swagger.yaml 内容保持一致。
 */
class SwaggerController extends Controller
{
    /**
     * 返回 swagger.json
     *
     * @return Response
     */
    public function json(): Response
    {
        $content = $this->read('swagger.json');
        return response($content, 200)
            ->header('Content-Type', 'application/json')
            ->header('Cache-Control', 'public, max-age=300');
    }

    /**
     * 返回 swagger.yaml
     *
     * @return Response
     */
    public function yaml(): Response
    {
        $content = $this->read('swagger.yaml');
        return response($content, 200)
            ->header('Content-Type', 'application/x-yaml')
            ->header('Cache-Control', 'public, max-age=300');
    }

    /**
     * 返回 Swagger UI 页面 (CDN 加载 swagger-ui 资源, 与 Go/Rust/TS 端一致)
     *
     * @return Response
     */
    public function ui(): Response
    {
        $html = <<<'HTML'
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>API 接口文档</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
  <style>html{box-sizing:border-box;overflow-y:scroll}body{margin:0;background:#fafafa}</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/api/v1/swagger.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset]
      });
    };
  </script>
</body>
</html>
HTML;
        return response($html, 200)
            ->header('Content-Type', 'text/html; charset=utf-8')
            ->header('Cache-Control', 'public, max-age=300');
    }

    /**
     * 读取文档文件, 文件缺失时返回 404
     */
    private function read(string $filename): string
    {
        $path = public_path('swagger/' . $filename);
        $content = @file_get_contents($path);
        if ($content === false) {
            throw new NotFoundHttpException();
        }
        return $content;
    }
}
