# Leo 异步视频工作台设计

## 1. 目标

在 Sub2API 用户端新增可直接使用的 Leo 视频生成工作台。用户选择自己的 Leo API Key 后，可以提交纯文本视频、图片 URL 视频或本地图片视频，并在页面内查看多个异步任务的排队、生成、完成、失败和取消状态。

Sub2API 继续负责用户鉴权、分组权限、LeoStudio 实例调度、任务归属、余额冻结、实际输出计费和用量记录；LeoStudio 继续负责 Cookie 号池、持久化生成队列和 Leonardo 视频生成。

成功标准：

- 浏览器不接触 LeoStudio API Key、Cookie 或管理接口；
- 客户端提交后立即得到 Sub2API 任务 ID，不保持数分钟长连接；
- 多任务可连续提交，页面刷新、Sub2API 重启后仍可恢复；
- 文生视频、图片 URL 和本地图片三种输入均可用；
- 完成任务按 LeoStudio 实际分辨率和时长只计费一次；
- 失败或取消任务不计费并释放冻结余额；
- 本地输入图片不会无限占用磁盘，也不会公开暴露为可枚举资源。

## 2. 已确认决策

- 使用 LeoStudio `f822735629c51f15d115e3e60b161a93ec2e20ff` 已实现的异步协议。
- 同步 `POST /v1/videos/generations` 保持兼容；仅携带 `Prefer: respond-async` 时进入异步流程。
- Sub2API 持久化任务映射，不让浏览器直连 LeoStudio。
- 新增 `video_jobs` 数据表及数据库迁移已经用户批准。
- 客户端使用双栏工作台：左侧生成设置，右侧最新结果与任务列表；移动端改为上下布局。
- 支持连续提交多个任务，由 LeoStudio 队列和号池控制实际执行。
- 支持纯文本、远程图片 URL、本地图片上传。
- 第一版使用宿主机本地磁盘保存临时输入图片，不接入对象存储。
- Sub2API 与 LeoStudio 第一版均直接运行在同一宿主机。
- 只有 `pending` 任务允许取消；`running` 和终态任务不能取消。

## 3. 非目标

- 不修改 LeoStudio 仓库；直接使用其现有异步接口。
- 不在 Sub2API 重建视频生成队列，不调用 LeoStudio 同步接口模拟异步。
- 不增加视频编辑、视频扩展、Webhook、WebSocket 或 SSE。
- 不保存生成视频到本地磁盘，不建设视频素材库。
- 不支持跨主机或 Docker 网络下的本地图片地址自动发现。
- 不支持取消已经进入 `running` 的上游生成。
- 不提供任务优先级、人工重排、批量删除或管理员代提交。

## 4. 总体架构

### 4.1 推荐方案

Sub2API 作为异步任务代理和账务源：

1. 浏览器使用用户选择的 Sub2API API Key 提交生成请求。
2. Sub2API 校验权限、请求参数和余额，选择一个 Leo 账号。
3. Sub2API 冻结预计费用，向所选 LeoStudio 发送 `Prefer: respond-async`。
4. LeoStudio 返回上游 `job_id` 后，Sub2API 保存任务与账号映射并返回自己的不可猜测任务 ID。
5. Sub2API 后台协调器持续查询活动任务，将状态和终态结果写入数据库。
6. 完成后按实际输出结算；失败或取消后释放冻结余额。
7. 客户端轮询 Sub2API，不直接访问 LeoStudio。

### 4.2 未采用方案

- 浏览器直连 LeoStudio：会暴露 Leo Key，绕过 Sub2API 调度、权限、计费和审计。
- Sub2API 自建生成队列并调用同步接口：与 LeoStudio 的持久化队列重复，重启恢复和取消语义更复杂。
- 仅在浏览器保存上游 job ID：无法可靠恢复账号归属，也无法保证后台结算。

## 5. 数据模型

新增 `video_jobs` 表，至少包含：

| 字段 | 用途 |
|---|---|
| `job_id` | Sub2API 生成的不可猜测公开任务 ID，唯一 |
| `user_id`、`api_key_id`、`group_id` | 用户归属、查询边界和结算归属 |
| `account_id` | 固定的 LeoStudio 实例，不在状态轮询时重新调度 |
| `upstream_job_id` | LeoStudio 返回的整数任务 ID |
| `status` | `pending`、`running`、`completed`、`failed`、`canceled`、`settling` |
| `requested_model`、`upstream_model` | 外部模型与映射后模型 |
| `prompt` | 任务恢复和客户端历史展示 |
| `resolution`、`duration`、`aspect_ratio`、`audio` | 请求规格快照 |
| `image_source` | `none`、`url` 或 `local` |
| `image_url` | 远程图片 URL，或本地回环临时 URL |
| `local_input_name` | 服务端生成的相对文件名，不保存用户路径 |
| `result` | LeoStudio 终态 `result` JSON |
| `error_message` | 脱敏后的公开错误 |
| `hold_amount`、`actual_cost` | 冻结金额和最终金额 |
| `billing_snapshot` | 三档视频单价、倍率和计费归属快照 |
| `request_hash` | 冻结与结算的稳定指纹 |
| `settled_at` | 幂等结算标记 |
| `created_at`、`updated_at`、`started_at`、`finished_at` | 生命周期时间 |

索引覆盖 `job_id`、`user_id + created_at`、`api_key_id + created_at`、`status + updated_at` 和 `account_id + upstream_job_id`。上游任务 ID 只在同一 Leo 账号内唯一，因此使用账号与上游 ID 的联合唯一约束。

任务记录保留用于账务与历史，不因客户端删除而物理删除。第一版不提供删除任务记录功能。

## 6. 公共 API

所有用户操作继续使用 Sub2API API Key：

```http
Authorization: Bearer <SUB2_API_KEY>
```

### 6.1 提交异步任务

```http
POST /v1/videos/generations
Prefer: respond-async
Content-Type: application/json
```

请求字段沿用同步 Leo 协议：

```json
{
  "model": "seedance-2.0",
  "prompt": "A paper kite crossing a bright coastal sky",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": false,
  "image_url": "http://127.0.0.1:8080/internal/video-inputs/random-token"
}
```

Sub2API 返回自己的任务标识，不向客户端暴露上游账号或上游 job ID：

```http
HTTP/1.1 202 Accepted
Preference-Applied: respond-async
Location: /v1/videos/jobs/vidjob_xxx
Cache-Control: no-store
```

```json
{
  "job_id": "vidjob_xxx",
  "status": "pending",
  "status_url": "/v1/videos/jobs/vidjob_xxx"
}
```

同步请求不携带 `Prefer: respond-async` 时继续走既有同步 handler，响应和计费不变。

### 6.2 上传本地图片

```http
POST /v1/videos/uploads
Content-Type: multipart/form-data
```

只接受一个 `image` 文件，支持 PNG、JPEG、WebP，最大 `10 MiB`。服务端根据文件内容识别 MIME，不信任文件扩展名。成功后返回供后续生成请求使用的内部 URL和 `upload_id`。

上传接口只保存文件，不冻结余额。若后续任务创建失败，文件由孤儿清理器回收。

### 6.3 查询任务

```http
GET /v1/videos/jobs?limit=50&status=running
GET /v1/videos/jobs/:job_id
```

API Key 只能查询由自身创建的任务。列表默认按 `created_at` 倒序，最大 `100` 条。详情在完成时包含与同步响应一致的 `result`；失败时仅包含脱敏 `error.message`。

### 6.4 取消任务

```http
DELETE /v1/videos/jobs/:job_id
```

只有 `pending` 可取消。Sub2API 向原 LeoStudio 账号发送取消请求，成功后持久化 `canceled` 并释放余额。`running` 或终态返回 `409`；未知任务或不属于当前 API Key 的任务统一返回 `404`。

所有异步响应均设置 `Cache-Control: no-store`。

## 7. 上游协议与账号固定

LeoStudio 协议为：

- 创建：`POST {base_url}/videos/generations`，携带 `Prefer: respond-async`；
- 查询：`GET {base_url}/videos/jobs/{upstream_job_id}`；
- 取消：`DELETE {base_url}/videos/jobs/{upstream_job_id}`；
- 鉴权：`Authorization: Bearer <leo_api_key>`。

任务创建成功后固定 `account_id`。查询和取消必须回到同一个 LeoStudio 实例，不能重新调度到其他账号。

创建阶段的 `401`、`403`、`429`、`5xx` 和传输失败仍可在没有收到 `202` 时按现有规则切换其他 Leo 账号。收到 `202` 后发生响应不确定性时不自动重复创建，避免双重消耗。状态查询失败只记录并重试，不切换账号。

## 8. 后台协调与恢复

新增一个轻量协调器，不执行生成，只同步上游状态：

- 每 `2` 秒批量领取需要轮询的 `pending`、`running`、`settling` 任务；
- 使用数据库条件更新或租约避免同一进程内重复处理；
- `pending`、`running` 写回状态和更新时间；
- `completed` 写入完整结果并进入 `settling`；
- `failed`、`canceled` 写入终态并释放余额；
- `settling` 使用固定请求 ID重试幂等结算；
- 进程启动后自动扫描未终态任务，不依赖浏览器在线；
- 临时网络错误使用退避重试，不把任务误标为失败；
- 上游明确 `404` 记录为失败并释放余额，因为映射已无法恢复。

协调器停止时等待当前短请求结束，不中断 LeoStudio 生成。Sub2API 重启不会导致 LeoStudio 任务重新提交。

## 9. 余额冻结与计费

多任务提交必须在创建上游任务前冻结预计费用。冻结金额使用请求规格和提交时的价格快照：

```text
预计费用 = 请求分辨率每秒价格 x 请求时长 x 有效视频倍率
```

流程：

1. 校验余额或订阅资格。
2. 生成 Sub2API job ID和稳定 `request_hash`。
3. 以 `video_hold:<job_id>` 幂等冻结预计费用。
4. 调用 LeoStudio 创建异步任务。
5. 创建失败或没有收到可解析的 `202 + job_id` 时，以 `video_release:<job_id>` 释放冻结金额；不自动重提。
6. 完成后从 `result.provider.resolution` 和 `result.provider.duration` 读取实际输出。
7. 以 `video_capture:<job_id>` 幂等结算实际费用并释放差额。
8. 失败或取消时全额释放。
9. 写入一条 `billing_mode=video` 的用量记录，`request_id` 使用稳定任务结算 ID。

若实际费用高于冻结金额，任务进入可重试的结算失败状态，不重复扣款。第一版按 LeoStudio 对请求规格只做归一化、不提高分辨率或时长的契约冻结；后台日志明确记录此异常，避免静默透支。

价格、分组倍率、视频独立倍率和账号倍率使用提交时快照，避免任务运行期间管理员改价改变已提交任务费用。

## 10. 本地图片存储

临时目录位于现有数据目录下的 `video-inputs/`，不允许客户端指定路径。文件名和访问 token 均使用密码学安全随机值。

LeoStudio 与 Sub2API 同宿主机，因此返回：

```text
http://127.0.0.1:<sub2_port>/internal/video-inputs/<token>
```

内部读取路由：

- 使用真实 `RemoteAddr` 校验调用方是 loopback，不信任 `X-Forwarded-For`；
- token 不存在、过期或来源非 loopback 时统一返回 `404`；
- 只允许 `GET`，设置 `Cache-Control: no-store` 和正确图片 MIME；
- 不返回服务器文件路径、原始文件名或目录列表。

生命周期：

- 活动任务引用的文件不会按普通 TTL 清理；
- 任务终态后保留 `1` 小时，便于结果页复用参数，再删除文件；
- 未关联任务的孤儿上传在 `24` 小时后删除；
- 后台每天扫描一次过期文件；启动时也执行一次；
- 删除失败记录脱敏日志并在下次扫描重试。

第一版仅支持同宿主机直接进程部署。未来改为 Docker、跨主机或多实例时，应切换到 R2/S3 或显式配置内部可达地址，不隐式猜测网络拓扑。

## 11. 客户端工作台

新增用户路由 `/video-generation` 和侧边栏入口“视频生成”。页面不做营销说明，直接进入工作台。

### 11.1 左侧生成设置

- Leo API Key 下拉框：只显示状态可用、分组平台为 `leo` 且允许视频生成的 Key；默认选择最近使用项。
- Prompt 文本域：必填，支持 `Ctrl/Cmd + Enter` 提交。
- 模型：从 Key 所属分组可用模型中取 Seedance 模型；无可用模型时禁止提交。
- 比例：`16:9`、`9:16`、`1:1`、`4:3`、`3:4`、`21:9`、`9:21`。
- 分辨率：`480p`、`720p`、`1080p`。
- 时长：`4` 到 `15` 秒步进器。
- 音频：开关。
- 图片输入：分段控制“无 / 本地上传 / 图片 URL”；本地图片显示预览、替换和移除。
- 提交按钮只在当前表单有效时启用，提交成功后保留参数并清空临时提交状态，方便连续生成。

### 11.2 右侧任务与结果

- 顶部显示最新选中任务的视频预览、规格、费用和下载链接。
- 下方显示最近 `50` 个任务，按时间倒序。
- 状态使用稳定尺寸徽标：排队、生成中、结算中、完成、失败、已取消。
- `pending` 显示取消图标按钮；完成任务显示播放和下载；失败任务显示脱敏错误。
- 页面存在活动任务时每 `2` 秒轮询列表；全部终态后停止。
- 页面刷新后从服务端列表恢复，不把 API Key、Prompt、图片 URL 或结果写入浏览器持久存储。
- 桌面使用约 `400px + minmax(0, 1fr)` 双栏；移动端按设置、预览、任务列表顺序纵向排列。

空状态明确区分：没有 Leo API Key、分组未授权、无模型、尚未提交任务。错误通过现有 Toast 与任务行内错误共同呈现。

## 12. 权限和安全

- 上传、提交、列表、详情和取消均使用现有 API Key 鉴权和 Leo 分组边界。
- 任务查询按 `api_key_id` 隔离；不以可猜测数据库主键授权。
- Leo Key、Cookie、内网管理地址和上游错误正文不进入客户端响应、任务结果或普通日志。
- Prompt 和远程图片 URL只作为任务业务数据保存，不写入生命周期日志。
- 本地上传限制数量、大小、MIME 和读取来源；拒绝 SVG、HTML 和任意文件。
- 远程 `image_url` 只允许 `http` 或 `https`，禁止 userinfo；实际下载由 LeoStudio 完成。
- 所有异步状态和媒体响应禁止缓存。

## 13. 错误处理

| 场景 | 行为 |
|---|---|
| 请求或图片无效 | `400` 或 `413`，不冻结、不创建任务 |
| 余额不足 | `402`，不调用 LeoStudio |
| 无可用 Leo 账号 | 释放冻结，返回现有无账号错误 |
| 创建阶段可确认失败 | 切换账号或最终释放冻结并返回错误 |
| 创建响应不确定 | 不自动重提；本地任务标记失败并释放冻结。上游可能留下无法追踪的任务，但不会造成重复生成或无限冻结 |
| 状态轮询临时失败 | 保持原状态并退避重试 |
| 上游任务失败 | 保存脱敏错误，释放冻结 |
| pending 取消成功 | 保存 `canceled`，释放冻结 |
| 结算暂时失败 | 保持 `settling` 并幂等重试，不向用户重复扣费 |
| 本地文件读取失败 | 上游任务最终失败，释放冻结并清理孤儿文件 |

## 14. 测试与验收

后端必须先写失败测试并覆盖：

- `video_jobs` schema、迁移、索引和仓储状态转换；
- `Prefer` 精确解析，同步路径保持不变；
- 异步创建的账号选择、模型映射、响应头和任务映射；
- 列表、详情、API Key 隔离、未知任务和 pending 取消；
- 重启恢复、后台轮询、完成、失败、取消和临时上游错误；
- 冻结、结算、差额释放、失败释放、幂等重试和用量日志；
- 本地上传的 MIME、大小、随机文件名、loopback 读取、孤儿和终态清理；
- 错误与日志不泄露 Leo Key、Prompt、图片 URL 或本地路径。

前端必须先写失败测试并覆盖：

- 路由和侧边栏入口；
- Leo Key 过滤、默认选择和无 Key 空状态；
- 文本、URL、本地上传三种请求；
- 参数校验、异步提交和任务轮询；
- 多任务状态、取消、完成预览、失败提示和刷新恢复；
- 桌面与移动端布局不重叠、不溢出。

最终验证包括后端全量测试、前端全量测试、类型检查、Lint、生产构建、迁移检查、`git diff --check`，以及使用 mock LeoStudio 的端到端验收。真实 Leonardo 生成会消耗额度，未得到单独授权前不执行。

## 15. 文档、发布与回滚

正式实现同步更新 `docs/LEO_VIDEO_CHANNEL.md`，补充异步 API、本地图片目录、同宿主机限制、余额冻结和客户端使用说明。

回滚顺序：

1. 停止新任务提交并等待现有任务进入终态；
2. 回滚客户端入口和异步 handler；
3. 停止协调器；
4. 保留 `video_jobs` 表作为审计数据，不自动删除；
5. 恢复同步 Leo 渠道版本。

数据库迁移为新增表，不改写已有业务表数据。旧版本不会读取该表；需要彻底删除时必须另行批准数据清理操作。
