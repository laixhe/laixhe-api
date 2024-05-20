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

#### JWT 的使用方式

> 放在 HTTP 请求的头信息Authorization字段里面

```
Authorization: Bearer <token>
```
