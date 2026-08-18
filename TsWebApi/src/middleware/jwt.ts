// JWT 令牌工具
// 提供 Token 签发、校验、载荷提取

import { jwtVerify, SignJWT } from "jose";
import config from "../config";
import { info, warn } from "../util/logger";

// JwtPayload JWT 令牌载荷，存储用户 UID
interface JwtPayload {
  uid: number;
}

// generateToken 创建 JWT 载荷并设置过期时间、发布时间
// 注意：TS 版本不设置 NotBefore（生效时间），与 Go 略有差异；
// 本项目校验端是 jose 的 jwtVerify（见本文件 getJwtClaims），未要求 nbf 字段，因此不影响功能
export async function generateToken(uid: number): Promise<string> {
  const token = await new SignJWT({ uid })
    .setProtectedHeader({ alg: "HS256" })
    .setIssuedAt() // 发布时间（创建时间）
    .setExpirationTime(
      // 过期时间（jose v6 需要绝对时间）
      new Date(Date.now() + config.jwt.expireTime * 1000)
    )
    .sign(config.jwt.secretKey);
  info("jwt", "Token签发成功", { uid });
  return token;
}

// 从请求头中解析 JWT Token 并返回 payload
async function getJwtClaims(
  authHeader: string | null
): Promise<JwtPayload | null> {
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    return null;
  }

  // 截掉 "Bearer " 前缀（7 个字符），仅保留 Token 本身
  const token = authHeader.slice(7);
  try {
    const { payload } = await jwtVerify(token, config.jwt.secretKey);
    // 防御畸形载荷: uid 必须是正数, 否则视为无效 (对齐 Go 端 uid<=0 视为伪造的语义)
    const uid = payload.uid;
    if (typeof uid !== "number" || !Number.isInteger(uid) || uid <= 0) {
      warn("jwt", "Token载荷uid非法", { uid });
      return null;
    }
    return { uid };
  } catch (err) {
    warn("jwt", "Token验签失败", { error: String(err) });
    return null;
  }
}

// 从 headers 对象中获取 JWT claims
export async function getJwtClaimsFromHeaders(
  headers: Record<string, string | undefined>
): Promise<JwtPayload | null> {
  const authHeader = headers["authorization"] || null;
  return getJwtClaims(authHeader);
}
