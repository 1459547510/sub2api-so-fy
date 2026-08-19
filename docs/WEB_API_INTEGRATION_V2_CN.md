# 视频图片统一 API 对接文档 V2

> 版本：REST API v2
> 适用范围：拥有对应分组权限的 Sub2API Key
> 兼容说明：V1 视频接口继续有效，已有客户端无需修改调用路径

## 1. 接入约定

### Base URL

将 `https://your-sub2-domain.example` 替换为实际站点地址：

```text
https://your-sub2-domain.example
```

所有接口都使用同一个 Bearer API Key：

```http
Authorization: Bearer <SUB2_API_KEY>
```

视频创建请求建议声明异步响应：

```http
Content-Type: application/json
Prefer: respond-async
```

客户端只需要提交公开模型 ID 和标准字段，其余选择由服务端完成。

## 2. V1 兼容性

V2 在原有 V1 视频契约上增加图片接口和模型发现，以下行为保持不变：

- `POST /v1/videos/generations` 的请求字段和 `202 Accepted` 异步响应保持不变。
- `POST /v1/videos/uploads`、任务查询、任务列表、内容下载和取消路径保持不变。
- 原有 Bearer API Key 继续使用。Key 获得新分组权限后，旧客户端无需改代码即可调用新模型。
- 模型调度、权限和计费由服务端处理，客户端不需要感知实际服务实现。

V2 的兼容边界：

- 视频模型必须遵守本页的异步任务、轮询和内容下载契约。
- 图片模型必须遵守 OpenAI Images 的生成和编辑请求/响应契约。
- 未在当前 Key 的模型列表中出现的模型 ID 不应提交；服务端会返回模型或参数错误。

## 3. 模型发现

### `GET /v1/models`

返回当前 API Key 可用的公开模型 ID。客户端启动和定期刷新时都应重新读取，不要把模型清单永久写死在客户端。

```bash
curl "$BASE_URL/v1/models" \
  -H "Authorization: Bearer $SUB2_API_KEY"
```

响应示例：

```json
{
  "object": "list",
  "data": [
    { "id": "<video-model-id>", "type": "model" },
    { "id": "<image-model-id>", "type": "model" }
  ]
}
```

只将 `data[].id` 原样填写到后续请求的 `model` 字段。模型 ID 只用于选择可用模型，不能用来拼接其他接口路径。

### 当前模型价格

Web 端 V2 文档页会在打开时读取当前账号可见的媒体分组和模型目录，并显示图片、视频模型的当前价格。价格列表不是静态写在文档代码中的：

- 管理员修改分组或模型定价后，刷新价格按钮会立即重新读取并展示新价格。
- 页面在打开期间会定期刷新价格；页面只展示当前账号有权限使用的分组和模型。
- 图片价格按图片档位显示；视频价格按模型和分辨率显示，视频价格单位为美元/秒。
- 页面展示值包含当前分组生效倍率，最终扣费仍以请求实际命中的分组、模型、分辨率和时长为准。

客户端不应从文档页面抓取价格或把价格写死；应以管理端配置和实际账单结果为准。

## 4. 图片接口

图片接口遵循 OpenAI Images 格式。部分模型支持 JSON URL 和 multipart；部分模型只接受可访问的 HTTP(S) 参考图 URL。具体模型能力以 `/v1/models` 和服务端校验结果为准。

### 4.1 图片生成

```http
POST /v1/images/generations
```

常用字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `model` | string | 否 | 图片模型 ID；可省略使用服务端默认值，建议显式传入 `/v1/models` 返回的 ID |
| `prompt` | string | 是 | 生成说明 |
| `n` | integer | 否 | 生成数量，`1-10` |
| `size` | string | 否 | 例如 `1024x1024`，以模型支持范围为准 |
| `response_format` | string | 否 | `url` 或 `b64_json`，以模型支持范围为准 |
| `quality` | string | 否 | 可选质量档位，仅在模型支持时发送 |
| `image_urls` | string[] | 否 | 参考图绝对 HTTP(S) URL；也接受单数 `image_url` |

```bash
curl -X POST "$BASE_URL/v1/images/generations" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "<image-model-id>",
    "prompt": "A clean product photo on a white background",
    "n": 1,
    "size": "1024x1024",
    "response_format": "b64_json"
  }'
```

### 4.2 图片编辑

```http
POST /v1/images/edits
```

JSON 请求使用 `images[].image_url`，OpenAI 图片模型还可选 `mask.image_url`。`file_id` 不属于统一契约。

若模型不接受独立编辑请求，网关会把 `/v1/images/edits` 的 `images[].image_url` 转成 `image_urls`，再转发到 `/v1/images/generations`。这类模型不接受 `mask` 或 multipart；本地文件请先放到可访问的 HTTP(S) URL。

```bash
curl -X POST "$BASE_URL/v1/images/edits" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "<image-model-id>",
    "prompt": "Replace the background with a quiet studio",
    "images": [{ "image_url": "https://media.example/input.png" }],
    "mask": { "image_url": "https://media.example/mask.png" }
  }'
```

multipart 请求示例：

```bash
curl -X POST "$BASE_URL/v1/images/edits" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -F "model=<image-model-id>" \
  -F "prompt=Turn the sketch into a polished illustration" \
  -F "image[]=@./input.png" \
  -F "mask=@./mask.png"
```

图片响应示例：

```json
{
  "created": 1780000000,
  "data": [
    { "b64_json": "<base64-image>" }
  ]
}
```

## 5. 视频接口

### 5.1 创建任务

```http
POST /v1/videos/generations
```

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `model` | string | 是 | 从 `/v1/models` 读取的视频模型 ID |
| `prompt` | string | 是 | 场景、动作、镜头和风格描述 |
| `resolution` | string | 否 | 模型支持的分辨率，例如 `720p` |
| `duration` | integer | 否 | 模型支持的整数秒数 |
| `aspect_ratio` | string | 否 | 模型支持的画面比例，例如 `16:9` |
| `audio` | boolean | 否 | 是否请求生成音频，默认 `false` |
| `prompt_enhance` | string | 否 | 模型支持时可用的提示词增强档位 |
| `image_url` | string | 否 | 单张首帧绝对 HTTP(S) URL |
| `start_frame_url` | string | 否 | 首帧绝对 HTTP(S) URL |
| `end_frame_url` | string | 否 | 尾帧绝对 HTTP(S) URL |
| `guidances` | object | 否 | 参考图片、视频和音频对象 |

```bash
curl -X POST "$BASE_URL/v1/videos/generations" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: respond-async" \
  -d '{
    "model": "<video-model-id>",
    "prompt": "A slow aerial shot over a coastal city at sunrise",
    "resolution": "720p",
    "duration": 8,
    "aspect_ratio": "16:9",
    "audio": false
  }'
```

成功响应为 HTTP `202`：

```json
{
  "job_id": "vidjob_example",
  "status": "pending",
  "status_url": "/v1/videos/jobs/vidjob_example"
}
```

收到 `job_id` 后不要重复创建同一任务，继续轮询 `status_url`。

### 5.2 上传参考媒体

```http
POST /v1/videos/uploads
```

表单字段名为 `image`、`video` 或 `audio`。图片使用 PNG/JPEG/WebP，视频使用 MP4/MOV，音频使用 MP3 或 PCM WAV。文件大小、编码、时长和数量以所选模型限制为准。

```bash
curl -X POST "$BASE_URL/v1/videos/uploads" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -F "image=@./reference.png"
```

响应示例：

```json
{
  "upload_id": "upload_example",
  "media_url": "https://media.example/uploaded/reference.png",
  "media_type": "image",
  "content_type": "image/png",
  "size": 428516
}
```

把 `media_url` 放入对应的参考对象：

```json
{
  "guidances": {
    "image_reference": [
      {
        "image": { "url": "<MEDIA_URL>", "type": "UPLOADED" },
        "strength": "MID",
        "order": 0
      }
    ],
    "video_reference_base": [
      { "video": { "url": "<MEDIA_URL>", "type": "UPLOADED" } }
    ],
    "audio_reference": [
      { "audio": { "url": "<MEDIA_URL>", "type": "UPLOADED" } }
    ]
  }
}
```

首尾帧、参考图片、参考视频和参考音频能否组合由模型能力决定。创建 JSON 不接受 data URL、Base64 或 multipart 媒体。

### 5.3 任务查询、下载和取消

```http
GET    /v1/videos/jobs?limit=50&status=running
GET    /v1/videos/jobs/{job_id}
GET    /v1/videos/jobs/{job_id}/content
DELETE /v1/videos/jobs/{job_id}
```

任务状态：

| 状态 | 含义 |
| --- | --- |
| `pending` | 已接收，等待执行 |
| `running` | 正在生成 |
| `settling` | 正在保存结果 |
| `completed` | 结果可下载 |
| `failed` | 任务失败 |
| `canceled` | 任务已取消 |

```bash
curl "$BASE_URL/v1/videos/jobs/vidjob_example" \
  -H "Authorization: Bearer $SUB2_API_KEY"

curl -o output.mp4 "$BASE_URL/v1/videos/jobs/vidjob_example/content" \
  -H "Authorization: Bearer $SUB2_API_KEY"
```

只有创建任务的同一个 API Key 才能查询、下载或取消该任务。任务进入终态后不要继续提交取消请求。

### 5.4 OpenAI 视频兼容入口

面向使用 OpenAI 视频接口格式的中转/聚合网关（渠道类型选 OpenAI，`base_url` 填本站地址，密钥填本站 API Key）。原生 JSON 契约不受影响，两种方言按路径和 Content-Type 区分。

创建（`multipart/form-data`）：

```http
POST /v1/videos
```

| 表单字段 | 对应原生字段 | 说明 |
| --- | --- | --- |
| `model` | `model` | 公开模型 ID 原样填写 |
| `prompt` | `prompt` | 文本提示词 |
| `seconds` | `duration` | 正整数秒 |
| `size` | `resolution` + `aspect_ratio` | `1920x1080` 这类宽x高，或直接填 `720p` 等档位 |
| `input_reference` | 参考图 | PNG/JPEG/WebP 文件，自动上传并作为参考图片 |
| `metadata` | 高级字段 | JSON 字符串，可携带 `resolution`、`aspect_ratio`、`audio`、`duration`、`image_urls`、`start_frame_url`、`end_frame_url`、`guidances`、`prompt_enhance` |

```bash
curl "$BASE_URL/v1/videos" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -F "model=<video-model-id>" \
  -F "prompt=A slow aerial shot over a coastal city at sunrise" \
  -F "seconds=8" \
  -F "size=1280x720"
```

创建响应为 OpenAI 视频对象：

```json
{
  "id": "vidjob_example",
  "object": "video",
  "model": "<video-model-id>",
  "status": "queued",
  "progress": 0,
  "created_at": 1755500000,
  "seconds": "8",
  "size": "1280x720"
}
```

查询与下载：

```http
GET /v1/videos/{id}
GET /v1/videos/{id}/content
```

`GET /v1/videos/{id}` 返回同一视频对象，状态依次为 `queued`、`in_progress`、`completed`；失败或取消为 `failed`，并附 `error.message`。`completed` 后从 `/content` 下载 MP4。

兼容边界：

- 分辨率、时长、画面比例仍以所选模型能力为准，超出会返回参数错误。
- 多参考图/视频/音频、首尾帧等高级参数通过 `metadata` 传入，或改用 5.1 的原生 JSON 接口。
- 该入口始终异步；不支持通过 `DELETE /v1/videos/{id}` 取消，取消请使用 5.3 的任务接口。

## 6. 错误和重试

统一错误格式：

```json
{
  "error": {
    "type": "invalid_request_error",
    "message": "Detailed error message"
  }
}
```

| HTTP | 含义 | 建议 |
| --- | --- | --- |
| `400` | JSON、必填字段或参数组合无效 | 修正请求后再发 |
| `401` | API Key 缺失或无效 | 检查 Key |
| `402` | 余额不足 | 充值后重试 |
| `403` | 当前 Key 未开放该能力 | 检查分组权限 |
| `404` | 模型、任务或成品不存在 | 检查 ID 和 Key |
| `409` | 当前任务状态不允许操作 | 先查询状态 |
| `422` | 参考媒体或模型参数校验失败 | 按错误信息修改参数 |
| `429` | 频率或并发超限 | 使用指数退避 |
| `502` / `5xx` | 服务暂时不可用 | 限次重试 |

建议：

1. 视频创建使用异步模式，收到 `job_id` 后每 3-5 秒轮询一次。
2. `400`、`401`、`403`、`422` 不要原样重复提交，先修正请求或权限。
3. `429`、`502` 和网络错误最多重试 3 次，间隔 2、5、10 秒。
4. 已收到 `job_id` 时只轮询原任务，不要重复创建。

## 7. 接入检查清单

- API Key 已获得目标分组权限并有足够余额。
- 客户端启动时读取 `/v1/models` 并保存公开模型 ID。
- 图片调用使用 `/v1/images/generations` 或 `/v1/images/edits`。
- 视频调用使用 `/v1/videos/generations` 并携带 `Prefer: respond-async`。
- 本地图片、视频和音频先上传，再使用返回的 `media_url`。
- 只使用模型允许的分辨率、时长、画面比例和参考媒体组合。
- 错误重试遵守状态码和退避规则。
