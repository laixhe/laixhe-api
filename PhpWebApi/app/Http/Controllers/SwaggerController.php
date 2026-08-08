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
