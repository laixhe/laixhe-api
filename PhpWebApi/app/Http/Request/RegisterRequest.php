<?php

namespace App\Http\Request;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\Rule;

/**
 * 注册请求参数
 */
class RegisterRequest implements IRequest
{
    public string $nickname = '';
    public string $email = '';
    public string $password = '';

    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'nickname' => ['required', 'string', 'min:2', 'max:20'],
            'email' => ['required', 'email'],
            'password' => ['required', 'string', 'min:6', 'max:20'],
//            'type_id' => ['required', Rule::in([1,2,3])]
        ],
            [
                'nickname' => '用户名长度在2~20之间！',
                'email' => '邮箱格式不正确！',
                'password' => '密码长度在6~20之间！',
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
        $this->email = $params['email'];
        $this->password = $params['password'];
    }
}
