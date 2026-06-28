# Token 激励计划

Token 激励计划按自然周统计用户在 `usage_logs` 中产生的 Token 消耗，并允许用户在本周内按已达标档位自行领取奖励。

## 状态接口行为

- 用户端状态接口：`GET /user/token-incentive`
- `tokens` 始终表示当前自然周实时累计 Token 消耗。
- `claimed=true` 只表示本周已经领取过至少一个档位，不会冻结 `tokens` 的展示值。
- `claimable=true` 表示当前还有已达标但未领取的档位，可以继续领取。
- `reward_amount` 表示当前可领取档位的奖励金额；`claimed_reward_amount` 表示本周已领取奖励总额。
- `claimed_threshold_tokens` 表示本周已领取过的档位阈值，用于前端区分已领取和可领取档位。

## 领取规则

- 用户端领取接口：`POST /user/token-incentive/claim`
- 每个用户每个 `week_start` 下，每个 `threshold_tokens` 档位只能领取一次。
- 达到多个档位时，按从低到高依次领取，每次领取该档位配置的完整金额，不做补差。
- 例如默认规则下，本周达到 5000 万 Token 可领 2 元；继续达到 1 亿 Token 后可再领 5 元；达到 5 亿 Token 后可再领 10 元。
- 领取时会在数据库内重新统计本周 Token 并校验是否达到当前未领取档位，避免前端或并发请求绕过资格判断。
- 奖励需在本周内领取，周窗口结束后不补发上一周未领取奖励。
