<?php

namespace App\Http\Request;

use App\Result\Result;

interface IRequest
{
    public function validator(array $params): ?Result;

    public function param(array $params) : void;
}
