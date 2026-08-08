<?php

namespace App\Http\Requests;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;

/**
 * 修改用户信息请求参数 (与 Go 端 controllers/user.go 校验对齐)
 */
class UserUpdateRequest implements IRequest
{
    public string $nickname = '';
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
            return new Result(ResultCode::Param, $validator->errors()->first());
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
