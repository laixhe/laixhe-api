<?php

namespace Tests\Unit;

use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;
use RuntimeException;
use Tests\TestCase;

/**
 * 全局辅助函数单元测试 (纯逻辑, 无需数据库)
 *
 * 覆盖:
 * - validation_error: 字段类型非法 → 400, 业务校验失败 → 422 (与 Go/Rust/TS 端绑定层一致)
 * - response_exception_safe: 已知业务错误透传; 未知异常 (含 QueryException 的完整 SQL) 不泄露,
 *   只返回固定 500 文案 (与 bootstrap/app.php 统一 500 防护对齐)
 * - jwt.expire_time 默认值: 30 天, 与 README/.env.example 及其它三端一致
 */
class HelpersTest extends TestCase
{
    public function test_jwt_expire_time_defaults_to_30_days(): void
    {
        // 兜底值 2592000 (30 天) 与 README/.env.example 及其它三端配置一致
        $this->assertSame(2592000, (int) config('jwt.expire_time'));
    }

    public function test_validation_error_type_error_returns_400(): void
    {
        // 字段类型非字符串 (如数字) → 400 (与 Go/Rust/TS 端绑定层行为一致)
        $v = Validator::make(['nickname' => 123], ['nickname' => ['string']]);
        $v->fails();
        $result = validation_error($v);
        $this->assertSame(ResultCode::BadRequest->value, $result->code->value);
        $this->assertSame('Bad Request', $result->message);
    }

    public function test_validation_error_business_rule_returns_422(): void
    {
        // 纯业务校验失败 (长度不足) → 422 具体文案
        $v = Validator::make(
            ['nickname' => 'a'],
            ['nickname' => ['string', 'min:2']],
            ['nickname.min' => '昵称长度不能小于2位']
        );
        $v->fails();
        $result = validation_error($v);
        $this->assertSame(ResultCode::Param->value, $result->code->value);
        $this->assertSame('昵称长度不能小于2位', $result->message);
    }

    public function test_response_exception_safe_passes_known_business_error(): void
    {
        // 已知业务错误 (如"用户不存在" 422): 消息透传客户端
        $e = new RuntimeException('用户不存在', ResultCode::Param->value);
        $resp = response_exception_safe($e, 'user/update');
        $this->assertSame(422, $resp->getStatusCode());
        // JSON 内容默认转义非 ASCII, 解码后断言 message (避免 \uXXXX 转义干扰)
        $json = json_decode((string) $resp->getContent(), true);
        $this->assertSame('用户不存在', $json['message']);
    }

    public function test_response_exception_safe_masks_unknown_error_sql(): void
    {
        // 未知异常 (QueryException 的完整 SQL 已含在 getMessage 中): 只返回固定 500 文案,
        // 不得把 SQL 等内部细节泄露给客户端
        $e = new RuntimeException('SQLSTATE[42S22]: Column not found: 1054 Unknown column... SQL: select * from `user`');
        $resp = response_exception_safe($e, 'auth/register');
        $this->assertSame(500, $resp->getStatusCode());
        $content = (string) $resp->getContent();
        $this->assertStringNotContainsString('SQLSTATE', $content);
        $this->assertStringNotContainsString('select * from', $content, 'SQL 不应出现在响应中');
    }

    public function test_response_exception_safe_auth_invalid_falls_back_to_default_text(): void
    {
        // 已知错误码但空消息: 回落 ResultCode 默认文案 (与 response_exception 行为一致)
        $e = new RuntimeException('', ResultCode::AuthInvalid->value);
        $resp = response_exception_safe($e, 'auth/refresh');
        $this->assertSame(401, $resp->getStatusCode());
        $this->assertStringContainsString('Unauthorized', (string) $resp->getContent());
    }
}
