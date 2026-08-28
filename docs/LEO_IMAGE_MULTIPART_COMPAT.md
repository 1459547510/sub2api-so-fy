# Leo 图片 multipart 兼容

Sub2API 的 Leo 图片渠道支持客户继续使用 OpenAI 兼容的
`POST /v1/images/edits` multipart 上传。Sub2API 会保存参考图、生成临时的
HTTPS 地址，并把请求转换为 LeoStudio 支持的 JSON `image_urls` 请求；客户不需要
先把图片上传到外部图床。

## 配置

生产环境必须配置外部可访问的 Sub2API HTTPS 根地址：

```yaml
gateway:
  leo_image_upload_public_base_url: "https://api.example.com"
```

该地址不能包含用户名、密码、查询参数或片段。反向代理必须把
`/media/image-inputs/*` 转发到 Sub2API。配置缺失或无效时，普通 JSON URL 请求
继续可用，Leo multipart 请求返回 `503`。

## 请求兼容

multipart 请求中的 `image` 和 `image[...]` 文件字段会按原顺序转换为
`image_urls`。以下标准字段会保留：

- `model`、`prompt`、`n`、`size`、`response_format`
- `quality`、`background`、`output_format`、`moderation`
- `input_fidelity`、`style`、`stream`
- `output_compression`、`partial_images`

参考图只接受 PNG、JPEG 和 WebP，每张最大 10 MiB。`mask` 仍不受 Leo 图片渠道
支持；包含 mask 的请求会在 Sub2API 参数校验阶段返回 `400`。

转换只在最终路由平台为 Leo 时执行，其他图片渠道及 JSON `image_urls` /
`images[].image_url` 调用链不变。

## 临时文件和安全边界

- 临时文件使用 192 位随机令牌作为文件名和访问凭据，不暴露原始文件名。
- 公开读取仅允许服务端识别为图片的内容，并返回
  `X-Content-Type-Options: nosniff` 和 `Cache-Control: no-store`。
- 请求结束后文件被标记为终态，现有清理任务在一小时后删除；未完成标记的孤儿
  文件最多保留 24 小时。
- 公开地址必须使用 HTTPS。令牌 URL 不应写入客户错误响应、业务日志或监控标签。

## 验收

部署后使用非敏感测试图调用：

```bash
curl --fail-with-body \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -F "model=GPT Image-2" \
  -F "prompt=multipart compatibility validation" \
  -F "n=1" \
  -F "size=1024x1024" \
  -F "image=@reference.png;type=image/png" \
  "https://api.example.com/v1/images/edits"
```

验收需同时确认 Sub2API 返回 `200`、请求日志记录 `multipart=true`、Leo 账号被正常
选中、LeoStudio 产生成功记录且结果图片可下载。还需再执行一次 JSON URL 请求，
确认原调用链没有回归。

## 回滚

源码回滚使用本功能提交的父提交。生产回滚时恢复部署前的 Sub2API 二进制和
`config.yaml`，重启服务；保留新增配置不会影响旧版本读取其他配置项。
