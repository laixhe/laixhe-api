<?php

namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

/**
 * 用户接口集成测试 (用户信息/列表/更新 + 分页边界钳制)
 *
 * 使用 sqlite 内存库 + RefreshDatabase, 与 AuthTest 相同设施, 不依赖外部 MySQL。
 * 分页钳制语义与 Go/Rust/TS 端一致: page<=0→1, page_size<=0→12, page_size>100→100。
 */
class UserTest extends TestCase
{
    use RefreshDatabase;

    private function uniqueEmail(): string
    {
        return 'user_' . time() . '_' . random_int(1000, 9999) . '@example.com';
    }

    /**
     * 注册一个测试用户, 返回 token (保证列表接口有数据)
     */
    private function registerToken(string $email): string
    {
        $resp = $this->postJson('/api/v1/auth/register', [
            'nickname' => '测试用户',
            'email' => $email,
            'password' => 'pass123',
        ]);
        $resp->assertStatus(200);
        return $resp->json('token');
    }

    public function test_list_default_pagination(): void
    {
        $this->registerToken($this->uniqueEmail());

        // 正常分页 (page=1, page_size=12 为默认值)
        $this->getJson('/api/v1/user/list?page=1&page_size=12')
            ->assertStatus(200)
            ->assertJsonStructure(['total', 'page', 'page_size', 'list'])
            ->assertJson(['page' => 1, 'page_size' => 12]);
    }

    public function test_list_page_zero_normalized_to_one(): void
    {
        $this->getJson('/api/v1/user/list?page=0&page_size=10')
            ->assertStatus(200)
            ->assertJson(['page' => 1, 'page_size' => 10]);
    }

    public function test_list_page_size_zero_uses_default(): void
    {
        // page_size<=0 回落默认 12 (与 Go/Rust/TS 端一致)
        $this->getJson('/api/v1/user/list?page=2&page_size=0')
            ->assertStatus(200)
            ->assertJson(['page' => 2, 'page_size' => 12]);
    }

    public function test_list_page_size_capped_at_100(): void
    {
        // 超大 page_size 钳制为 100, 防止恶意大分页触发全量查询
        $this->getJson('/api/v1/user/list?page=1&page_size=999')
            ->assertStatus(200)
            ->assertJson(['page' => 1, 'page_size' => 100]);
    }

    public function test_list_negative_params_normalized(): void
    {
        // 负数参数: page→1, page_size→12
        $this->getJson('/api/v1/user/list?page=-3&page_size=-5')
            ->assertStatus(200)
            ->assertJson(['page' => 1, 'page_size' => 12]);
    }

    public function test_list_non_numeric_params_returns_400(): void
    {
        // 非数字分页参数按请求格式错误返回 400 (与 Go/Rust 端绑定层行为一致)
        $this->getJson('/api/v1/user/list?page=abc&page_size=12')->assertStatus(400);
        $this->getJson('/api/v1/user/list?page=1&page_size=xyz')->assertStatus(400);
    }

    public function test_info_non_numeric_uid_returns_400(): void
    {
        // 非数字 uid 按请求格式错误返回 400 (与 Go/Rust 端绑定层行为一致)
        $this->getJson('/api/v1/user/info?uid=abc')->assertStatus(400);
    }

    public function test_info_oversized_uid_returns_400(): void
    {
        // 超出整数范围的 uid: filter_var FILTER_VALIDATE_INT 失败 → 400 (与 Go/Rust/TS 端溢出行为一致)
        $this->getJson('/api/v1/user/info?uid=99999999999999999999')->assertStatus(400);
    }

    public function test_update_nickname_emoji_count_by_codepoint(): void
    {
        // 昵称长度按 Unicode 字符计数 (Laravel min/max 用 mb_strlen, 与 Go/Rust/TS 端码点计数一致):
        // 20 个 emoji 应通过, 21 个应 422
        $email = $this->uniqueEmail();
        $token = $this->registerToken($email);

        $this->postJson('/api/v1/user/update', [
            'nickname' => str_repeat('😀', 20),
        ], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(200);

        $this->postJson('/api/v1/user/update', [
            'nickname' => str_repeat('😀', 21),
        ], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(422)
            ->assertJson(['message' => '昵称长度不能超过20位']);
    }

    public function test_info_invalid_uid_returns_422(): void
    {
        $this->getJson('/api/v1/user/info?uid=0')
            ->assertStatus(422)
            ->assertJson(['message' => '无效的用户ID']);
    }

    public function test_info_not_found_returns_422(): void
    {
        $this->getJson('/api/v1/user/info?uid=999999')
            ->assertStatus(422)
            ->assertJson(['message' => '用户不存在']);
    }

    public function test_info_success(): void
    {
        $email = $this->uniqueEmail();
        $resp = $this->postJson('/api/v1/auth/register', [
            'nickname' => '测试用户',
            'email' => $email,
            'password' => 'pass123',
        ]);
        $uid = $resp->json('user.uid');

        $this->getJson('/api/v1/user/info?uid=' . $uid)
            ->assertStatus(200)
            ->assertJson(['uid' => $uid, 'email' => $email])
            // 响应不得包含密码
            ->assertJsonMissing(['password' => '*']);
    }

    public function test_update_without_token_returns_401(): void
    {
        $this->postJson('/api/v1/user/update', [
            'nickname' => '新昵称',
            'avatar_url' => '',
        ])->assertStatus(401);
    }

    public function test_update_numeric_field_returns_400(): void
    {
        // 字段类型非字符串 (如数字): 请求格式错误 400 (与 Go/Rust/TS 端绑定层行为一致)
        $email = $this->uniqueEmail();
        $token = $this->registerToken($email);

        $this->postJson('/api/v1/user/update', [
            'nickname' => 123,
        ], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(400);
    }

    public function test_update_boolean_and_array_fields_return_400(): void
    {
        // 布尔字段 → 400 (与 Go/Rust/TS 端绑定层行为一致)
        $email = $this->uniqueEmail();
        $token = $this->registerToken($email);

        $this->postJson('/api/v1/user/update', [
            'nickname' => true,
        ], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(400);

        // 数组字段 → 400
        $this->postJson('/api/v1/user/update', [
            'nickname' => ['a', 'b'],
        ], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(400);
    }

    public function test_update_top_level_array_returns_400(): void
    {
        // 顶层 body 非对象 (数组) → 400 (与 Go/Rust 端绑定层行为一致)
        $email = $this->uniqueEmail();
        $token = $this->registerToken($email);

        $this->postJson('/api/v1/user/update', [1, 2], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(400);
    }

    public function test_update_with_token_success(): void
    {
        $email = $this->uniqueEmail();
        $token = $this->registerToken($email);

        $this->postJson('/api/v1/user/update', [
            'nickname' => '新昵称',
            'avatar_url' => 'https://example.com/a.png',
        ], ['Authorization' => 'Bearer ' . $token])
            ->assertStatus(200)
            ->assertJson(['nickname' => '新昵称', 'avatar_url' => 'https://example.com/a.png']);
    }
}
