<?php
/**
 * 根据代码注解重新生成 swagger 文档 (zircote/swagger-php v6)
 *
 * 用法: composer swagger  (或直接运行 php scripts/generate-swagger.php)
 * 输出: public/swagger/swagger.json / public/swagger/swagger.yaml
 *
 * 与 Go 端 `swag init` 对应: 修改接口/参数/响应注解后重新运行即可。
 */

require __DIR__ . '/../vendor/autoload.php';

use OpenApi\Builder;

$result = (new Builder())
    ->addSource(__DIR__ . '/../app')
    ->build();

$jsonPath = __DIR__ . '/../public/swagger/swagger.json';
$yamlPath = __DIR__ . '/../public/swagger/swagger.yaml';

file_put_contents($jsonPath, $result->toJson() . PHP_EOL);
file_put_contents($yamlPath, $result->toYaml() . PHP_EOL);

echo "swagger docs regenerated -> public/swagger/swagger.json / swagger.yaml\n";
