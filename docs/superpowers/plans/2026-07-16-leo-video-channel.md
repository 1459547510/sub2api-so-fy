# Leo Video Channel Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an independent `leo` platform that lets Sub2API configure authenticated LeoStudio instances and synchronously proxy `POST /v1/videos/generations` with scheduling, failover, billing, and admin UI support.

**Architecture:** Keep Leo protocol code in dedicated backend files while reusing the existing `OpenAIGatewayService` infrastructure for account scheduling, concurrency slots, HTTP transport, usage recording, and billing. A Leo account is an existing API-key account with validated `base_url`, `api_key`, and `model_mapping`; no database columns or async job system are added.

**Tech Stack:** Go 1.26, Gin, Ent, PostgreSQL-compatible JSON credentials, Vue 3, TypeScript, Vitest, Tailwind CSS.

---

## File Map

Backend platform and validation:

- Modify `backend/internal/domain/constants.go`: define `PlatformLeo`.
- Modify `backend/internal/service/domain_constants.go`: export Leo and allow platform quota accounting.
- Modify `backend/internal/model/error_passthrough_rule.go`: allow Leo error-passthrough rules.
- Modify `backend/ent/schema/user_platform_quota.go`: accept `leo` in Ent validation without changing table shape.
- Regenerate `backend/ent/`: update generated validators only.
- Modify `backend/internal/handler/admin/group_handler.go`: accept Leo group payloads.
- Modify `backend/internal/service/scheduler_snapshot_service.go`: hydrate Leo scheduling buckets.

Backend Leo account and protocol:

- Create `backend/internal/service/leo_account.go`: Leo identity, credential, URL, model, and validation helpers.
- Create `backend/internal/service/leo_account_test.go`: helper and credential validation tests.
- Modify `backend/internal/service/admin_account.go`: validate Leo credentials on create and merged updates.
- Modify `backend/internal/service/admin_group.go`: expose Leo model candidates and enforce video prices.
- Modify `backend/internal/service/account_test_service.go`: Bearer-authenticated `/health` probe.
- Create `backend/internal/service/account_test_service_leo_test.go`: health-probe tests.
- Create `backend/internal/service/leo_video.go`: synchronous request forwarding and response metadata parsing.
- Create `backend/internal/service/leo_video_test.go`: protocol, auth, response, and error tests.

Backend routing, scheduling, and billing:

- Create `backend/internal/handler/leo_video.go`: request validation, permissions, scheduling, failover, and usage submission.
- Create `backend/internal/handler/leo_video_test.go`: handler and failover tests.
- Modify `backend/internal/server/routes/gateway.go`: dispatch only video generation to Leo.
- Modify `backend/internal/service/account.go`: classify Leo accounts for shared scheduling only.
- Modify `backend/internal/service/openai_gateway_scheduling.go`: preserve `PlatformLeo` instead of normalizing it to OpenAI.
- Modify `backend/internal/service/openai_gateway_usage.go`: treat any explicit video result as video billing.
- Modify `backend/internal/service/video_billing_resolution.go`: recognize LeoStudio `RESOLUTION_*` values.
- Extend existing backend tests for platform isolation and video billing.

Frontend:

- Modify `frontend/src/types/index.ts`: add `leo` to account and group platform unions.
- Modify `frontend/src/api/admin/settings.ts`: include Leo in platform quota normalization.
- Modify `frontend/src/components/user/UserPlatformQuotaCell.vue` and `frontend/src/components/user/dashboard/UserDashboardStats.vue`: order and display Leo usage.
- Modify `frontend/src/components/common/PlatformIcon.vue`, `PlatformTypeBadge.vue`, and `GroupBadge.vue`: render Leo consistently.
- Modify `frontend/src/components/admin/account/AccountTableFilters.vue`: add Leo filtering.
- Modify `frontend/src/views/admin/GroupsView.vue` and `groupsImagePricing.ts`: add Leo groups, required video prices, and video-specific permission copy.
- Modify `frontend/src/components/account/CreateAccountModal.vue` and `EditAccountModal.vue`: Leo API-key forms.
- Modify `frontend/src/composables/useModelWhitelist.ts`: add Seedance model presets.
- Modify `frontend/src/i18n/locales/{zh,en}/admin/accounts.ts` and `overview.ts`: add Leo copy.
- Add or extend focused frontend tests for Leo account and group forms.

Documentation:

- Create `docs/LEO_VIDEO_CHANNEL.md`: deployment and request examples.
- Append `progress.md` after every completed task as required by repository policy.

---

### Task 1: Register the Leo platform end to end

**Files:**
- Modify: `backend/internal/domain/constants.go`
- Modify: `backend/internal/service/domain_constants.go`
- Modify: `backend/internal/model/error_passthrough_rule.go`
- Modify: `backend/ent/schema/user_platform_quota.go`
- Modify: generated files under `backend/ent/`
- Modify: `backend/internal/handler/admin/group_handler.go`
- Modify: `backend/internal/service/scheduler_snapshot_service.go`
- Create: `backend/internal/service/leo_platform_test.go`
- Modify: `progress.md`

- [ ] **Step 1: Write the failing platform registration test**

Create `backend/internal/service/leo_platform_test.go`:

```go
package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeoPlatformRegistration(t *testing.T) {
	require.Equal(t, "leo", PlatformLeo)
	require.True(t, IsAllowedQuotaPlatform(PlatformLeo))
	require.Contains(t, AllowedQuotaPlatforms, PlatformLeo)
}
```

- [ ] **Step 2: Run the test and verify the missing constant fails compilation**

Run from `backend/`:

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run TestLeoPlatformRegistration -count=1
```

Expected: FAIL because `PlatformLeo` is undefined.

- [ ] **Step 3: Add Leo to the authoritative backend platform lists**

Add the domain constant and service alias:

```go
// backend/internal/domain/constants.go
PlatformLeo = "leo"

// backend/internal/service/domain_constants.go
PlatformLeo = domain.PlatformLeo
```

Append `PlatformLeo` to `AllowedQuotaPlatforms`, both scheduler snapshot platform lists, and the error-passthrough allowed platform list. Extend both `CreateGroupRequest.Platform` and `UpdateGroupRequest.Platform` bindings to:

```go
binding:"omitempty,oneof=anthropic openai gemini antigravity grok leo"
```

Extend the Ent validation switch:

```go
case "anthropic", "openai", "gemini", "antigravity", "grok", "leo":
	return nil
```

- [ ] **Step 4: Regenerate Ent validators and verify no migration is produced**

Run from `backend/`:

```powershell
$env:GOTOOLCHAIN='auto'; go generate ./ent
git status --short ent migrations
```

Expected: generated Ent Go files may change; `backend/migrations/` remains unchanged.

- [ ] **Step 5: Run the platform and quota tests**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service ./internal/handler/admin -run 'TestLeoPlatformRegistration|PlatformQuota|Group' -count=1
```

Expected: PASS.

- [ ] **Step 6: Append the Task 1 verification record to `progress.md`**

Record the platform registration, Ent generation command, tests, changed files, and rollback point. Do not rewrite earlier entries.

- [ ] **Step 7: Commit Task 1**

```powershell
git add backend/internal/domain/constants.go backend/internal/service/domain_constants.go backend/internal/model/error_passthrough_rule.go backend/ent backend/internal/handler/admin/group_handler.go backend/internal/service/scheduler_snapshot_service.go backend/internal/service/leo_platform_test.go
git commit -m "feat: register leo platform"
```

---

### Task 2: Define and validate Leo account credentials

**Files:**
- Create: `backend/internal/service/leo_account.go`
- Create: `backend/internal/service/leo_account_test.go`
- Modify: `backend/internal/service/account.go`
- Modify: `backend/internal/service/admin_account.go`
- Modify: `backend/internal/service/admin_group.go`
- Modify: `backend/internal/service/admin_service_group_test.go`
- Modify: `progress.md`

- [ ] **Step 1: Write failing Leo URL and credential tests**

Create table-driven tests covering:

```go
func TestNormalizeLeoBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{"docker http", "http://leostudio:8000/v1/", "http://leostudio:8000/v1", false},
		{"private ip", "http://10.0.0.20:8000/v1", "http://10.0.0.20:8000/v1", false},
		{"https", "https://leo.internal.example/v1", "https://leo.internal.example/v1", false},
		{"missing v1", "http://leostudio:8000", "", true},
		{"userinfo", "http://user:pass@leostudio:8000/v1", "", true},
		{"query", "http://leostudio:8000/v1?x=1", "", true},
		{"fragment", "http://leostudio:8000/v1#x", "", true},
		{"bad scheme", "file:///tmp/v1", "", true},
	}
	// Run NormalizeLeoBaseURL and assert want/wantErr for every row.
}

func TestValidateLeoAccountCredentials(t *testing.T) {
	valid := map[string]any{
		"base_url": "http://leostudio:8000/v1",
		"api_key": "leo-secret",
		"model_mapping": map[string]any{
			"seedance-2.0": "seedance-2.0",
		},
	}
	require.NoError(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeAPIKey, valid))
	require.Error(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeOAuth, valid))
	require.Error(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeAPIKey, map[string]any{}))
}
```

- [ ] **Step 2: Run tests and verify they fail**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run 'TestNormalizeLeoBaseURL|TestValidateLeoAccountCredentials' -count=1
```

Expected: FAIL because Leo helpers do not exist.

- [ ] **Step 3: Implement the minimal Leo account contract**

Create `leo_account.go` with these public helpers:

```go
var LeoDefaultVideoModelIDs = []string{"seedance-2.0", "seedance-2.0-fast"}

func (a *Account) IsLeo() bool {
	return a != nil && a.Platform == PlatformLeo
}

func (a *Account) GetLeoAPIKey() string {
	if !a.IsLeo() || a.Type != AccountTypeAPIKey {
		return ""
	}
	return strings.TrimSpace(a.GetCredential("api_key"))
}

func (a *Account) GetLeoBaseURL() string {
	if !a.IsLeo() {
		return ""
	}
	baseURL, err := NormalizeLeoBaseURL(a.GetCredential("base_url"))
	if err != nil {
		return ""
	}
	return baseURL
}

func BuildLeoVideosGenerationsURL(baseURL string) (string, error) {
	baseURL, err := NormalizeLeoBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	return baseURL + "/videos/generations", nil
}

func BuildLeoHealthURL(baseURL string) (string, error) {
	baseURL, err := NormalizeLeoBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(baseURL, "/v1") + "/health", nil
}
```

`NormalizeLeoBaseURL` must use `net/url`, accept only `http`/`https`, require a host and exact `/v1` path, reject userinfo/query/fragment, and trim only the trailing slash. `ValidateLeoAccountCredentials` must be a no-op for non-Leo platforms and must reject non-API-key Leo accounts, missing URLs, missing keys, and empty/non-string model mappings with `infraerrors.BadRequest`.

- [ ] **Step 4: Enforce validation on create and merged update**

In `CreateAccount`, call the validator before constructing the account:

```go
if err := ValidateLeoAccountCredentials(input.Platform, input.Type, input.Credentials); err != nil {
	return nil, err
}
```

In `UpdateAccount`, call it after sensitive credential preservation and after applying an optional type change:

```go
if err := ValidateLeoAccountCredentials(account.Platform, account.Type, account.Credentials); err != nil {
	return nil, err
}
```

Add `IsLeo` to `Account.IsOpenAICompatible` only so Leo can reuse the existing scheduler; do not add Leo to chat, Responses, images, OAuth, quota-refresh, or WebSocket branches.

- [ ] **Step 5: Add default Leo model candidates**

Extend `defaultModelsListCandidateIDs`:

```go
case PlatformLeo:
	return append([]string(nil), LeoDefaultVideoModelIDs...)
```

- [ ] **Step 6: Run focused account and group tests**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run 'Leo|ModelMapping|GroupModelsList' -count=1
```

Expected: PASS.

- [ ] **Step 7: Append Task 2 to `progress.md` and commit**

```powershell
git add backend/internal/service/leo_account.go backend/internal/service/leo_account_test.go backend/internal/service/account.go backend/internal/service/admin_account.go backend/internal/service/admin_group.go backend/internal/service/admin_service_group_test.go
git commit -m "feat: validate leo accounts"
```

---

### Task 3: Add Bearer-authenticated Leo health checks

**Files:**
- Modify: `backend/internal/service/account_test_service.go`
- Create: `backend/internal/service/account_test_service_leo_test.go`
- Modify: `progress.md`

- [ ] **Step 1: Write the failing health-check test**

Use `httptest.NewServer` and an account repository stub. The server must assert:

```go
require.Equal(t, http.MethodGet, r.Method)
require.Equal(t, "/health", r.URL.Path)
require.Equal(t, "Bearer leo-secret", r.Header.Get("Authorization"))
require.Equal(t, "application/json", r.Header.Get("Accept"))
```

Return `{"status":"ok"}` and assert the account test stream ends with `test_complete` and `success=true`. Add sibling tests for `401`, invalid `base_url`, and a transport failure.

- [ ] **Step 2: Run the test and verify it routes incorrectly before implementation**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run TestAccountTestService_Leo -count=1
```

Expected: FAIL because Leo currently falls through to the Claude probe.

- [ ] **Step 3: Route Leo accounts to a dedicated health probe**

Add before the existing OpenAI branch:

```go
if account.IsLeo() {
	return s.testLeoAccountConnection(c, account)
}
```

Implement `testLeoAccountConnection` using `BuildLeoHealthURL`, `s.httpUpstream`, the account proxy, and the configured concurrency. On `200`, emit a short content event and a successful completion event. On `401/403`, call `accountRepo.SetError` with a redacted message. Do not include the API key or full response body in emitted errors.

- [ ] **Step 4: Run the focused health tests**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run TestAccountTestService_Leo -count=1
```

Expected: PASS.

- [ ] **Step 5: Append Task 3 to `progress.md` and commit**

```powershell
git add backend/internal/service/account_test_service.go backend/internal/service/account_test_service_leo_test.go
git commit -m "feat: test leo upstream health"
```

---

### Task 4: Implement synchronous Leo video forwarding

**Files:**
- Create: `backend/internal/service/leo_video.go`
- Create: `backend/internal/service/leo_video_test.go`
- Modify: `backend/internal/service/video_billing_resolution.go`
- Create: `backend/internal/service/video_billing_resolution_test.go`
- Modify: `progress.md`

- [ ] **Step 1: Write failing request and response tests**

Cover these cases with a recording `HTTPUpstream` stub:

```go
func TestForwardLeoVideoMapsModelAndAddsBearer(t *testing.T) {
	body := []byte(`{"model":"seedance","prompt":"city","resolution":"720p","duration":8,"audio":false}`)
	account := &Account{
		ID: 1, Platform: PlatformLeo, Type: AccountTypeAPIKey, Concurrency: 1,
		Credentials: map[string]any{
			"base_url": "http://leo.internal:8000/v1",
			"api_key": "leo-secret",
			"model_mapping": map[string]any{"seedance": "seedance-2.0"},
		},
	}
	// Assert POST /v1/videos/generations, Bearer header, mapped model,
	// and unchanged prompt/resolution/duration/audio fields.
}
```

For a successful response containing `provider.resolution="RESOLUTION_720"` and `provider.duration=12`, assert:

```go
require.Equal(t, 1, result.VideoCount)
require.Equal(t, "720p", result.VideoResolution)
require.Equal(t, 12, result.VideoDurationSeconds)
require.Equal(t, "seedance", result.Model)
require.Equal(t, "seedance-2.0", result.UpstreamModel)
```

Add a fallback test where provider metadata is missing and request resolution/duration are used. Add tests that `400/422` are written through without failover and that `401/403/429/5xx` return `UpstreamFailoverError` without committing the response.

- [ ] **Step 2: Run tests and verify the forwarding method is missing**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run 'TestForwardLeoVideo|TestNormalizeVideoBillingResolutionLeo' -count=1
```

Expected: FAIL because `ForwardLeoVideo` and Leo resolution aliases do not exist.

- [ ] **Step 3: Implement request parsing and model rewriting**

Define a focused request type:

```go
type LeoVideoRequestInfo struct {
	Model           string
	Prompt          string
	Resolution      string
	DurationSeconds int
	ImageURL        string
}
```

Parse with `gjson`; map only the `model` field using `account.ResolveMappedModel`, and rewrite only that field with `sjson.SetBytes`. Reject invalid JSON before creating the upstream request.

- [ ] **Step 4: Implement the synchronous forwarder**

The method signature is:

```go
func (s *OpenAIGatewayService) ForwardLeoVideo(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
) (*OpenAIForwardResult, error)
```

Required behavior:

```go
targetURL, err := BuildLeoVideosGenerationsURL(account.GetLeoBaseURL())
// Build POST request with detached upstream context.
req.Header.Set("Authorization", "Bearer "+account.GetLeoAPIKey())
req.Header.Set("Accept", "application/json")
req.Header.Set("Content-Type", "application/json")
// Use account proxy and concurrency through s.httpUpstream.Do.
```

Use `ReadUpstreamResponseBody` and the existing filtered passthrough response-header writer. For success, write the response and build `OpenAIForwardResult` with `VideoCount`, normalized actual response resolution/duration, response generation ID, requested model, mapped upstream model, response headers, and total duration.

For `401/403`, call `rateLimitService.HandleUpstreamError` so API-key accounts become error-state accounts, then return `UpstreamFailoverError`. For `429/5xx`, reconcile the existing rate-limit state and return `UpstreamFailoverError`. For `400/422`, write the sanitized upstream response and return `(nil, nil)` so no usage is charged.

- [ ] **Step 5: Recognize LeoStudio resolution values**

Extend `NormalizeVideoBillingResolutionOrDefault` cases:

```go
case "480", "480p", "sd", "resolution_480":
	return VideoBillingResolution480P
case "720", "720p", "hd", "resolution_720":
	return VideoBillingResolution720P
case "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080":
	return VideoBillingResolution1080P
```

- [ ] **Step 6: Run service tests**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run 'LeoVideo|NormalizeVideoBilling' -count=1
```

Expected: PASS.

- [ ] **Step 7: Append Task 4 to `progress.md` and commit**

```powershell
git add backend/internal/service/leo_video.go backend/internal/service/leo_video_test.go backend/internal/service/video_billing_resolution.go backend/internal/service/video_billing_resolution_test.go
git commit -m "feat: forward leo video requests"
```

---

### Task 5: Route, schedule, and fail over Leo video requests

**Files:**
- Create: `backend/internal/handler/leo_video.go`
- Create: `backend/internal/handler/leo_video_test.go`
- Modify: `backend/internal/server/routes/gateway.go`
- Modify: `backend/internal/service/openai_gateway_scheduling.go`
- Modify: `backend/internal/service/openai_account_scheduler_test.go`
- Modify: `backend/internal/handler/endpoint_test.go`
- Modify: `progress.md`

- [ ] **Step 1: Write failing scheduler isolation and route tests**

Add a scheduler test with one OpenAI, one Grok, and one Leo account. Call `SelectAccountWithSchedulerForCapability` with `PlatformLeo` and assert only the Leo account is selected.

Add route assertions:

```go
// Leo group
POST /v1/videos/generations -> LeoVideoGeneration
POST /videos/generations    -> LeoVideoGeneration

// Leo group
POST /v1/videos/edits       -> 404 not_found_error
GET  /v1/videos/:id         -> 404 not_found_error
POST /v1/images/generations -> 404 not_found_error
```

- [ ] **Step 2: Run tests and verify Leo normalizes to OpenAI before the fix**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service ./internal/server/routes -run 'Leo|VideoGeneration' -count=1
```

Expected: FAIL because the scheduler and route dispatcher do not preserve Leo.

- [ ] **Step 3: Preserve Leo in the shared scheduler**

Change the normalizer to:

```go
func normalizeOpenAICompatiblePlatform(platform string) string {
	switch platform {
	case PlatformGrok, PlatformLeo:
		return platform
	default:
		return PlatformOpenAI
	}
}
```

`Account.IsOpenAICompatible` already includes Leo from Task 2; keep all OpenAI endpoint capability arguments empty for Leo so no chat/image capability is inferred.

- [ ] **Step 4: Implement `LeoVideoGeneration` on `OpenAIGatewayHandler`**

The handler must follow the existing Grok media lifecycle without importing Grok request normalization:

```go
func (h *OpenAIGatewayHandler) LeoVideoGeneration(c *gin.Context) {
	// Recover panic; load API key and auth subject.
	// Read and validate JSON body; require model and prompt.
	// Check billing eligibility and GroupAllowsImageGeneration.
	// Acquire user and image-generation slots.
	// Select PlatformLeo account for the requested model.
	// Acquire the selected account concurrency slot.
	// Call ForwardLeoVideo.
	// Retry another account only for UpstreamFailoverError.
	// Submit OpenAIRecordUsageInput after a successful result.
}
```

Use `GenerateExplicitSessionHash(c, body)` for deterministic balancing, `PlatformLeo` in scheduler selection and no sticky response ID. Use existing `maxAccountSwitches`, `failedAccountIDs`, `handleFailoverExhausted`, ops timing, request logging, and `ChannelUsageFields`.

- [ ] **Step 5: Dispatch only generation routes to Leo**

Update the route switch:

```go
switch getGroupPlatform(c) {
case service.PlatformGrok:
	h.OpenAIGateway.GrokVideoGeneration(c)
case service.PlatformLeo:
	h.OpenAIGateway.LeoVideoGeneration(c)
default:
	// Existing not-supported response.
}
```

Leave status, edit, extension, image, messages, Responses, embeddings, and WebSocket dispatch unchanged so they reject Leo.

- [ ] **Step 6: Run handler and route tests**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/handler ./internal/server/routes ./internal/service -run 'Leo|VideoGeneration' -count=1
```

Expected: PASS.

- [ ] **Step 7: Append Task 5 to `progress.md` and commit**

```powershell
git add backend/internal/handler/leo_video.go backend/internal/handler/leo_video_test.go backend/internal/server/routes/gateway.go backend/internal/service/openai_gateway_scheduling.go backend/internal/service/openai_account_scheduler_test.go backend/internal/handler/endpoint_test.go
git commit -m "feat: route leo video generation"
```

---

### Task 6: Enforce Leo video pricing and record actual output metadata

**Files:**
- Modify: `backend/internal/service/admin_group.go`
- Modify: `backend/internal/service/admin_service_group_test.go`
- Modify: `backend/internal/service/openai_gateway_usage.go`
- Modify: `backend/internal/service/openai_gateway_record_usage_test.go`
- Create: `backend/internal/handler/admin/group_handler_leo_test.go`
- Modify: `progress.md`

- [ ] **Step 1: Write failing group-price tests**

Add create tests proving all three prices are required for Leo and zero is valid:

```go
zero := 0.0
one := 0.1

_, err := svc.CreateGroup(ctx, &CreateGroupInput{
	Name: "leo", Platform: PlatformLeo, RateMultiplier: 1,
	VideoPrice480P: &one, VideoPrice720P: &one,
})
require.ErrorContains(t, err, "video_price_1080p is required")

group, err := svc.CreateGroup(ctx, &CreateGroupInput{
	Name: "leo-free", Platform: PlatformLeo, RateMultiplier: 1,
	VideoPrice480P: &zero, VideoPrice720P: &zero, VideoPrice1080P: &zero,
})
require.NoError(t, err)
require.NotNil(t, group.VideoPrice480P)
```

Add update tests proving a Leo group cannot clear one tier and cannot become Leo without all three effective prices.

- [ ] **Step 2: Write the failing generic video-usage test**

Create an `OpenAIForwardResult` with model `seedance-2.0`, `VideoCount=1`, `VideoResolution="720p"`, and `VideoDurationSeconds=12`. Assert the usage log has billing mode `video`, 720p, duration 12, and cost `group_720_price * 12 * video_multiplier`.

- [ ] **Step 3: Run tests and verify current Grok-only billing detection fails**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service -run 'LeoGroupVideoPrice|LeoVideoUsage' -count=1
```

Expected: FAIL because group prices are optional and `isGrokVideoUsageResult` rejects Seedance.

- [ ] **Step 4: Add effective Leo group-price validation**

Implement:

```go
func validateLeoVideoPrices(group *Group) error {
	if group == nil || group.Platform != PlatformLeo {
		return nil
	}
	if group.VideoPrice480P == nil {
		return errors.New("video_price_480p is required for leo groups")
	}
	if group.VideoPrice720P == nil {
		return errors.New("video_price_720p is required for leo groups")
	}
	if group.VideoPrice1080P == nil {
		return errors.New("video_price_1080p is required for leo groups")
	}
	return nil
}
```

Call it after normalized prices are assigned in `CreateGroup`, and after the final merged state is assembled in `UpdateGroup`, immediately before repository persistence.

- [ ] **Step 5: Generalize video billing detection**

Replace the Grok-name gate with explicit result metadata:

```go
func isVideoUsageResult(result *OpenAIForwardResult) bool {
	return result != nil && result.VideoCount > 0
}
```

Update all internal calls in `openai_gateway_usage.go`. Keep Grok default pricing helpers unchanged; Leo always reaches configured group video prices because group validation makes the relevant tier non-nil.

- [ ] **Step 6: Run group and usage tests**

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/service ./internal/handler/admin -run 'LeoGroupVideoPrice|LeoVideoUsage|GrokVideo' -count=1
```

Expected: PASS, including existing Grok video billing regressions.

- [ ] **Step 7: Append Task 6 to `progress.md` and commit**

```powershell
git add backend/internal/service/admin_group.go backend/internal/service/admin_service_group_test.go backend/internal/service/openai_gateway_usage.go backend/internal/service/openai_gateway_record_usage_test.go backend/internal/handler/admin
git commit -m "feat: bill leo videos by output"
```

---

### Task 7: Add complete Leo admin UI support

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/utils/platformColors.ts`
- Modify: `frontend/src/api/admin/settings.ts`
- Modify: `frontend/src/api/admin/users.ts`
- Modify: `frontend/src/api/__tests__/settings.authSourceDefaults.spec.ts`
- Modify: `frontend/src/components/user/UserPlatformQuotaCell.vue`
- Modify: `frontend/src/components/user/dashboard/UserDashboardStats.vue`
- Modify: `frontend/src/components/common/PlatformIcon.vue`
- Modify: `frontend/src/components/common/PlatformTypeBadge.vue`
- Modify: `frontend/src/components/common/GroupBadge.vue`
- Modify: `frontend/src/components/admin/account/AccountTableFilters.vue`
- Modify: `frontend/src/components/admin/user/UserPlatformQuotaModal.vue`
- Modify: `frontend/src/components/admin/user/__tests__/UserPlatformQuotaModal.spec.ts`
- Modify: `frontend/src/components/admin/ErrorPassthroughRulesModal.vue`
- Modify: `frontend/src/components/admin/channel/types.ts`
- Modify: `frontend/src/views/admin/GroupsView.vue`
- Modify: `frontend/src/views/admin/groupsImagePricing.ts`
- Modify: `frontend/src/views/admin/__tests__/groupsImagePricing.spec.ts`
- Modify: `frontend/src/views/admin/ops/components/OpsDashboardHeader.vue`
- Modify: `frontend/src/components/account/CreateAccountModal.vue`
- Modify: `frontend/src/components/account/EditAccountModal.vue`
- Modify: `frontend/src/components/account/__tests__/CreateAccountModal.spec.ts`
- Modify: `frontend/src/components/account/__tests__/EditAccountModal.spec.ts`
- Modify: `frontend/src/composables/useModelWhitelist.ts`
- Modify: `frontend/src/i18n/locales/zh/admin/accounts.ts`
- Modify: `frontend/src/i18n/locales/en/admin/accounts.ts`
- Modify: `frontend/src/i18n/locales/zh/admin/overview.ts`
- Modify: `frontend/src/i18n/locales/en/admin/overview.ts`
- Modify: `progress.md`

- [ ] **Step 1: Write failing frontend platform and form tests**

Extend platform pricing tests:

```ts
expect(supportsVideoPricingPlatform('leo')).toBe(true)
expect(supportsImagePricingPlatform('leo')).toBe(false)
expect(getVideoPricePlaceholder('leo', 'video_price_480p')).toBe('')
expect(getDefaultVideoPreviewPrice('leo', 'video_price_480p')).toBeNull()
```

Add mounted component tests that select Leo and assert:

```ts
expect(wrapper.get('[data-testid="platform-leo"]').exists()).toBe(true)
expect(wrapper.get('[data-testid="leo-base-url"]').attributes('required')).toBeDefined()
expect(wrapper.get('[data-testid="leo-api-key"]').attributes('type')).toBe('password')
expect(wrapper.text()).toContain('seedance-2.0')
expect(wrapper.text()).not.toContain('OAuth')
```

For edit, assert an existing redacted Leo account preserves the stored API key when the password input stays empty and submits a replacement only when entered.

- [ ] **Step 2: Run focused frontend tests and verify they fail**

Run from `frontend/`:

```powershell
pnpm exec vitest run src/views/admin/__tests__/groupsImagePricing.spec.ts src/components/account/__tests__/CreateAccountModal.spec.ts src/components/account/__tests__/EditAccountModal.spec.ts src/api/__tests__/settings.authSourceDefaults.spec.ts
```

Expected: FAIL because Leo is absent from types, options, and forms.

- [ ] **Step 3: Add Leo to shared frontend platform identity**

Extend unions and normalization arrays:

```ts
export type GroupPlatform = 'anthropic' | 'openai' | 'gemini' | 'antigravity' | 'grok' | 'leo'
export type AccountPlatform = GroupPlatform
export type PlatformType = 'anthropic' | 'openai' | 'gemini' | 'antigravity' | 'grok' | 'leo'
```

Append Leo after Grok in quota/usage ordering arrays and account/group filter options. Add explicit `Leo` labels and a restrained teal/neutral badge class. Reuse the existing icon system with a familiar play/video mark and do not add a custom brand SVG.

- [ ] **Step 4: Add Leo group controls and validation**

Change `supportsVideoPricingPlatform`:

```ts
export const supportsVideoPricingPlatform = (platform: string): boolean =>
  platform === 'grok' || platform === 'leo'
```

Do not add Leo to `imagePricingPlatforms`. In `GroupsView.vue`, add Leo to create/edit/filter platform options. For Leo, render the legacy `allow_image_generation` checkbox with Leo-specific text “允许视频生成 / Allow video generation”. Before create or update submission, require finite non-negative values for all three video price fields:

```ts
const hasCompleteLeoVideoPrices = (form: VideoPricingFormState) =>
  form.platform !== 'leo' ||
  [form.video_price_480p, form.video_price_720p, form.video_price_1080p]
    .every((value) => value !== null && value !== '' && Number.isFinite(Number(value)) && Number(value) >= 0)
```

Show field-level or existing toast validation and do not substitute Grok placeholder prices.

- [ ] **Step 5: Add Leo create and edit account forms**

Add a `data-testid="platform-leo"` platform button. Platform switching must force:

```ts
accountCategory.value = 'apikey'
form.type = 'apikey'
apiKeyBaseUrl.value = ''
modelRestrictionMode.value = 'mapping'
modelMappings.value = [
  { from: 'seedance-2.0', to: 'seedance-2.0' },
  { from: 'seedance-2.0-fast', to: 'seedance-2.0-fast' },
]
```

For Leo, require a non-empty Base URL, API Key, and at least one valid mapping. Use `http://leostudio:8000/v1` only as placeholder copy, not as an implicit submitted default. Keep the existing sensitive credential preservation behavior in `EditAccountModal.vue`.

- [ ] **Step 6: Add concise bilingual copy**

Add keys for Leo Base URL, Bearer API Key, required model mapping, video-only scope, and manual USD/second prices. Update the generic video pricing description so it branches by platform: Grok retains official defaults; Leo states that all prices are operator supplied.

- [ ] **Step 7: Run focused tests, typecheck, and lint**

```powershell
pnpm exec vitest run src/views/admin/__tests__/groupsImagePricing.spec.ts src/components/account/__tests__/CreateAccountModal.spec.ts src/components/account/__tests__/EditAccountModal.spec.ts src/api/__tests__/settings.authSourceDefaults.spec.ts
pnpm run typecheck
pnpm run lint:check
```

Expected: all commands exit 0.

- [ ] **Step 8: Append Task 7 to `progress.md` and commit**

```powershell
git add frontend/src/types frontend/src/utils/platformColors.ts frontend/src/api/admin/settings.ts frontend/src/api/admin/users.ts frontend/src/api/__tests__/settings.authSourceDefaults.spec.ts frontend/src/components/user frontend/src/components/common frontend/src/components/admin/account/AccountTableFilters.vue frontend/src/components/admin/user frontend/src/components/admin/ErrorPassthroughRulesModal.vue frontend/src/components/admin/channel/types.ts frontend/src/views/admin/GroupsView.vue frontend/src/views/admin/groupsImagePricing.ts frontend/src/views/admin/__tests__/groupsImagePricing.spec.ts frontend/src/views/admin/ops/components/OpsDashboardHeader.vue frontend/src/components/account frontend/src/composables/useModelWhitelist.ts frontend/src/i18n/locales/zh/admin frontend/src/i18n/locales/en/admin
git commit -m "feat: configure leo channels in admin"
```

---

### Task 8: Document, integrate, and verify the complete workflow

**Files:**
- Create: `docs/LEO_VIDEO_CHANNEL.md`
- Create: `backend/internal/handler/leo_video_integration_test.go`
- Modify: `progress.md`
- Test: all backend and frontend suites affected above

- [ ] **Step 1: Write the operator documentation**

Document:

```text
LeoStudio requirements:
- GET /health
- POST /v1/videos/generations
- Authorization: Bearer <api_key>

Sub2API setup:
- create a Leo group
- enable video generation
- configure 480p/720p/1080p USD-per-second prices
- create a Leo API-key account with base_url, api_key, model_mapping, concurrency
- run the account health test
```

Include a working request example:

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

State explicitly that the call is synchronous, can take several minutes, continues upstream after client disconnect, and that Leo does not support image/edit/extend/status/chat endpoints in this release.

- [ ] **Step 2: Run all backend tests**

Run from `backend/`:

```powershell
$env:GOTOOLCHAIN='auto'; go test ./...
```

Expected: PASS with zero failed packages.

- [ ] **Step 3: Run frontend focused tests and full static verification**

Run from `frontend/`:

```powershell
pnpm exec vitest run src/views/admin/__tests__/groupsImagePricing.spec.ts src/components/account/__tests__/CreateAccountModal.spec.ts src/components/account/__tests__/EditAccountModal.spec.ts src/api/__tests__/settings.authSourceDefaults.spec.ts
pnpm run typecheck
pnpm run lint:check
pnpm run build
```

Expected: all commands exit 0; existing non-fatal Browserslist/chunk warnings may remain.

- [ ] **Step 4: Write and run the repeatable mock-LeoStudio integration test**

Create `backend/internal/handler/leo_video_integration_test.go`. The test must start `httptest.NewServer` that:

- rejects a wrong Bearer token;
- returns `{"status":"ok"}` for `/health`;
- validates the mapped model at `/v1/videos/generations`;
- waits briefly and returns LeoStudio-shaped `data` and `provider` fields.

Wire the handler with in-memory repository/cache stubs, a Leo group, a Leo account, and a Sub2API API key. Call the public video endpoint and verify:

```text
HTTP 200
response data[0].mp4_url preserved
provider.generation_id preserved
usage platform = leo
usage billing_mode = video
usage resolution/duration match provider output
charged amount = configured USD/s × duration × multiplier
```

Run from `backend/`:

```powershell
$env:GOTOOLCHAIN='auto'; go test ./internal/handler -run TestLeoVideoGenerationIntegration -count=1
```

Expected: PASS.

- [ ] **Step 5: Inspect the final diff and ensure scope remains exact**

```powershell
git diff --check
git status --short
git diff --stat
```

Expected: no whitespace errors, no migration files, no async task code, and no unrelated formatting churn.

- [ ] **Step 6: Append final verification and rollback instructions to `progress.md`**

List every changed file with a one-line purpose. The rollback point is the commit immediately before Task 1; rollback by reverting the Leo task commits in reverse order. No database rollback is required because there is no migration.

- [ ] **Step 7: Commit documentation and final integration changes**

```powershell
git add docs/LEO_VIDEO_CHANNEL.md backend/internal/handler/leo_video_integration_test.go
git commit -m "docs: explain leo video channels"
```

---

## Final Review Checklist

- [ ] `leo` is independent from Grok in persisted platform values and route dispatch.
- [ ] Only synchronous video generation is supported.
- [ ] Every Leo upstream request uses the configured Bearer API Key.
- [ ] Private HTTP Base URLs work only through validated admin-managed Leo credentials.
- [ ] API Key values are redacted in admin responses and logs.
- [ ] Scheduler selection never mixes Leo, OpenAI, and Grok accounts.
- [ ] Model mapping rewrites only `model`.
- [ ] `400/422` do not fail over or charge; `401/403/429/5xx` follow the approved behavior.
- [ ] Actual response resolution and duration drive billing when available.
- [ ] Leo groups cannot fall back to Grok or image default prices.
- [ ] Unsupported endpoints return clear 404 responses.
- [ ] No database migration or async task system was added.
- [ ] Backend tests, frontend tests, typecheck, lint, build, and local smoke test all pass.
