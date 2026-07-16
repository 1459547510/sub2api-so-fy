# Leo 同步视频渠道设计

## 1. 目标

在 Sub2API 中新增独立一级平台 `leo`，让管理员将一个 LeoStudio HTTP 服务实例配置为一个上游账号。LeoStudio 继续负责 Cookie 号池、余额检查、账号轮换和 Leonardo 视频生成；Sub2API 负责用户 API Key 鉴权、分组授权、LeoStudio 实例调度、请求转发、故障转移、计费和用量记录。

第一期只支持同步视频生成：

```http
POST /v1/videos/generations
```

客户端等待 LeoStudio 完成生成，Sub2API 原样返回 LeoStudio 的最终 JSON 响应。

## 2. 已确认决策

- `leo` 是独立平台，不复用 `grok`、`openai` 或其他平台标识。
- 一个 Sub2API Leo 账号对应一个 LeoStudio 服务实例，而不是一条 Leonardo Cookie。
- Leo 账号只支持 `apikey` 类型。
- 账号凭据包含 `base_url`、`api_key` 和现有 `model_mapping`。
- Sub2API 使用 `Authorization: Bearer <api_key>` 调用 LeoStudio。
- LeoStudio 通过内网 HTTP 或 HTTPS 地址访问。
- 第一期仅提供同步视频生成，不提供异步任务包装。
- 第一期不提供图片生成、视频编辑、视频扩展、任务状态查询或聊天接口。
- 不修改数据库结构；平台、凭据、模型映射和视频价格使用现有字段保存。

## 3. 非目标

- 不在 Sub2API 中导入、展示或管理 LeoStudio Cookie 池。
- 不代理 LeoStudio 管理界面或 Cookie 管理接口。
- 不为同步请求创建任务表、后台队列或重启恢复机制。
- 不兼容 xAI/Grok 的异步 `request_id` 和 `/v1/videos/{id}` 协议。
- 不为 Leo 增加 OAuth、Refresh Token、SSO 或浏览器登录流程。
- 不为 LeoStudio 假定美元默认成本。

## 4. 平台与账号模型

后端平台常量、分组平台校验、账号平台校验、调度快照、用量平台和前端平台类型均增加 `leo`。

Leo 账号使用现有账号实体，约束如下：

| 字段 | 要求 |
|---|---|
| `platform` | 固定为 `leo` |
| `type` | 固定为 `apikey` |
| `credentials.base_url` | 必填，指向以 `/v1` 结尾的 LeoStudio API 根地址 |
| `credentials.api_key` | 必填，非空 Bearer Token |
| `credentials.model_mapping` | 必填，至少包含一个外部模型到 LeoStudio 模型的映射 |
| `concurrency` | 使用现有账号并发控制 |
| `proxy_id` | 沿用现有可选代理配置 |
| `status` | 沿用现有启用、禁用和异常状态 |

管理端预置以下模型映射候选，但仍允许管理员手工填写：

- `seedance-2.0`
- `seedance-2.0-fast`

预置模型只是表单辅助，不限制 LeoStudio 二次开发后可接受的其他模型。

## 5. URL 与鉴权规则

`base_url` 支持内网 HTTP 和 HTTPS，例如：

```text
http://leostudio:8000/v1
http://10.0.0.20:8000/v1
https://leo.internal.example/v1
```

保存和使用前执行以下校验：

- scheme 只能是 `http` 或 `https`；
- host 必须存在；
- 禁止 URL userinfo、query 和 fragment；
- 规范化末尾斜杠；
- API 根路径必须归一为 `/v1`，避免重复拼接；
- 仅管理员可以创建和修改 Leo 账号。

生成请求发送到：

```text
{normalized_base_url}/videos/generations
```

健康检查发送到移除 `/v1` 后的：

```text
{origin}/health
```

生成和健康检查请求均携带：

```http
Authorization: Bearer <api_key>
Accept: application/json
```

`api_key` 使用现有账号凭据保护和脱敏机制，不写入普通日志、错误信息或管理端列表响应。

## 6. 请求处理流程

1. 客户端使用 Sub2API API Key 调用 `POST /v1/videos/generations`。
2. 网关完成现有用户鉴权、分组分配、余额或订阅资格检查。
3. 路由只在分组平台为 `leo` 时进入 Leo 视频处理器。
4. 处理器仅接受 JSON，请求体必须非空且为合法 JSON。
5. `model` 和 `prompt` 必须为非空字符串。
6. 内容审核和媒体权限沿用现有媒体生成入口；底层继续使用现有 `allow_image_generation` 字段，Leo 管理界面显示为“允许视频生成”。
7. 调度器按平台 `leo`、模型、账号状态、并发和现有调度策略选择 Leo 账号。
8. 使用账号 `model_mapping` 将客户端模型名改写为 LeoStudio 模型名。
9. 除 `model` 外的 JSON 字段保持原样，随后添加 Bearer Token 并发送给 LeoStudio。
10. 同步等待 LeoStudio 响应，将允许的响应头、HTTP 状态和响应体原样返回客户端。
11. 成功后记录所选账号、模型、分辨率、时长、费用和请求耗时。

首期明确支持的 LeoStudio 请求字段为：

```json
{
  "model": "seedance-2.0",
  "prompt": "A cinematic city at night",
  "aspect_ratio": "16:9",
  "resolution": "720p",
  "duration": 8,
  "audio": false,
  "image_url": "https://example.com/start.png"
}
```

Sub2API 不重新定义 LeoStudio 参数规则。除必填字段和模型映射外，其余业务校验由 LeoStudio 完成，以便兼容其后续新增参数。

## 7. 路由边界

Leo 分组只允许：

```http
POST /v1/videos/generations
POST /videos/generations
```

第二条路径沿用网关现有的无 `/v1` 兼容入口，内部行为与第一条一致。

以下接口对 Leo 分组返回明确的 `404 not_found_error`，消息说明当前平台不支持该能力：

- 图片生成与编辑；
- 视频编辑与扩展；
- 视频状态查询；
- Chat Completions；
- Responses；
- Embeddings；
- WebSocket。

## 8. 同步、超时与客户端断开

LeoStudio 当前视频流水线最长等待约 8 分钟。Sub2API 使用现有 OpenAI 上游响应头超时配置，默认 600 秒，可覆盖该范围，不新增 Leo 专属超时配置。

请求在 LeoStudio 返回前持续占用所选账号的一个并发槽。客户端断开后，上游请求继续执行到完成或超时，避免已消耗 Leonardo 额度的任务被 Sub2API 主动中止；完成后仍记录实际用量，但无法向已断开的客户端返回结果。

## 9. 错误分类与故障转移

| 场景 | Sub2API 行为 |
|---|---|
| 请求 JSON 非法、缺少 `model` 或 `prompt` | 返回 `400 invalid_request_error`，不选择上游账号 |
| LeoStudio `400/422` | 透传请求错误，不切换账号 |
| LeoStudio `401/403` | 记录鉴权错误并将账号标记异常，将该账号移出本次候选后切换其他 Leo 账号 |
| LeoStudio `429` | 视为实例暂时不可用，按现有规则切换账号 |
| LeoStudio `5xx` | 视为上游故障，按现有规则切换账号 |
| DNS、连接、TLS 或响应读取失败 | 按现有上游故障规则切换账号 |
| 全部账号失败 | 返回统一 `502 upstream_error` |

只有在尚未向客户端提交响应时才能故障转移。错误正文按现有上游错误脱敏规则处理，不回传凭据、内网地址或完整上游响应。

## 10. 响应与计费

成功响应保持 LeoStudio 当前结构，不转换为 Grok 异步响应：

```json
{
  "created": 1784160000,
  "data": [
    {
      "url": "https://example.com/video.mp4",
      "mp4_url": "https://example.com/video.mp4"
    }
  ],
  "provider": {
    "generation_id": "generation-id",
    "used_cookie_id": 12,
    "model": "seedance-2.0",
    "resolution": "RESOLUTION_720",
    "duration": 8,
    "aspect_ratio": "16:9",
    "audio": false
  }
}
```

计费公式为：

```text
费用 = 对应分辨率每秒价格 × 实际时长 × 当前有效分组倍率
```

规则如下：

- Leo 分组创建和更新时必须配置 `video_price_480p`、`video_price_720p`、`video_price_1080p`；
- 三档价格单位均为 USD/秒，允许 `0` 表示免费；
- 优先从成功响应的 `provider.resolution` 和 `provider.duration` 读取实际值；
- 响应缺少实际值时，使用经过现有规范化规则处理的请求分辨率和时长；
- 使用现有视频独立倍率开关和 `video_rate_multiplier`；未启用时使用分组常规倍率；
- 用量记录平台为 `leo`，同时保存外部模型、映射后的上游模型、账号、分辨率、时长和视频数量；
- 不使用 Grok 默认视频价格，也不回退到图片默认价格。

## 11. 管理端设计

分组管理：

- 平台选择、筛选和徽章增加 `Leo`；
- 创建 Leo 分组时显示“允许视频生成”和三档视频价格；
- 缺少任一视频价格时阻止保存并显示字段错误；
- 不显示 Leo 不支持的聊天或 OAuth 专属配置。

账号管理：

- 平台选择、筛选、列表徽章增加 `Leo`；
- Leo 仅提供 API Key 账号类型；
- 表单显示账号名称、Base URL、API Key、模型映射、并发、代理、状态和分组；
- API Key 输入使用密码框，编辑时保持既有秘密值保留语义；
- 连接测试调用 `/health` 并显示成功、鉴权失败、连接失败或响应异常；
- 平台图标使用现有图标体系中的视频图标，不引入自绘 SVG。

所有平台枚举、筛选项、平台配额视图、渠道展示、错误透传规则和运维筛选中需要明确识别 `leo`，避免创建后在管理界面中变成未知平台。

## 12. 组件边界

后端按现有结构增加以下职责，不复制 Grok/xAI 专属逻辑：

- Leo 平台常量和账号辅助方法：识别 Leo API Key 账号并读取规范化凭据；
- Leo URL 构造与校验：只负责 Base URL、生成地址和健康检查地址；
- Leo 视频服务：构造 Bearer 请求、转发同步响应、解析计费元数据；
- Leo 视频处理器：执行鉴权、权限、调度、故障转移和用量记录；
- 网关路由分发：根据分组平台将视频生成请求交给 Leo 处理器；
- 账号连接测试：调用 LeoStudio `/health`。

前端只扩展现有平台选项和账号/分组表单，不创建新的独立管理页面。

## 13. 验证策略

后端单元与集成测试覆盖：

- `leo` 通过分组和账号平台校验；
- 非 API Key Leo 账号被拒绝；
- Base URL 正规化及非法 URL 拒绝；
- 生成与健康检查请求携带正确 Bearer Token；
- 模型映射只改写 `model`，其他 JSON 字段保持不变；
- 同步成功响应原样透传；
- 响应中的实际分辨率和时长用于计费；
- 缺少实际元数据时使用请求值；
- `400/422` 不故障转移；
- `401/403/429/5xx` 和传输错误按设计故障转移；
- 所有账号失败时返回统一错误；
- Leo 分组只有视频生成路由可用；
- 视频价格必填且费用公式正确；
- API Key 不出现在日志和管理响应中。

前端测试覆盖：

- 分组和账号平台选项出现 Leo；
- Leo 账号只显示 API Key 表单；
- Base URL、API Key、模型映射和并发字段正确提交；
- API Key 回显脱敏；
- Leo 分组缺少视频价格时无法保存；
- 平台徽章、筛选和视频权限文案正确。

端到端验收使用可控的模拟 LeoStudio HTTP 服务验证：

1. `/health` 校验 Bearer Token；
2. `/v1/videos/generations` 校验请求体和模型映射；
3. 模拟同步成功、鉴权失败、限流、服务异常和超时；
4. 核对客户端响应、账号切换、用量记录和余额扣减。

## 14. 文档与部署说明

实现时同步更新 `docs/`，至少说明：

- LeoStudio 需要实现 Bearer Token 校验；
- LeoStudio 服务地址必须能被 Sub2API 后端访问；
- 示例账号、分组、模型映射和三档视频价格配置；
- 同步请求可能持续数分钟；
- 客户端断开后任务仍可能消耗 Leonardo 额度；
- 仅支持视频生成的接口边界；
- 使用 `curl` 完成健康检查和视频生成验证的方法。

## 15. 成功标准

- 管理员可以创建 Leo 分组和 Leo API Key 账号，无需手工修改数据库；
- 账号连接测试可以验证 LeoStudio 地址与 Bearer Token；
- 客户端通过 Sub2API API Key 同步生成 Leo 视频并获得 LeoStudio 最终响应；
- 模型映射、调度、并发、代理和故障转移均沿用现有平台能力；
- 计费按实际分辨率和时长执行，不使用错误的默认价格；
- 非支持接口被明确拒绝；
- 自动化测试、文档和 `progress.md` 均与行为一致；
- 本次实现不引入数据库迁移和异步任务系统。
