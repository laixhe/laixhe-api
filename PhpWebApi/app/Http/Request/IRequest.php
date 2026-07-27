<?php

namespace App\Http\Request;

use App\Result\Result;

interface IRequest
{
    /**
     * 参数验证
     * @param array $params
     * @return Result|null
     */
    public function validator(array $params): ?Result;

    /**
     * 参数填充
     * @param array $params
     * @return void
     */
    public function param(array $params) : void;
}
