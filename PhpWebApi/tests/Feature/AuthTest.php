<?php

namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

/**
 * 鉴权接口集成测试 (注册/登录/限流)
 *
 * 使用 sqlite 内存库 + RefreshDatabase, 每个用例自动重建表结构,
 * 不依赖外部 MySQL。限流中间件在测试中通过覆盖 config 触发 429。
 */
class AuthTest extends TestCase
{
    use RefreshDatabase;

    private function uniqueEmail(): string
    {
        return 'test_' . time() . '_' . random_int(1000, 9999) . '@example.com';
    }

    private function registerPayload(string $email): array
    {
        return [
            'nickname' => '测试用户',
            'email' => $email,
            'password' => 'pass123',
        ];
    }

    public function test_register_success(): void
    {
        $resp = $this->postJson('/api/v1/auth/register', $this->registerPayload($this->uniqueEmail()));

        $resp->assertStatus(200)
            ->assertJsonStructure(['token', 'user' => ['uid', 'email', 'nickname']]);
        // 响应不得包含密码
        $resp->assertJsonMissing(['user' => ['password' => '*']]);
    }

    public function test_register_duplicate_email_returns_422(): void
    {
        $email = $this->uniqueEmail();

        $this->postJson('/api/v1/auth/register', $this->registerPayload($email))->assertStatus(200);
        $resp = $this->postJson('/api/v1/auth/register', $this->registerPayload($email));

        $resp->assertStatus(422)
            ->assertJson(['message' => '邮箱已存在']);
    }

    public function test_login_success_and_wrong_password(): void
    {
        $email = $this->uniqueEmail();
        $this->postJson('/api/v1/auth/register', $this->registerPayload($email))->assertStatus(200);

        // 正确密码
        $this->postJson('/api/v1/auth/login', [
            'email' => $email,
            'password' => 'pass123',
        ])->assertStatus(200)->assertJsonStructure(['token', 'user']);

        // 错误密码: 与账号不存在同文案, 防账号探测
        $this->postJson('/api/v1/auth/login', [
            'email' => $email,
            'password' => 'wrong123',
        ])->assertStatus(422)->assertJson(['message' => '邮箱或密码错误']);
    }

    public function test_login_invalid_params_returns_422(): void
    {
        // 邮箱格式错误
        $this->postJson('/api/v1/auth/login', [
            'email' => 'not-an-email',
            'password' => 'pass123',
        ])->assertStatus(422)->assertJson(['message' => '邮箱格式错误']);
    }

    public function test_register_numeric_field_returns_400(): void
    {
        // 字段类型非字符串 (如数字): 请求格式错误 400 (与 Go/Rust/TS 端绑定层行为一致)
        $this->postJson('/api/v1/auth/register', [
            'nickname' => 123,
            'email' => 'a@b.com',
            'password' => 'pass123',
        ])->assertStatus(400);
    }

    public function test_login_numeric_field_returns_400(): void
    {
        $this->postJson('/api/v1/auth/login', [
            'email' => 'a@b.com',
            'password' => 123,
        ])->assertStatus(400);
    }

    public function test_register_top_level_array_returns_400(): void
    {
        // 顶层 body 非对象 (数组) → 400 (与 Go/Rust 端绑定层行为一致)
        $this->postJson('/api/v1/auth/register', [1, 2])->assertStatus(400);
    }

    public function test_register_top_level_null_returns_400(): void
    {
        // 顶层 body 为 JSON null → 400 (与 Go/Rust 端绑定层行为一致)
        $this->call('POST', '/api/v1/auth/register', [], [], [], ['CONTENT_TYPE' => 'application/json'], 'null')
            ->assertStatus(400);
    }

    public function test_not_found_returns_unified_json(): void
    {
        // 未匹配路由 → 统一 JSON 格式 (与 Rust/TS 端对齐)
        $this->getJson('/api/v1/no-such-route')
            ->assertStatus(404)
            ->assertJson(['code' => 404, 'message' => 'Not Found']);
    }

    public function test_cors_headers_present(): void
    {
        // CORS 响应头 (与 Go/Rust/TS 端对齐)
        $this->get('/api/v1/health')
            ->assertHeader('Access-Control-Allow-Origin', '*')
            ->assertHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, PATCH, DELETE, OPTIONS');
    }

    public function test_rate_limit_returns_429(): void
    {
        // 覆盖限流阈值为 2, 使测试可快速触发 429 (生产默认 1000)
        // 语义与 Go/Rust/TS 端一致: 第 max+1 个请求被拒绝
        config(['rate_limit.max' => 2]);
        $email = $this->uniqueEmail();

        // 第 1、2 个请求放行
        $this->postJson('/api/v1/auth/register', $this->registerPayload($email))->assertStatus(200);
        $this->postJson('/api/v1/auth/login', ['email' => $email, 'password' => 'wrong'])->assertStatus(422);
        // 第 3 个请求触发限流 429
        $this->postJson('/api/v1/auth/login', ['email' => $email, 'password' => 'wrong'])
            ->assertStatus(429)
            ->assertJson(['code' => 429]);
    }

    public function test_refresh_without_token_returns_401(): void
    {
        // 刷新接口需要 JWT: 无 token 返回 401 统一 JSON (与 Go/Rust/TS 端一致)
        $this->postJson('/api/v1/auth/refresh')
            ->assertStatus(401)
            ->assertJson(['code' => 401]);
    }

    public function test_refresh_with_invalid_token_returns_401(): void
    {
        $this->postJson('/api/v1/auth/refresh', [], ['Authorization' => 'Bearer invalid.token.value'])
            ->assertStatus(401)
            ->assertJson(['code' => 401]);
    }

    public function test_refresh_with_token_success(): void
    {
        $email = $this->uniqueEmail();
        $token = $this->postJson('/api/v1/auth/register', $this->registerPayload($email))->json('token');

        $this->postJson('/api/v1/auth/refresh', [], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(200)
            ->assertJsonStructure(['token', 'user' => ['uid', 'email']]);
    }

    public function test_refresh_user_not_found_returns_401(): void
    {
        // 用有效 JWT 但用户已被删除: refresh 返回 401 (与 Go 端一致)
        $email = $this->uniqueEmail();
        $resp = $this->postJson('/api/v1/auth/register', $this->registerPayload($email));
        $token = $resp->json('token');
        $uid = $resp->json('user.uid');
        \App\Models\User::query()->where('id', $uid)->delete();

        $this->postJson('/api/v1/auth/refresh', [], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(401)
            ->assertJson(['code' => 401]);
    }
}
