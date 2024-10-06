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
    public string $uname = '';
    public string $login_at = '';

    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'uname' => ['required', 'string', 'min:2', 'max:30'],
            'login_at' => ['required', 'date_format:Y-m-d H:i:s'],
            //'login_at' => ['required', 'date_format:U'], // 验证时间戳
        ],
            [
                'uname' => '用户名长度在2~30之间！',
                'login_at' => '登录时间格式不对！',
            ]);
        if ($validator->fails()) {
            return new Result(ResultCode::Param, $validator->errors()->first(), []);
        }

        $this->param($params);

        return null;
    }

    public function param(array $params): void
    {
        $this->uname = $params['uname'];
        $this->login_at = $params['login_at'];
    }
}