### JWT 的数据结构

> 它是一个很长的字符串，中间用点 . 分隔成三个部分 ( Header.Payload.Signature )

- Header（头部）
- Payload（负载）
- Signature（签名）

##### Header

> 是一个 JSON 对象，描述 JWT 的元数据，使用 Base64URL 算法转成字符串

```
{
  "alg": "HS256",
  "typ": "JWT"
}
```

在代码中，`alg`属性表示签名的算法（algorithm），默认是 HMAC SHA256（写成 HS256）；`typ`属性表示这个令牌（token）的类型（type），JWT 令牌统一写为JWT

##### Payload

> 一个 JSON 对象，用来存放实际需要传递的数据。JWT 规定了 7 个官方字段，也可以私有字段
> 使用 Base64URL 算法转成字符串

- iss (issuer)：签发人
- exp (expiration time)：过期时间
- sub (subject)：主题
- aud (audience)：受众
- nbf (Not Before)：生效时间
- iat (Issued At)：签发时间
- jti (JWT ID)：编号

##### Signature

> 是对前两部分的签名，防止数据篡改
> 需要指定一个密钥（secret）。这个密钥只有服务器才知道，不能泄露给用户。然后，使用 Header 里面指定的签名算法（默认是 HMAC SHA256）

```
HMACSHA256(base64UrlEncode(header) + "." +base64UrlEncode(payload), secret)
```

#### Base64URL

> Base64URL 是 Base64 的 URL 安全变体：把标准 Base64 中的 `+` 换成 `-`、`/` 换成 `_`，并去掉末尾的 `=` 填充符。这样编码结果可安全地出现在 URL 查询参数、HTTP 请求头等场景中，无需再转义。

#### JWT 的使用方式

> 放在 HTTP 请求的头信息Authorization字段里面

```
Authorization: Bearer <token>
```

## 本仓库中的实际使用

> 仓库内六种语言实现（Go / Java / PHP / Python / Rust / TypeScript）的 JWT 约定完全一致，可对照各自源码阅读：

- **签名算法**：HS256（HMAC-SHA256）。密钥为各端配置文件中的默认值，**生产环境务必更换**（位置见仓库 README「安全提醒」）。
- **载荷（Payload）**：只放 `uid`（用户 id）和三个注册声明 `iat`（签发时间）、`nbf`（生效时间）、`exp`（过期时间，默认 30 天 / 2592000 秒），不使用 `iss` / `sub` / `aud` / `jti`（注：仅 TS 版不设置 `nbf`，各端校验器均不要求该字段）：

```json
{
  "uid": 1,
  "iat": 1780000000,
  "nbf": 1780000000,
  "exp": 1782592000
}
```

- **校验**：`uid` 从 1 开始计数，`uid <= 0` 视为无效（防御伪造 `{"uid":0}` 的令牌）；签名与过期校验由各端 JWT 库完成。
- **传输**：携带在请求头 `Authorization: Bearer <token>`，由各端 JWT 中间件统一解析。
- **刷新**：`POST /api/v1/auth/refresh`（需 Bearer）从已验证的 token 中读取 `uid`，确认用户仍存在且状态正常后签发新 token。

