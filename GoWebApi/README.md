### API 接口
```
http://127.0.0.1:6600/api/v1/swagger.yaml
http://127.0.0.1:6600/api/v1/swagger.json
```

#### 文件上传文档注释
```
// @Accept   mpfd
// @Produce  image/png
// @Param    file  formData  file  true  "图片文件"
// @Success  200   {file}    binary  "返回图片资源"

@Accept mpfd        表示文件上传，接收 multipart/form-data 格式的请求
@Produce image/png  表示返回图片二进制流，也可用 octet-stream 表示通用二进制
```
