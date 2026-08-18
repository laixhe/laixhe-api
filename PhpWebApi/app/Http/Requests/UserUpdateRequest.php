<?php

namespace App\Http\Requests;

use App\Result\Result;
use Illuminate\Support\Facades\Validator;
use OpenApi\Attributes as OA;

/**
 * 修改用户信息请求参数 (与 Go 端 controllers/user.go 校验对齐)
 */
#[OA\Schema(
    schema: 'UserUpdateRequest',
    description: '修改用户信息请求参数 (Uid 由 JWT 提供)',
    required: ['nickname'],
)]
class UserUpdateRequest implements IRequest
{
    #[OA\Property(description: '昵称', example: '张三')]
    public string $nickname = '';
    #[OA\Property(description: '头像地址', example: '')]
    public string $avatar_url = '';

    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            // 昵称按"字符"计数 2~20 位
            'nickname' => ['required', 'string', 'min:2', 'max:20'],
            // 头像地址可选, 非空时由控制器校验必须以 http/https 开头
            'avatar_url' => ['string', 'max:255'],
        ],
            [
                'nickname' => '昵称长度不能小于2位',
                'nickname.min' => '昵称长度不能小于2位',
                'nickname.max' => '昵称长度不能超过20位',
                'avatar_url.max' => '头像地址长度不能超过255位',
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
        $this->nickname = $params['nickname'] ?? '';
        $this->avatar_url = $params['avatar_url'] ?? '';
    }
}
