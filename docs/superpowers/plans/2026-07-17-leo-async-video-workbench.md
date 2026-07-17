# Leo Async Video Workbench Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a durable, billable Leo asynchronous video API and a user-facing Vue workbench supporting text, remote-image, and local-image multi-task generation.

**Architecture:** Sub2API owns an opaque `video_jobs` record that fixes every upstream job to the Leo account selected at submission. A background runtime polls LeoStudio's existing async job API and performs idempotent terminal billing; local input files are served only over a loopback-only token route. The browser uses the user's selected Sub2API API Key and never receives LeoStudio credentials or upstream job identifiers.

**Tech Stack:** Go 1.24, Gin, Ent/PostgreSQL, Wire, Vue 3, TypeScript, Axios/fetch, Vitest, Tailwind CSS.

---

## File Structure

**Backend data and repository**

- Create `backend/ent/schema/video_job.go`: Ent definition for durable video jobs.
- Create `backend/migrations/182_video_jobs.sql`: additive PostgreSQL table and indexes.
- Create `backend/internal/service/video_job.go`: domain types, state constants, repository contract, DTO conversion.
- Create `backend/internal/repository/video_job_repo.go`: Ent-backed create/read/list/claim/transition methods.
- Modify `backend/internal/repository/wire.go`: provide the repository.

**Backend Leo integration and billing**

- Create `backend/internal/service/leo_video_async.go`: `Prefer` parsing and upstream create/get/delete calls.
- Create `backend/internal/service/video_job_billing.go`: request pricing snapshot, balance hold, terminal settlement, release, and usage log.
- Modify `backend/internal/service/usage_billing.go`: add optional async video balance-hold commands.
- Modify `backend/internal/repository/usage_billing_repo.go`: implement video hold/capture/release using existing dedup tables.
- Modify `backend/internal/service/openai_gateway_usage.go`: allow a precomputed cost override for durable settlement without changing normal requests.

**Backend orchestration and HTTP**

- Create `backend/internal/service/video_job_service.go`: submit/list/get/cancel orchestration.
- Create `backend/internal/service/video_job_runtime.go`: poll active jobs, settle terminal jobs, resume after restart.
- Create `backend/internal/service/video_input_store.go`: local file validation, token lookup, and cleanup.
- Create `backend/internal/handler/leo_video_async.go`: async API handlers and upload handler.
- Modify `backend/internal/handler/leo_video.go`: dispatch exact `Prefer: respond-async` requests to the async service while preserving sync behavior.
- Modify `backend/internal/server/routes/gateway.go`: register list/detail/cancel/upload and loopback input routes.
- Modify `backend/internal/service/wire.go`, `backend/cmd/server/wire.go`, `backend/cmd/server/wire_gen.go`: construct/start/stop the runtime.

**Frontend**

- Create `frontend/src/api/videoGeneration.ts`: typed upload, submit, list, detail, and cancel calls.
- Create `frontend/src/views/user/VideoGenerationView.vue`: responsive two-column workbench.
- Create `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: UI flow coverage.
- Create `frontend/src/api/__tests__/videoGeneration.spec.ts`: API contract coverage.
- Modify `frontend/src/router/index.ts`: `/video-generation` route.
- Modify `frontend/src/components/layout/AppSidebar.vue`: video navigation item and icon.
- Modify `frontend/src/i18n/locales/zh/*.ts`, `frontend/src/i18n/locales/en/*.ts`: navigation and workbench copy.

**Docs and logs**

- Modify `docs/LEO_VIDEO_CHANNEL.md`: async API, local files, billing, and workbench instructions.
- Modify `progress.md`: append one tested task record per implementation commit.

### Task 1: Durable Video Job Model

**Files:**
- Create: `backend/ent/schema/video_job.go`
- Create: `backend/migrations/182_video_jobs.sql`
- Create: `backend/internal/service/video_job.go`
- Create: `backend/internal/repository/video_job_repo.go`
- Test: `backend/internal/repository/video_job_repo_test.go`
- Modify: `backend/internal/repository/wire.go`

- [ ] **Step 1: Write the failing repository test**

Define a job with `JobID: "vidjob_test"`, `AccountID: 9`, `UpstreamJobID: 42`, `Status: pending`, then assert create, API-key-scoped lookup, descending list, and `pending -> running` conditional transition. Also assert another API key receives `ErrVideoJobNotFound`.

```go
func TestVideoJobRepositoryCreateListAndTransition(t *testing.T) {
	ctx := context.Background()
	repo := newVideoJobRepositoryTestRepo(t)
	job := &service.VideoJob{JobID: "vidjob_test", UserID: 1, APIKeyID: 2, GroupID: 3, AccountID: 9, UpstreamJobID: 42, Status: service.VideoJobPending, RequestedModel: "seedance-2.0", UpstreamModel: "seedance-2.0", Prompt: "waves", Resolution: "720p", DurationSeconds: 8, AspectRatio: "16:9"}
	require.NoError(t, repo.CreateVideoJob(ctx, job))
	got, err := repo.GetVideoJobForAPIKey(ctx, "vidjob_test", 2)
	require.NoError(t, err)
	require.Equal(t, int64(42), got.UpstreamJobID)
	require.NoError(t, repo.TransitionVideoJob(ctx, "vidjob_test", []string{service.VideoJobPending}, service.VideoJobRunning, service.VideoJobTransition{}))
	_, err = repo.GetVideoJobForAPIKey(ctx, "vidjob_test", 999)
	require.ErrorIs(t, err, service.ErrVideoJobNotFound)
}
```

- [ ] **Step 2: Run the test and verify RED**

Run: `cd backend && go test ./internal/repository -run TestVideoJobRepositoryCreateListAndTransition -count=1`

Expected: compile failure because `VideoJobRepository` and `VideoJob` do not exist.

- [ ] **Step 3: Add the schema, migration, domain contract, and repository**

Use a string public ID and JSON columns for `result` and `billing_snapshot`. State changes must use `UPDATE ... WHERE status = ANY($allowed)` so two pollers cannot settle the same row.

```go
const (
	VideoJobPending   = "pending"
	VideoJobRunning   = "running"
	VideoJobSettling  = "settling"
	VideoJobCompleted = "completed"
	VideoJobFailed    = "failed"
	VideoJobCanceled  = "canceled"
)

var ErrVideoJobNotFound = errors.New("video job not found")
var ErrVideoJobTransitionConflict = errors.New("video job transition conflict")
```

Run `go generate ./ent` after adding the Ent schema, then add `NewVideoJobRepository` to `repository.ProviderSet`.

- [ ] **Step 4: Verify GREEN and schema consistency**

Run:

```powershell
cd backend
go test ./internal/repository -run 'VideoJob|Migration' -count=1
go test ./ent/schema -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/ent backend/migrations/182_video_jobs.sql backend/internal/service/video_job.go backend/internal/repository/video_job_repo.go backend/internal/repository/video_job_repo_test.go backend/internal/repository/wire.go
git commit -m "feat: persist leo video jobs"
```

### Task 2: LeoStudio Async Client

**Files:**
- Create: `backend/internal/service/leo_video_async.go`
- Test: `backend/internal/service/leo_video_async_test.go`

- [ ] **Step 1: Write failing protocol tests**

Cover exact `Prefer` parsing, mapped model submit, `202` response decoding, status result decoding, pending cancellation, Bearer auth, and no secret in public errors.

```go
func TestPrefersLeoRespondAsync(t *testing.T) {
	h := http.Header{}
	h.Add("Prefer", "wait=5, Respond-Async")
	require.True(t, PrefersLeoRespondAsync(h))
	h.Set("Prefer", "respond-async=false")
	require.False(t, PrefersLeoRespondAsync(h))
}
```

- [ ] **Step 2: Verify RED**

Run: `cd backend && go test ./internal/service -run 'PrefersLeoRespondAsync|LeoAsync' -count=1`

Expected: compile failure for missing async helpers.

- [ ] **Step 3: Implement the narrow client**

Add:

```go
type LeoAsyncAccepted struct { JobID int64 `json:"job_id"`; Status string `json:"status"`; StatusURL string `json:"status_url"` }
type LeoAsyncJob struct { JobID int64 `json:"job_id"`; Status string `json:"status"`; Result json.RawMessage `json:"result,omitempty"`; Error *LeoAsyncJobError `json:"error,omitempty"` }

func (s *OpenAIGatewayService) CreateLeoAsyncVideo(ctx context.Context, account *Account, body []byte) (*LeoAsyncAccepted, error)
func (s *OpenAIGatewayService) GetLeoAsyncVideo(ctx context.Context, account *Account, upstreamJobID int64) (*LeoAsyncJob, error)
func (s *OpenAIGatewayService) CancelLeoAsyncVideo(ctx context.Context, account *Account, upstreamJobID int64) (*LeoAsyncJob, error)
```

Only create sends `Prefer: respond-async`. Status and cancel return sanitized typed errors and never write directly to Gin.

- [ ] **Step 4: Verify GREEN**

Run: `cd backend && go test ./internal/service -run 'PrefersLeoRespondAsync|LeoAsync' -count=1`

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/internal/service/leo_video_async.go backend/internal/service/leo_video_async_test.go
git commit -m "feat: call leo async video jobs"
```

### Task 3: Idempotent Balance Hold and Settlement

**Files:**
- Create: `backend/internal/service/video_job_billing.go`
- Test: `backend/internal/service/video_job_billing_test.go`
- Modify: `backend/internal/service/usage_billing.go`
- Modify: `backend/internal/repository/usage_billing_repo.go`
- Modify: `backend/internal/service/openai_gateway_usage.go`
- Test: `backend/internal/repository/usage_billing_repo_unit_test.go`
- Test: `backend/internal/service/openai_gateway_record_usage_test.go`

- [ ] **Step 1: Write failing hold and settlement tests**

Assert balance-backed jobs reserve once, completed jobs charge once and release the full hold, failed/canceled jobs only release, and a second settlement is a no-op. Assert `CostOverride` keeps submission-time prices after group price changes.

```go
func TestVideoJobSettlementIsIdempotent(t *testing.T) {
	job := testVideoJob(service.VideoJobSettling)
	billing := newFakeVideoBillingRepo()
	svc := service.VideoJobBillingService{BillingRepo: billing, UsageRecorder: billing}
	require.NoError(t, svc.SettleCompleted(context.Background(), job, testLeoResult("RESOLUTION_720", 12)))
	require.NoError(t, svc.SettleCompleted(context.Background(), job, testLeoResult("RESOLUTION_720", 12)))
	require.Len(t, billing.usageCommands, 1)
	require.Len(t, billing.releaseCommands, 1)
}
```

- [ ] **Step 2: Verify RED**

Run: `cd backend && go test ./internal/service ./internal/repository -run 'VideoJobBilling|VideoBalanceHold|CostOverride' -count=1`

Expected: compile failure for missing billing types.

- [ ] **Step 3: Implement video hold and optional cost override**

Add `VideoBalanceHoldCommand` and repository methods parallel to the existing batch-image hold path, but with `JobID` and request prefixes `video_hold:`, `video_release:`. Keep normal request billing unchanged when `OpenAIRecordUsageInput.CostOverride == nil`.

Settlement order for balance billing is intentional: idempotent normal usage billing charges actual cost and updates all quota dimensions, then idempotent hold release returns the frozen estimate. Net balance change equals actual cost even after a crash between the two operations.

- [ ] **Step 4: Verify GREEN and existing billing regressions**

Run:

```powershell
cd backend
go test ./internal/service ./internal/repository -run 'VideoJobBilling|VideoBalanceHold|CostOverride|BatchImageBalanceHold' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/internal/service/video_job_billing.go backend/internal/service/video_job_billing_test.go backend/internal/service/usage_billing.go backend/internal/repository/usage_billing_repo.go backend/internal/repository/usage_billing_repo_unit_test.go backend/internal/service/openai_gateway_usage.go backend/internal/service/openai_gateway_record_usage_test.go
git commit -m "feat: settle async leo video billing"
```

### Task 4: Video Job Submission and Access Service

**Files:**
- Create: `backend/internal/service/video_job_service.go`
- Test: `backend/internal/service/video_job_service_test.go`

- [ ] **Step 1: Write failing service tests**

Cover validation, scheduler account selection, hold-before-upstream ordering, model mapping, durable mapping, account failover before `202`, at-most-once ambiguous failure, API-key-scoped list/detail, and pending-only cancel.

- [ ] **Step 2: Verify RED**

Run: `cd backend && go test ./internal/service -run TestVideoJobService -count=1`

Expected: compile failure for `VideoJobService`.

- [ ] **Step 3: Implement minimal orchestration**

Use:

```go
type CreateVideoJobInput struct {
	APIKey *APIKey
	User *User
	Subscription *UserSubscription
	Body []byte
	LocalInputName string
}

func (s *VideoJobService) Create(ctx context.Context, in CreateVideoJobInput) (*VideoJob, error)
func (s *VideoJobService) List(ctx context.Context, apiKeyID int64, limit int, status string) ([]*VideoJob, error)
func (s *VideoJobService) Get(ctx context.Context, jobID string, apiKeyID int64) (*VideoJob, error)
func (s *VideoJobService) Cancel(ctx context.Context, jobID string, apiKeyID int64) (*VideoJob, error)
```

Generate IDs as `vidjob_` plus 24 random bytes encoded with base64url. Do not return `AccountID` or `UpstreamJobID` in public DTOs.

- [ ] **Step 4: Verify GREEN**

Run: `cd backend && go test ./internal/service -run TestVideoJobService -count=1`

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/internal/service/video_job_service.go backend/internal/service/video_job_service_test.go
git commit -m "feat: manage leo async video jobs"
```

### Task 5: Background Status Runtime

**Files:**
- Create: `backend/internal/service/video_job_runtime.go`
- Test: `backend/internal/service/video_job_runtime_test.go`
- Modify: `backend/internal/service/wire.go`
- Modify: `backend/cmd/server/wire.go`
- Regenerate: `backend/cmd/server/wire_gen.go`

- [ ] **Step 1: Write failing runtime tests**

Use a fake clock and fake Leo client. Cover pending/running polling, completed-to-settling-to-completed, failed release, transient retry, restart scan, and `Stop()` waiting for an in-flight poll.

- [ ] **Step 2: Verify RED**

Run: `cd backend && go test ./internal/service -run TestVideoJobRuntime -count=1`

Expected: compile failure for runtime types.

- [ ] **Step 3: Implement runtime**

```go
type VideoJobRuntime struct {
	repo VideoJobRepository
	accounts AccountRepository
	gateway *OpenAIGatewayService
	billing *VideoJobBillingService
	pollInterval time.Duration
	cancel context.CancelFunc
	done chan struct{}
}
```

Each tick claims at most 50 active rows. Use account affinity from the row, never scheduler selection. Start it from `ProvideVideoJobRuntime`; add `Stop()` to application cleanup before Ent closes.

- [ ] **Step 4: Verify GREEN and Wire**

Run:

```powershell
cd backend
go test ./internal/service -run TestVideoJobRuntime -count=1
go generate ./cmd/server
go test ./cmd/server -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/internal/service/video_job_runtime.go backend/internal/service/video_job_runtime_test.go backend/internal/service/wire.go backend/cmd/server/wire.go backend/cmd/server/wire_gen.go
git commit -m "feat: reconcile leo video jobs"
```

### Task 6: Local Video Input Files

**Files:**
- Create: `backend/internal/service/video_input_store.go`
- Test: `backend/internal/service/video_input_store_test.go`
- Create: `backend/internal/handler/video_input.go`
- Test: `backend/internal/handler/video_input_test.go`

- [ ] **Step 1: Write failing storage and HTTP tests**

Cover 10 MiB limit, PNG/JPEG/WebP detection, SVG rejection, random names, loopback GET, non-loopback `404`, orphan 24-hour cleanup, and terminal one-hour cleanup.

- [ ] **Step 2: Verify RED**

Run: `cd backend && go test ./internal/service ./internal/handler -run 'VideoInput' -count=1`

Expected: compile failure for video input types.

- [ ] **Step 3: Implement store and handlers**

Store beneath the configured data directory in `video-inputs/`. The returned URL is constructed from the configured server port:

```go
func (s *VideoInputStore) InternalURL(token string) string {
	return "http://127.0.0.1:" + strconv.Itoa(s.port) + "/internal/video-inputs/" + url.PathEscape(token)
}
```

The public upload handler uses API Key auth. The internal GET handler checks `net.SplitHostPort(req.RemoteAddr)` and `net.ParseIP(host).IsLoopback()` without consulting forwarded headers.

- [ ] **Step 4: Verify GREEN**

Run: `cd backend && go test ./internal/service ./internal/handler -run 'VideoInput' -count=1`

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/internal/service/video_input_store.go backend/internal/service/video_input_store_test.go backend/internal/handler/video_input.go backend/internal/handler/video_input_test.go
git commit -m "feat: host temporary leo video inputs"
```

### Task 7: Gateway Async Routes

**Files:**
- Create: `backend/internal/handler/leo_video_async.go`
- Test: `backend/internal/handler/leo_video_async_test.go`
- Modify: `backend/internal/handler/leo_video.go`
- Modify: `backend/internal/handler/openai_gateway_handler.go`
- Modify: `backend/internal/server/routes/gateway.go`
- Test: `backend/internal/server/routes/gateway_test.go`
- Test: `backend/internal/handler/leo_video_integration_test.go`

- [ ] **Step 1: Write failing HTTP contract tests**

Assert sync remains `200`, async returns `202`, public `Location` hides upstream ID, list/detail are API-key isolated, upload is multipart, internal input GET is loopback-only, and cancel conflicts after running.

- [ ] **Step 2: Verify RED**

Run: `cd backend && go test ./internal/handler ./internal/server/routes -run 'LeoAsync|VideoJob|VideoInput' -count=1`

Expected: route tests return `404` and handler symbols are missing.

- [ ] **Step 3: Register routes and dispatch**

Register before the generic `/videos/:request_id` route:

```go
gateway.POST("/videos/uploads", h.OpenAIGateway.LeoVideoUpload)
gateway.GET("/videos/jobs", h.OpenAIGateway.LeoVideoJobs)
gateway.GET("/videos/jobs/:job_id", h.OpenAIGateway.LeoVideoJob)
gateway.DELETE("/videos/jobs/:job_id", h.OpenAIGateway.CancelLeoVideoJob)
```

Register `/internal/video-inputs/:token` outside API Key middleware and enforce loopback in the handler. `LeoVideoGeneration` checks `PrefersLeoRespondAsync` before acquiring long-lived sync concurrency slots.

- [ ] **Step 4: Verify GREEN and sync regression**

Run:

```powershell
cd backend
go test ./internal/handler ./internal/server/routes ./internal/service -run 'Leo|VideoJob|VideoInput|VideoGeneration' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add backend/internal/handler backend/internal/server/routes/gateway.go backend/internal/server/routes/gateway_test.go
git commit -m "feat: expose leo async video api"
```

### Task 8: Frontend Video API

**Files:**
- Create: `frontend/src/api/videoGeneration.ts`
- Test: `frontend/src/api/__tests__/videoGeneration.spec.ts`

- [ ] **Step 1: Write failing API tests**

Mock `fetch` and assert upload multipart does not set a manual boundary, submit adds `Prefer: respond-async`, list includes the selected API Key, and cancel uses `DELETE`.

- [ ] **Step 2: Verify RED**

Run: `cd frontend && pnpm exec vitest run src/api/__tests__/videoGeneration.spec.ts`

Expected: module import failure.

- [ ] **Step 3: Implement typed API wrapper**

Export `uploadVideoInput`, `createVideoJob`, `listVideoJobs`, `getVideoJob`, and `cancelVideoJob`. Reuse `buildGatewayUrl`, a single error parser, and `Authorization: Bearer <selected key>`.

- [ ] **Step 4: Verify GREEN**

Run: `cd frontend && pnpm exec vitest run src/api/__tests__/videoGeneration.spec.ts`

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add frontend/src/api/videoGeneration.ts frontend/src/api/__tests__/videoGeneration.spec.ts
git commit -m "feat: add leo video client api"
```

### Task 9: Responsive Video Workbench

**Files:**
- Create: `frontend/src/views/user/VideoGenerationView.vue`
- Create: `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/components/layout/AppSidebar.vue`
- Modify: `frontend/src/i18n/locales/zh/common.ts`
- Modify: `frontend/src/i18n/locales/en/common.ts`
- Modify: `frontend/src/i18n/locales/zh/dashboard.ts`
- Modify: `frontend/src/i18n/locales/en/dashboard.ts`

- [ ] **Step 1: Write failing component tests**

Mount with two API keys and assert only active Leo keys appear. Cover text submit, URL submit, local upload then submit, two-second polling while active, cancel pending, playable completed result, failed row, and no-key empty state.

- [ ] **Step 2: Verify RED**

Run: `cd frontend && pnpm exec vitest run src/views/user/__tests__/VideoGenerationView.spec.ts`

Expected: component import failure.

- [ ] **Step 3: Implement the two-column workbench**

Use existing `card`, `input`, button, `Icon`, Toast, dark-mode, and responsive conventions. Stable desktop grid:

```html
<div class="grid min-w-0 gap-4 xl:grid-cols-[400px_minmax(0,1fr)]">
  <section class="min-w-0">...</section>
  <section class="min-w-0">...</section>
</div>
```

Use a segmented mode control for no image/local file/URL, native video controls for results, icon buttons with titles for cancel/download, and no nested cards.

- [ ] **Step 4: Verify GREEN, types, and navigation**

Run:

```powershell
cd frontend
pnpm exec vitest run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts
pnpm run typecheck
pnpm run lint:check
```

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts frontend/src/router/index.ts frontend/src/components/layout/AppSidebar.vue frontend/src/i18n/locales
git commit -m "feat: add leo video workbench"
```

### Task 10: Documentation and Final Acceptance

**Files:**
- Modify: `docs/LEO_VIDEO_CHANNEL.md`
- Modify: `progress.md`
- Test: all changed backend and frontend packages.

- [ ] **Step 1: Update operator and user documentation**

Document upstream commit, async API examples, same-host requirement, `video-inputs/` lifecycle, 10 MiB formats, pending-only cancellation, balance hold/settlement, client route, and rollback.

- [ ] **Step 2: Run backend verification**

```powershell
cd backend
go test ./... -count=1
go vet ./...
```

Expected: PASS with no failures.

- [ ] **Step 3: Run frontend verification**

```powershell
cd frontend
pnpm run test:run
pnpm run typecheck
pnpm run lint:check
pnpm run build
```

Expected: PASS. Existing Browserslist, dynamic-import, and chunk-size warnings are non-failing.

- [ ] **Step 4: Run repository and migration checks**

```powershell
git diff --check
git status --short
```

Expected: no whitespace errors; only task files plus pre-existing user changes remain.

- [ ] **Step 5: Append the final progress record and commit docs**

```powershell
git add -f docs/LEO_VIDEO_CHANNEL.md
git commit -m "docs: explain leo async video workbench"
```

Do not stage the pre-existing update-service, update-policy, locale, `.superpowers/`, or unrelated `progress.md` changes. Report them as retained user work.
