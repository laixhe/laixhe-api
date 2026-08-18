<?php

namespace App\Http\Requests;

use App\Result\Result;
use Illuminate\Support\Facades\Validator;
use OpenApi\Attributes as OA;

/**
 * 登录请求参数 (与 Go 端 controllers/auth.go 校验对齐)
 */
#[OA\Schema(
    schema: 'AuthLoginRequest',
    description: '登录请求参数',
    required: ['email', 'password'],
)]
class LoginRequest implements IRequest
{
    #[OA\Property(description: '邮箱', example: 'user@example.com')]
    public string $email = '';
    #[OA\Property(description: '密码', example: 'abc123')]
    public string $password = '';


    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'email' => ['required', 'email'],
            // 密码只能包含字母 数字 _ @ $, 长度 6~64 位 (上限防 bcrypt 72 字节静默截断)
            'password' => ['required', 'string', 'min:6', 'max:64', 'regex:/^[a-zA-Z0-9_@$]+$/'],
        ],
            [
                'email' => '邮箱格式错误',
                'password.min' => '密码长度不能小于6位',
                'password.max' => '密码长度不能超过64位',
                'password.regex' => '密码格式错误，只能包含字母 数字 _ @ $',
            ]);
        if ($validator->fails()) {
            // 字段类型非字符串 (如数字): 请求格式错误, 返回 400 (与 Go/Rust/TS 端绑定层行为一致);
            // 其余业务校验失败仍为 422; 收敛在 Helpers::validation_error, 避免三处复制粘贴
            return validation_error($validator);
        }

        $this->param($params);

        return null;
    }

    public function param(array $params): void
    {
        $this->email = $params['email'] ?? '';
        $this->password = $params['password'] ?? '';
    }
}
