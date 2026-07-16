# Leo 同步视频渠道

Sub2API 的 `leo` 渠道用于连接二次开发的 LeoStudio Cookie 号池。每个 Sub2API Leo 账号对应一个 LeoStudio 服务实例，Sub2API 负责账号级调度、并发、代理、故障转移和计费，LeoStudio 继续负责其内部 Cookie 号池。

本版本仅支持同步视频生成，不复用 Grok 渠道协议，也不包含异步任务系统。

## LeoStudio 接口要求

LeoStudio 实例必须提供以下接口：

- `GET /health`
- `POST /v1/videos/generations`
- 所有请求使用 `Authorization: Bearer <api_key>`

在 LeoStudio 二次开发中配置 API Key，并确保错误 Key 对 `/health` 和视频生成接口返回 `401` 或 `403`。Base URL 必须精确指向 `/v1`，例如：

```text
http://leostudio:8000/v1
```

Sub2API 的连接测试会请求去掉 `/v1` 后的 `/health`；视频生成会请求 `<base_url>/videos/generations`。

成功的视频响应应保留 LeoStudio 的 `data` 和 `provider` 字段。Sub2API 使用以下输出字段记录实际计费元数据：

```json
{
  "data": [
    { "mp4_url": "https://cdn.example/video.mp4" }
  ],
  "provider": {
    "generation_id": "generation-id",
    "resolution": "RESOLUTION_720",
    "duration": 8
  }
}
```

## Sub2API 配置

1. 在管理端创建平台为 `Leo` 的分组。
2. 开启“允许当前分组生成视频”。
3. 手动填写 480p、720p、1080p 三档 USD/秒价格。三档都必填，`0` 表示免费。
4. 创建平台为 `Leo`、类型为 `API Key` 的账号。
5. 填写 LeoStudio `/v1` Base URL、Bearer API Key、代理和并发。
6. 保留或调整默认模型映射：

```text
seedance-2.0      -> seedance-2.0
seedance-2.0-fast -> seedance-2.0-fast
```

7. 执行账号连接测试，确认 `/health` 返回成功。

模型映射只改写请求 JSON 中的 `model`，提示词、分辨率、时长、音频等字段原样转发。

## 调用示例

```bash
curl -X POST "$SUB2_BASE_URL/v1/videos/generations" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "seedance-2.0",
    "prompt": "A cinematic city at night",
    "aspect_ratio": "16:9",
    "resolution": "720p",
    "duration": 8,
    "audio": false
  }'
```

该请求是同步调用，上游生成可能持续数分钟。客户端断开后，Sub2API 仍会等待 LeoStudio 完成，并持续占用所选账号的并发槽直至上游结束。

## 调度与计费

- Leo、OpenAI 和 Grok 账号分别调度，不会混选。
- `400` 和 `422` 视为请求错误，不切换账号，也不计费。
- `401`、`403`、`429`、`5xx` 和传输错误会尝试切换到其他可用 Leo 账号。
- 优先使用响应中的 `provider.resolution` 和 `provider.duration` 计费；缺失时回退到请求值。
- 费用为“对应分辨率 USD/秒价格 x 实际时长 x 视频倍率”。

## 能力边界

本版本不支持 Leo 图片生成、视频编辑、视频扩展、生成状态查询、Messages、Responses、Chat、Embeddings 或 WebSocket 接口。这些请求不会回退到 Grok 或 OpenAI 账号。
