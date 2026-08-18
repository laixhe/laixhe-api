<?php

namespace Tests\Unit;

use App\Utils\JwtUtil;
use DateTimeImmutable;
use Lcobucci\JWT\Configuration;
use Lcobucci\JWT\Signer\Hmac\Sha256;
use Lcobucci\JWT\Signer\Key\InMemory;
use RuntimeException;
use Tests\TestCase;

/**
 * JWT 工具单元测试 (纯逻辑, 无需数据库)
 */
class JwtUtilTest extends TestCase
{
    public function test_create_and_parse_token(): void
    {
        $jwt = JwtUtil::getInstance()->createToken(1);
        $this->assertIsString($jwt);
        $this->assertSame(1, (int) JwtUtil::getInstance()->parseToken($jwt)->get('uid'));
    }

    public function test_validator_accepts_valid_token(): void
    {
        $jwt = JwtUtil::getInstance()->createToken(2);
        $claims = JwtUtil::getInstance()->validatorToken($jwt);
        $this->assertSame(2, (int) $claims->get('uid'));
    }

    public function test_validator_rejects_tampered_payload(): void
    {
        $jwt = JwtUtil::getInstance()->createToken(3);
        $parts = explode('.', $jwt);
        // 篡改 payload 再重新 base64url 编码, 签名必然与内容不匹配
        $decoded = base64_decode(strtr($parts[1], '-_', '+/'));
        $decoded[0] = $decoded[0] === '1' ? '0' : '1';
        $parts[1] = rtrim(strtr(base64_encode($decoded), '+/', '-_'), '=');
        $tampered = implode('.', $parts);

        $this->expectException(RuntimeException::class);
        JwtUtil::getInstance()->validatorToken($tampered);
    }

    public function test_validator_rejects_uid_zero(): void
    {
        // 对齐 Go 端语义: uid<=0 视为伪造 token, 一律拒绝
        $jwt = JwtUtil::getInstance()->createToken(0);
        $this->expectException(RuntimeException::class);
        JwtUtil::getInstance()->validatorToken($jwt);
    }

    public function test_validator_rejects_garbage_token(): void
    {
        $this->expectException(RuntimeException::class);
        JwtUtil::getInstance()->validatorToken('not-a-jwt');
    }

    public function test_validator_rejects_non_integer_uid(): void
    {
        // 载荷 uid 为字符串 (如 "5"): 严格类型检查应拒绝 (与 Go/Rust/TS 端一致,
        // 而非旧版 (int) 强转后放行)
        $config = Configuration::forSymmetricSigner(
            new Sha256(),
            InMemory::plainText((string) config('jwt.secret'))
        );
        $now = new DateTimeImmutable();
        $token = $config->builder()
            ->withClaim('uid', '5')
            ->expiresAt($now->modify('+3600 second'))
            ->issuedAt($now)
            ->getToken($config->signer(), $config->signingKey())
            ->toString();

        $this->expectException(RuntimeException::class);
        JwtUtil::getInstance()->validatorToken($token);
    }
}
