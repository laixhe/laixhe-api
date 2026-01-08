<?php

namespace App\Http\Request;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;

/**
 * 修改用户信息请求参数
 */
class UserUpdateRequest implements IRequest
{
    public string $nickname = '';
    public string $avatar_url = '';

    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'nickname' => ['required', 'string', 'min:2', 'max:20'],
            'avatar_url' => ['required'],
        ],
            [
                'nickname' => '用户名长度在2~20之间！',
                'avatar_url' => '头像地址不能为空！',
            ]);
        if ($validator->fails()) {
            return new Result(ResultCode::Param, $validator->errors()->first(), []);
        }
        $this->param($params);
        return null;
    }

    public function param(array $params): void
    {
        $this->nickname = $params['nickname'];
        $this->avatar_url = $params['avatar_url'];
    }
}