<?php

namespace App\Utils;

use DateTimeZone;
use Lcobucci\JWT\Configuration;
use Lcobucci\JWT\Token;
use Lcobucci\JWT\Token\DataSet;
use Lcobucci\JWT\Signer\Hmac\Sha256;
use Lcobucci\JWT\Signer\Key\InMemory;

use App\Result\ResultCode;

class JwtUtil
{
    private static ?JwtUtil $instance = null; // 单例模式
    private Configuration $config;
    private string $secretKey;
    private int $expireTime; // 过期时长(单位秒)

    private function __construct()
    {
        $this->secretKey = env('JWT_SECRET', '6Kbj0VFeXYMp60lEyiFoVq4UzqX8Z0GSSfnvTh2VuAQn0oHgQNYexU6yYVTk4xf9');
        $this->expireTime = env('JWT_EXPIRE_TIME', 604800);

        $config = Configuration::forSymmetricSigner(
            new Sha256(),
            InMemory::plainText($this->secretKey)
        );
        $this->config = $config;
    }

    private function __clone()
    {
    }

    // 单例模式
    public static function getInstance(): JwtUtil
    {
        if (is_null(self::$instance)) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    // 创建 JWT Token
    public function createToken(array $claims): string
    {
        $config = $this->config;
        $builder = $config->builder();
        if (is_array($claims) && count(array_filter(array_keys($claims), 'is_string')) > 0) {
            $now = new \DateTimeImmutable();
            $builder = $builder->expiresAt($now->modify('+' . $this->expireTime . 'second'))
                ->issuedAt($now)->canOnlyBeUsedAfter($now);

            foreach ($claims as $k => $item) {
                $builder = $builder->withClaim($k, $item);
            }
            // 生成新令牌
            return $builder->getToken($config->signer(), $config->signingKey())->toString();
        }

        throw new \RuntimeException('claims 参数必须为关联数组', ResultCode::Unknown->value);
    }

    // 解析 token
    public function parseToken(string $jwt): DataSet
    {
        try {
            $config = $this->config;
            return $config->parser()->parse($jwt)->claims();
        } catch (\Exception $e) {
            throw new \RuntimeException('', ResultCode::AuthInvalid->value);
        }
    }

    // 验证令牌
    public function validatorToken($jwt): DataSet
    {
        $now = new \DateTimeImmutable();
        $config = $this->config;
        $token = null;
        try {
            $token = $config->parser()->parse($jwt);
        } catch (\Exception $e) {
            throw new \RuntimeException('', ResultCode::AuthInvalid->value);
        }

        $claims = $token->claims();

        $exp = $claims->get('exp')->setTimezone(new DateTimeZone(date_default_timezone_get()));
        if ($exp < $now) {
            throw new \RuntimeException('', ResultCode::AuthExpire->value);
        }
        $uid = (int)$claims->get('uid');
        if($uid <= 0) {
            throw new \RuntimeException('', ResultCode::AuthInvalid->value);
        }
        return $claims;
    }

}
