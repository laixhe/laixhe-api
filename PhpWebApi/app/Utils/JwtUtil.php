<?php

namespace App\Utils;

use DateTimeImmutable;
use Throwable;
use RuntimeException;

use Lcobucci\Clock\SystemClock;
use Lcobucci\JWT\Configuration;
use Lcobucci\JWT\Token\DataSet;
use Lcobucci\JWT\Signer\Hmac\Sha256;
use Lcobucci\JWT\Signer\Key\InMemory;
use Lcobucci\JWT\Validation\Constraint\SignedWith;
use Lcobucci\JWT\Validation\Constraint\StrictValidAt;

use App\Result\ResultCode;

class JwtUtil
{
    private static ?JwtUtil $instance = null; // 单例模式
    private Configuration $config;
    private string $secretKey;
    private int $expireTime; // 过期时长(单位秒)

    private function __construct()
    {
        // 必须用 config() 读取 (见 config/jwt.php 顶部说明): env() 在 config:cache 后会返回 null
        $this->secretKey = (string)config('jwt.secret', '');
        // 密钥为空时直接报错, 避免用空密钥签发可被伪造的 token (新手最容易踩的坑)
        if ($this->secretKey === '') {
            throw new RuntimeException('JWT_SECRET 未配置, 请检查 .env', ResultCode::Service->value);
        }
        // 默认 30 天, 与 config/jwt.php 兜底值一致 (修改任一处需同步)
        $this->expireTime = (int)config('jwt.expire_time', 2592000);

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
     * 创建 JWT Token (HS256 签名, 含 uid/exp/iat/nbf 声明)
     *
     * @param int   $uid    用户id
     * @param array $claims 附加自定义声明 (可选, 必须为一维关联数组)
     * @return string
     *
     * @throws RuntimeException 过期时间生成失败 / claims 非关联数组时抛出 Service 错误
     */
    public function createToken(int $uid, array $claims = []): string
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

            // 附加自定义声明: 必须是一维关联数组 (key 全为字符串)
            if (!empty($claims)) {
                $claims_keys = array_keys($claims);
                if (count($claims_keys) !== count(array_filter($claims_keys, 'is_string'))) {
                    throw new RuntimeException('创建 JWT Token 参数 claims 必须为关联数组');
                }
            }
            foreach ($claims as $k => $item) {
                $builder = $builder->withClaim($k, $item);
            }
            // 生成新令牌
            return $builder->getToken($config->signer(), $config->signingKey())->toString();
        } catch (Throwable $e) {
            throw new RuntimeException('', ResultCode::Service->value);
        }
    }

    /**
     * 解析 token (仅解码, 不校验签名/有效期!)
     *
     * 注意: 该方法只做结构解析, 不能用于鉴权。
     * 需要校验 token 合法性请使用 validatorToken()。
     *
     * @param string $jwt
     * @return DataSet
     *
     * @throws RuntimeException token 结构不合法时抛出 AuthInvalid
     */
    public function parseToken(string $jwt): DataSet
    {
        try {
            $config = $this->config;
            return $config->parser()->parse($jwt)->claims();
        } catch (Throwable $e) {
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }
    }

    /**
     * 验证令牌 (与 Go 端 JWT 中间件对齐)
     *
     * 必须校验 HMAC 签名 (SignedWith) 与有效期 (StrictValidAt):
     * 仅解析不验签的话, 攻击者可伪造任意 uid 的令牌通过鉴权。
     *
     * @param string $jwt
     * @return DataSet
     *
     * @throws RuntimeException 签名无效 / 已过期 / 声明缺失 / uid 非法时抛出 AuthInvalid
     */
    public function validatorToken($jwt): DataSet
    {
        $config = $this->config;
        try {
            $token = $config->parser()->parse((string)$jwt);
        } catch (Throwable $e) {
            // 结构不合法 (base64 解码失败等)
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }

        try {
            // 校验签名 (密钥不匹配即拒绝) 与有效期 (StrictValidAt: exp/iat/nbf 缺失或非法均拒绝)。
            // 注意: 5.x 的 validate() 只返回 bool 不抛异常, 必须用 assert() 才会抛出 RequiredConstraintsViolated。
            $config->validator()->assert(
                $token,
                new SignedWith($config->signer(), $config->signingKey()),
                new StrictValidAt(new SystemClock(new \DateTimeZone('UTC')))
            );
        } catch (Throwable $e) {
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }

        $claims = $token->claims();
        $uid = $claims->get('uid');
        // 严格类型检查: uid 必须是正整数 (与 Go/Rust/TS 端一致, 对齐 parse_token/GetJwtClaims/jose 的
        // 类型校验语义); 载荷为字符串/浮点等畸形值时一律拒绝, 而非 (int) 强转后放行
        if (!is_int($uid) || $uid <= 0) {
            throw new RuntimeException('', ResultCode::AuthInvalid->value);
        }
        return $claims;
    }

}
