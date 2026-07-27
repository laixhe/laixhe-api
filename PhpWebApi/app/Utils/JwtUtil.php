<?php

namespace App\Utils;

use DateTimeZone;
use DateTimeImmutable;
use Throwable;
use RuntimeException;

use Lcobucci\JWT\Configuration;
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
        $this->secretKey = env('JWT_SECRET', '');
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

    /**
     * 单例模式
     * @return JwtUtil
     */
    public static function getInstance(): JwtUtil
    {
        if (is_null(self::$instance)) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    /**
     * 创建 JWT Token
     * @param array $claims
     * @return string
     * @throws RuntimeException
     */
    public function createToken(int $uid, array $claims=[]): string
    {
        $config = $this->config;
        $builder = $config->builder();
        try {
            $now = new DateTimeImmutable();
            $expiresAt = $now->modify('+' . $this->expireTime . ' second');
            if (empty($expiresAt)) {
                throw new RuntimeException('创建 JWT Token 过期时间生成失败');
            }

            $builder = $builder->expiresAt($expiresAt)->issuedAt($now)->canOnlyBeUsedAfter($now);
            $builder = $builder->withClaim('uid', $uid);

            // 判断是否为数组并且是一维关联数组
            if (!empty($claims)) {
                $claims_keys = array_keys($claims);
                if (count($claims_keys) !== count(array_filter($claims_keys, 'is_string'))){
                    throw new RuntimeException('创建 JWT Token 参数 claims 必须为关联数组');
                }
            }
            foreach ($claims as $k => $item) {
                $builder = $builder->withClaim($k, $item);
            }
            // 生成新令牌
            return $builder->getToken($config->signer(), $config->signingKey())->toString();
        } catch (Throwable $e) {
            // echo $e->getMessage();
            throw new RuntimeException('', ResultCode::Service->value);
        }
    }

    /**
     * 解析 token
     * @param string $jwt
     * @return DataSet
     *
     * @throws RuntimeException
     */
    public function parseToken(string $jwt): DataSet
    {
        try {
            $config = $this->config;
            return $config->parser()->parse($jwt)->claims();
        } catch (Throwable $e) {
            // echo $e->getMessage();
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }
    }

    /**
     * 验证令牌
     * @param $jwt
     * @return DataSet
     *
     * @throws RuntimeException
     */
    public function validatorToken($jwt): DataSet
    {
        $now = new DateTimeImmutable();
        $config = $this->config;
        $token = null;
        try {
            $token = $config->parser()->parse($jwt);
        } catch (Throwable $e) {
            // echo $e->getMessage();
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }

        $claims = $token->claims();

        $exp = $claims->get('exp')->setTimezone(new DateTimeZone(date_default_timezone_get()));
        if ($exp < $now) {
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }
        $uid = (int)$claims->get('uid');
        if ($uid <= 0) {
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }
        return $claims;
    }

}
