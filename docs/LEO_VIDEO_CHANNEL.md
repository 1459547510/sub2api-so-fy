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

本版本不支持 Leo 图片生成、视频编辑、视频扩展、Messages、Responses、Chat、Embeddings 或 WebSocket 接口。这些请求不会回退到 Grok 或 OpenAI 账号。

## 异步视频工作台

前端用户端新增 `/video-generation` 工作台。页面只接触用户选择的 Sub2API API Key，不会显示 LeoStudio API Key、Cookie、上游账号或上游任务 ID。工作台支持：

- 文生视频；
- 远程 `http`/`https` 图片 URL；
- PNG、JPEG、WebP 本地图片上传，最大 10 MiB；
- `pending`、`running`、`settling`、`completed`、`failed`、`canceled` 状态查询；
- 只有 `pending` 任务可取消，`running` 和终态任务不能取消；
- 完成任务的视频预览、打开和下载。

提交工作台任务时，浏览器请求 Sub2API：

```http
POST /v1/videos/generations
Prefer: respond-async
Authorization: Bearer <sub2_api_key>
Content-Type: application/json
```

Sub2API 返回 `202` 和自己的 `job_id`，后台协调器再轮询固定的 LeoStudio 账号。页面刷新后通过 `GET /v1/videos/jobs?limit=50` 恢复任务，不需要保持长连接。

## 本地图片生命周期

本地图片保存在现有数据目录的 `video-inputs/` 子目录，文件名使用随机 token，权限为仅宿主机进程可读。LeoStudio 与 Sub2API 必须运行在同一宿主机，上传接口返回的图片 URL 使用 `127.0.0.1` 回环地址；跨主机、Docker 网络或多实例部署时不要猜测地址，应改用对象存储或显式内部地址方案。

内部读取接口只接受真实 loopback `RemoteAddr`，不信任 `X-Forwarded-For`，也不返回目录列表或服务器路径。任务完成、失败或取消后文件会被标记为终态，至少保留 1 小时；没有关联任务的孤儿文件保留 24 小时。后台 runtime 启动时执行一次扫描，之后每天最多执行一次清理，删除失败会在下一轮重试。

## 异步计费

任务提交前按请求分辨率、时长和提交时价格快照冻结预计余额。完成任务按 LeoStudio 结果里的实际分辨率和时长只结算一次，并释放全部冻结额；失败或取消只释放冻结，不产生用量记录。任务记录和结算标记持久化在 `video_jobs` 表，Sub2API 重启后会继续查询和结算，不依赖浏览器在线。

## 运维边界

当前实现对应 LeoStudio 上游提交 `f822735629c51f15d115e3e60b161a93ec2e20ff` 的异步协议：`POST /v1/videos/generations`、`GET /v1/videos/jobs/:id` 和 `DELETE /v1/videos/jobs/:id`。第一版不保存生成视频到本地，不提供 Webhook、SSE、WebSocket、编辑、扩展或任务删除功能。
