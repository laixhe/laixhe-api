<?php

namespace App\Http\Request;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\Rule;

class RegisterRequest implements IRequest
{
    public string $email = '';
    public string $password = '';
    public string $uname = '';
    public int $age = 0;


    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'email' => ['required', 'email'],
            'password' => ['required', 'string', 'min:6', 'max:20'],
            'uname' => ['required', 'string', 'min:2', 'max:30'],
            'age' => ['required', 'int', 'gte:0', 'lte:200'],
//            'type_id' => ['required', Rule::in([1,2,3])]
        ],
            [
                'email' => '邮箱格式不正确！',
                'password' => '密码长度在6~20之间！',
                'uname' => '用户长度在2~30之间！',
                'age' => '年龄在0~200之间！',
            ]);
        if ($validator->fails()) {
            return new Result(ResultCode::Param, $validator->errors()->first(), []);
        }

        $this->param($params);

        return null;
    }

    public function param(array $params): void
    {
        $this->email = $params['email'];
        $this->password = $params['password'];
        $this->uname = $params['uname'];
        $this->age = (int)$params['age'];
    }
}
