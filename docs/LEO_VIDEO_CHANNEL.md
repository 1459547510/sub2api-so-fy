# Leo 同步视频渠道

Sub2API 的 `leo` 渠道用于连接二次开发的 LeoStudio Cookie 号池。每个 Sub2API Leo 账号对应一个 LeoStudio 服务实例，Sub2API 负责账号级调度、并发、代理、故障转移和计费，LeoStudio 继续负责其内部 Cookie 号池。

本版本保留同步视频生成兼容，同时新增基于 LeoStudio 异步协议的用户端视频工作台；不复用 Grok 渠道协议。

## LeoStudio 接口要求

LeoStudio 实例必须提供以下接口：

- `GET /health`
- `POST /v1/videos/generations`
- 所有请求使用 `Authorization: Bearer <api_key>`

在 LeoStudio 二次开发中配置 API Key，并确保错误 Key 对 `/health` 和视频生成接口返回 `401` 或 `403`。Base URL 必须精确指向 `/v1`，例如：

```text
http://127.0.0.1:18000/v1
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
3. 在分组中填写 480p、720p、1080p 三档 USD/秒回退价格。Leo 分组三档都必填，`0` 表示免费。
4. 在“渠道管理”中创建或编辑关联该分组的渠道，并进入 `Leo` 平台的模型定价。
5. 为已经配置计费的模型分别添加“视频（按秒）”定价。定价编辑器会按模型显示实际支持的分辨率：`happy-horse-1.1` 仅显示 720p、1080p，`grok-imagine-1.5` 仅显示 400p、544p、720p、960p，LTX 2.3 仅显示 1080p、1440p、2160p，Seedance 显示各自支持的档位。每条只能绑定一个精确模型，不能使用通配符；不要配置在账号统计定价规则中。
6. 创建平台为 `Leo`、类型为 `API Key` 的账号。
7. 填写 LeoStudio `/v1` Base URL、Bearer API Key、代理和并发。
8. 保留或调整默认模型映射：

```text
seedance-2.0      -> seedance-2.0
seedance-2.0-fast -> seedance-2.0-fast
seedance-2.0-mini -> seedance-2.0-mini
happy-horse-1.1   -> happy-horse-1.1
grok-imagine-1.5  -> grok-imagine-1.5
ltx-2.3-pro       -> ltxv-2.3-pro
ltx-2.3-fast      -> ltxv-2.3-fast
```

9. 执行账号连接测试，确认 `/health` 返回成功。

模型映射只改写请求 JSON 中的 `model`，提示词、分辨率、时长、音频和 guidance 字段原样转发。当前公开模型及其参数矩阵见 `docs/LEO_VIDEO_MODEL_SPECS.md`。同步和异步 Sub2 API 均使用同一套模型校验：

- `image_url`：兼容字段，表示单张首帧图片；
- `start_frame_url`、`end_frame_url`：分别表示首帧和尾帧；
- `image_urls`：图片参考 URL 数组，与 `guidances.image_reference` 合计最多 4 张；
- `guidances.start_frame`、`end_frame`、`image_reference`：使用上传接口返回的媒体 URL；
- `guidances.video_reference_base`：使用上传接口返回的 MP4/MOV 媒体 URL，URL 对象的 `type` 必须为 `UPLOADED`；
- `guidances.audio_reference`：使用上传接口返回的 MP3/WAV 媒体 URL，URL 对象的 `type` 必须为 `UPLOADED`，并且不在 URL 对象中提交客户端时长。

图片、视频和音频 URL 必须是视频服务端可以访问的绝对 `http`/`https` URL。生成接口不接受 data URL、Base64 或 multipart 媒体；参考音频必须同时搭配 `image_reference` 或 `video_reference_base`。Sub2 的 `/v1/videos/uploads` 支持图片、视频和音频三类上传，客户只需使用返回的 `media_url`。

首尾帧与参考图是两种独立输入模式：首帧和尾帧可以一起提交，但不能再附带 `image_urls` 或 `guidances.image_reference`；参考图请求也不能同时携带首帧或尾帧字段。

用户工作台提交单张图片时使用兼容性已验证的 `image_url` 首帧链路；提交多张图片时使用 `guidances.image_reference`，并按数组位置显式写入从 `0` 开始的连续 `order`。参考图模式与首尾帧模式互斥，服务端会再次校验该组合。

## 用户端菜单开关

管理员可在“系统设置 > 功能开关 > 视频生成”控制用户侧边栏中的“视频生成”菜单。设置键为 `video_generation_enabled`，默认开启；只有明确保存为 `false` 时才隐藏菜单，升级后无需重新开启。

该开关只控制侧边栏入口，不是访问控制：关闭后 `/video-generation` 路由仍可直接访问，Leo 同步与异步 API 也继续工作。需要停用视频能力时，应同时调整 Leo 分组、账号或 API Key 权限。

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
    "audio": false,
    "image_urls": [
      "https://media.example/reference-1.png",
      "https://media.example/reference-2.png"
    ]
  }'
```

该请求是同步调用，上游生成可能持续数分钟。客户端断开后，Sub2API 仍会等待 LeoStudio 完成，并持续占用所选账号的并发槽直至上游结束。

## 调度与计费

- Leo、OpenAI 和 Grok 账号分别调度，不会混选。
- `400` 和 `422` 视为请求错误，不切换账号，也不计费。
- `401`、`403`、`429`、`5xx` 和传输错误会尝试切换到其他可用 Leo 账号。
- 用户接口和任务记录中的错误统一使用 `Video provider`，不暴露 LeoStudio 或其底层供应商名称。
- 优先使用响应中的 `provider.resolution` 和 `provider.duration` 计费；缺失时回退到请求值。
- 精确命中的渠道模型视频定价优先于分组视频价格；渠道没有匹配的视频定价时，回退到分组的 480p、720p、1080p 价格。渠道的计费模型来源设置决定使用请求模型、渠道映射模型还是上游模型匹配价格。
- 费用为“对应分辨率 USD/秒单价 x 实际时长 x 视频数量 x 分组视频倍率”。三个 Seedance 模型均可配置独立单价，同步和异步请求使用相同优先级；Mini 按 720p 单价计费。

## 能力边界

本版本不支持 Leo 图片生成、视频编辑、视频扩展、Messages、Responses、Chat、Embeddings 或 WebSocket 接口。这些请求不会回退到 Grok 或 OpenAI 账号。

## 异步视频工作台

前端用户端新增 `/video-generation` 工作台。页面可以使用当前用户已保存的 Leo 分组 Sub2API API Key，也可以临时输入其他 Sub2API API Key；自定义 Key 只保存在当前页面内存，不写入浏览器存储或用户账号。切换 Key 来源或修改自定义 Key 会清空当前任务列表并停止旧轮询，提交、参考图上传、任务查询、取消以及视频预览下载统一使用当前 Key。页面不会显示 LeoStudio API Key、Cookie、上游账号或上游任务 ID。工作台支持：

- 文生视频；
- 按当前模型显示对应数量的远程 `http`/`https` 参考图片 URL；
- PNG、JPEG、WebP 本地参考图片，Seedance 最多 4 张、Happy Horse 最多 9 张，每张最大 10 MiB；
- 独立首帧和尾帧图片，可同时选择并并行上传，分别写入 `start_frame_url` 和 `end_frame_url`；首尾帧模式与参考图模式互斥；
- `pending`、`running`、`settling`、`completed`、`failed`、`canceled` 状态查询；
- 只有 `pending` 任务可取消，`running` 和终态任务不能取消；
- 完成任务的本地视频预览和下载。

### 工作台模型参数约束

工作台按当前 LeoStudio 模型能力注册表过滤参数；完整矩阵见
`docs/LEO_VIDEO_MODEL_SPECS.md`。模型切换后会重置分辨率、时长和画面比例，
提交前服务端会再次校验相同组合：

| 模型 | 可选分辨率 | 说明 |
| --- | --- | --- |
| `seedance-2.0` | `480p`、`720p`、`1080p` | 4–15 秒；1080p 为 4–12 秒。 |
| `seedance-2.0-fast` | `480p`、`720p` | 4–15 秒，不显示 1080p。 |
| `seedance-2.0-mini` | `480p`、`720p` | 两档均支持 `16:9`、`1:1`、`9:16`。 |
| `happy-horse-1.1` | `720p`、`1080p` | 3–15 秒，最多 9 张参考图；不支持尾帧、参考视频和参考音频。 |
| `grok-imagine-1.5` | `400p`、`544p`、`720p`、`960p` | 3–15 秒，必须提供一张首帧；不支持其他参考 guidance。 |
| `ltx-2.3-pro` | `1080p`、`1440p`、`2160p` | 固定 `16:9`；仅 6、8、10 秒；支持首尾帧、生成音频和提示词增强。 |
| `ltx-2.3-fast` | `1080p`、`1440p`、`2160p` | 固定 `16:9`；支持 6–20 秒偶数值；支持首尾帧、生成音频和提示词增强。 |

工作台会为每个模型设置一个可用的默认分辨率和时长；页面提交前会再次校验模型、分辨率、时长、画面比例和 guidance 组合，非法组合不会上传媒体或创建任务。直接调用公共 API 的客户端也会得到同样的服务端校验。

提交工作台任务时，浏览器请求 Sub2API：

```http
POST /v1/videos/generations
Prefer: respond-async
Authorization: Bearer <sub2_api_key>
Content-Type: application/json
```

Sub2API 返回 `202` 和自己的 `job_id`，后台协调器再轮询固定的 LeoStudio 账号。页面刷新后通过 `GET /v1/videos/jobs?limit=50` 恢复任务，不需要保持长连接。

工作台顶部提供“生成工作台 / API 对接文档”页签。独立子页面 `/video-generation/api-docs` 记录 Bearer 鉴权、本地图片上传、异步创建、单图、首尾帧与多参考图请求、两种图片模式的互斥约束、任务列表和详情、取消、成品下载、状态流转及常见 HTTP 错误；它不增加侧边栏一级菜单，也不展示供应商凭据、账号标识或供应商任务编号。

## 本地图片生命周期

本地图片保存在现有数据目录的 `video-inputs/` 子目录，文件名使用随机 token，权限为仅宿主机进程可读。LeoStudio 与 Sub2API 必须运行在同一宿主机，上传接口返回的图片 URL 使用 `127.0.0.1` 回环地址；跨主机、Docker 网络或多实例部署时不要猜测地址，应改用对象存储或显式内部地址方案。

内部读取接口只接受真实 loopback `RemoteAddr`，不信任 `X-Forwarded-For`，也不返回目录列表或服务器路径。Sub2 会识别 `image_url`、首尾帧、`image_urls` 和嵌套图片 guidance 中的全部本地 token。任务完成、失败或取消后，关联图片会一起标记为终态并至少保留 1 小时；没有关联任务的孤儿文件保留 24 小时。后台 runtime 启动时执行一次扫描，之后每天最多执行一次清理，删除失败会在下一轮重试。

使用 `embed` 标签构建内嵌前端时，`/internal/video-inputs/` 必须在 SPA fallback 之前交给后端路由处理；否则内部图片请求会收到 `200 text/html` 的前端页面，LeoStudio 会把 HTML 当作图片上传并导致视频创建失败。发布验收至少应确认：不存在的 token 返回 `404`，有效 token 返回实际图片 MIME，且响应大小与 SHA256 和上传源一致。

## 本地视频成品

LeoStudio 报告任务完成后，Sub2API 会先从结果中的 `source_url`、`mp4_url`、`video_url` 或 `url` 下载第一个视频。文件先写入临时文件，限制为最大 512 MiB，并校验 MP4 文件头 `ftyp`；校验通过后再原子移动到现有数据目录的 `video-outputs/<job_id>.mp4`。已存在且校验有效的文件会直接复用，服务重启后的结算重试不会重复下载。

任务结果中的原始 CDN 地址保存在 `source_url`，客户端使用的 `mp4_url`、`url` 和 `local_url` 会改写为 `/v1/videos/jobs/<job_id>/content`。该内容接口要求 `Authorization: Bearer <sub2_api_key>`，只允许创建该任务的 API Key 读取已完成任务，并通过 HTTP Range 返回 `video/mp4`。浏览器工作台先带 Bearer 下载 Blob，再使用本地 Blob URL 播放和下载，不直接播放 LeoStudio CDN 地址。

前端页面的 CSP 必须允许 `media-src 'self' blob:`，否则浏览器会成功下载成品但拒绝把 Blob URL 交给 `<video>` 播放。默认策略和安全中间件会自动补齐该指令，包括仍使用旧自定义 CSP 配置的部署；播放失败时页面仍保留原文件下载入口，重试会撤销旧 Blob 并重新获取成品。

当前不会自动清理 `video-outputs/` 中的成品；部署时应把 `pricing.data_dir` 所在磁盘容量纳入监控和备份策略。多实例部署必须共享同一数据目录，否则任务记录所在实例可能读取不到成品。

## 异步计费

任务提交前按请求分辨率、时长和提交时价格快照冻结预计余额。价格快照同时冻结渠道、计费模型、计费模型来源、渠道映射模型、定价来源、三档 USD/秒价格和视频倍率；任务运行期间修改渠道或分组价格不会改变该任务的结算依据。LeoStudio 报告完成后，只有视频下载、MP4 校验和本地保存全部成功，才按结果里的实际分辨率和时长确认一次用量并释放冻结额。结果没有视频 URL 时任务立即失败；下载或校验失败最多尝试 3 次，最终失败后释放冻结且不产生用量记录。上游生成失败或用户取消同样只释放冻结，不产生用量记录。

任务记录和结算标记持久化在 `video_jobs` 表，本地文件路径由 `job_id` 确定，Sub2API 重启后会继续保存和结算，不依赖浏览器在线。同步兼容接口也会拒绝 HTTP 成功但没有有效视频 URL 的上游响应，因此不会为这类空成功响应记录用量；同步接口返回的媒体仍使用 LeoStudio 原始响应，平台工作台使用上述异步本地成品链路。

成功结算会通过统一用量记录器写入一条视频使用记录，包含视频数量、实际分辨率、实际时长、费用以及提交时冻结的渠道和模型映射归属。使用自定义 Sub2API API Key 时，用量与扣费归属该 Key 对应的用户和分组，而不是当前网页登录账号。

## 运维边界

当前实现对应 LeoStudio 最新模型目录、媒体 URL guidance 与异步协议：`POST /v1/videos/generations`、`GET /v1/videos/jobs/:id` 和 `DELETE /v1/videos/jobs/:id`。Sub2API 将完成视频保存到本地数据目录并提供 API Key 鉴权读取；不提供 Webhook、SSE、WebSocket、编辑、扩展或任务删除功能。

### API key quota settlement

When a custom Sub2API API Key has a positive USD quota, a successfully settled Leo video also increments that key's `quota_used` through the shared atomic quota updater. Failed, canceled, or incomplete video jobs do not consume the key quota.

### Production account model mapping

The production Leo API-key account used by channel `5` must map every public
video model to the same upstream model name. The verified account mapping is:

```text
seedance-2.0      -> seedance-2.0
seedance-2.0-fast -> seedance-2.0-fast
seedance-2.0-mini -> seedance-2.0-mini
happy-horse-1.1   -> happy-horse-1.1
grok-imagine-1.5  -> grok-imagine-1.5
ltx-2.3-pro       -> ltxv-2.3-pro
ltx-2.3-fast      -> ltxv-2.3-fast
```

If a public model has channel pricing but is absent from the account mapping,
the scheduler returns `no available accounts` before creating a video job.
Adding a mapping exposes the model to `/v1/models`; it does not prove that the
LeoStudio upstream implementation is enabled, so a controlled generation test
is still required before customer traffic is enabled.

### Newly exposed models

`happy-horse-1.1` and `grok-imagine-1.5` are available in the video workbench and
public video API. Happy Horse supports `720p`/`1080p`, 3–15 seconds, up to nine
reference images, and prompt enhancement. Grok supports `400p`/`544p`/`720p`/`960p`,
3–15 seconds, and requires one start frame. Their unsupported guidance inputs
are rejected before dispatch. `seedance-2.0-mini` now supports both `480p` and
`720p`, with `16:9`, `1:1`, and `9:16` at each resolution.

`ltx-2.3-pro` and `ltx-2.3-fast` are also exposed. Both use `1080p` as the
default, accept `1440p` and `2160p`, and are fixed to `16:9`. Pro accepts
6/8/10 seconds; Fast accepts even durations from 6 through 20 seconds. Both
support a start/end frame pair, generated audio, and prompt enhancement, but
reject image, video, and audio references. Existing Leo accounts must add both
exact model mappings before these models can be scheduled; newly created
accounts include them by default.

Channel model pricing entries use each model's native resolution tiers:

- Happy Horse: `720p = 0.15` and `1080p = 0.19` USD/s.
- Grok Imagine 1.5: `400p = 0.10`, `544p = 0.10`, `720p = 0.17`, and
  `960p = 0.17` USD/s.
- Seedance models retain their model-specific `480p`/`720p`/`1080p` subset.
- LTX 2.3 Pro/Fast: configure independent `1080p`, `1440p`, and `2160p`
  channel prices. Reserve and settlement use the exact request/result tier;
  LTX pricing does not fall back to the group's legacy 480p/720p/1080p table.

### Production pricing snapshot

As of 2026-07-24, production channel `Seedance 2 视频专用渠道` (channel ID `5`) uses the following USD-per-second prices:

| Model | 480p | 720p | 1080p |
| --- | ---: | ---: | ---: |
| `seedance-2.0` | `$0.12` | `$0.25` | `$0.60` |
| `seedance-2.0-fast` | `$0.10` | `$0.20` | `$0.25`* |
| `seedance-2.0-mini` | - | `$0.17` | - |

`*` The Fast 1080p price remains configured because Seedance Fast uses the legacy three-tier pricing configuration. The Fast model capability matrix does not expose 1080p, so this tier is not advertised to users and cannot be selected in the video workbench. Exact channel pricing takes precedence over the unchanged group-level fallback prices; LTX uses its native three-tier channel pricing listed below.

As of 2026-07-29, the LTX channel prices are:

| Model | 1080p | 1440p | 2160p |
| --- | ---: | ---: | ---: |
| `ltx-2.3-fast` | `$0.06` | `$0.21` | `$0.24` |
| `ltx-2.3-pro` | `$0.09` | `$0.18` | `$0.36` |

Each LTX tier is frozen independently in the async job billing snapshot, so a
later channel edit does not change an already accepted job's settlement price.

### Customer media upload contract

`POST /v1/videos/uploads` accepts `image`, `video`, `audio`, or `file` with
`media_type`. Existing `image` clients remain compatible. The response exposes
an opaque `media_url` for use in the generation request; clients do not need to
handle upstream identifiers.

- Images: PNG, JPEG, or WebP, up to 10 MiB.
- Videos: the filename must end in MP4 or MOV, the bytes must be a readable ISO
  Base Media container with an `ftyp` header, the file is limited to 100 MiB,
  and a generation may contain at most three video references.
- Audio: MP3 with readable frames or RIFF/WAVE PCM 16/24-bit WAV, up to 15 MiB,
  between 2 and 30 seconds, with at most one audio reference.

Uploaded media is placed in `guidances.video_reference_base` or
`guidances.audio_reference` with `type: "UPLOADED"`. An audio reference must
be paired with an image reference or a video reference. Invalid extensions,
containers, WAV encoding, size, and audio duration are rejected before job
creation.
