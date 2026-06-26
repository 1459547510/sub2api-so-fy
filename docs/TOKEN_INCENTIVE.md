# Token 激励计划

Token 激励计划按自然周统计用户在 `usage_logs` 中产生的 Token 消耗，并允许用户在本周内按已达到的最高档位领取一次返现。

## 状态接口口径

- 用户端状态接口：`GET /user/token-incentive`
- `tokens` 始终表示当前自然周实时累计 Token 消耗。
- `claimed=true` 只表示本周已经领取过奖励，不会冻结 `tokens` 的展示值。
- 已领取后，`reward_amount` 表示本周实际领取时发放的金额；后续本周继续产生的新 Token 仍会计入 `tokens` 和进度条展示，但不能重复领取。

## 领取口径

- 用户端领取接口：`POST /user/token-incentive/claim`
- 每个用户每个 `week_start` 只能领取一次。
- 领取时会在数据库内重新统计本周 Token 并校验是否达到档位，避免前端或并发请求绕过资格判断。
- 奖励需在本周内领取，周窗口结束后不补发上一周未领取奖励。
