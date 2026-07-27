<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Support\Str;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Symfony\Component\HttpFoundation\Response;
use Godruoyi\Snowflake\Snowflake;

class AssignRequestId
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {

//        $requestId = (string) Str::uuid();
        $snowflake = new Snowflake();
        $requestId = $snowflake->id();
        Log::withContext([
            'request_id' => $requestId
        ]);

        $response =  $next($request);
        $response->headers->set('X-Request-Id', $requestId);

        return $response;
    }
}
