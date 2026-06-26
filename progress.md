## 2026-06-26 - Task: 修复 Token 激励本周消耗实时展示
### What was done
- 修复 Token 激励计划状态接口在用户已领取后继续显示领取时 Token 快照的问题，改为始终返回本周实时累计消耗。
- 保留已领取状态、领取时间和实际领取金额，避免影响每周只能领取一次的业务规则。
- 补充模块说明文档，明确状态接口和领取接口的统计口径。
### Testing
- `go test -tags unit ./internal/service -run TokenIncentive`（在 `backend` 目录）通过。
- `go test -tags unit ./internal/repository -run TokenIncentive`（在 `backend` 目录）通过。
### Notes
- `backend/internal/service/token_incentive_service.go`：状态构建改为使用实时周累计 Token，而不是已领取记录中的快照 Token。
- `backend/internal/service/token_incentive_service_test.go`：更新已领取场景测试，覆盖“领取金额保留、Token 进度实时更新”。
- `docs/TOKEN_INCENTIVE.md`：新增 Token 激励计划状态展示和领取口径说明。
- `progress.md`：追加本轮修复记录。
- 回滚方式：执行 `git checkout -- backend/internal/service/token_incentive_service.go backend/internal/service/token_incentive_service_test.go docs/TOKEN_INCENTIVE.md progress.md`，或回退包含本轮改动的提交。
