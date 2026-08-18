package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type splitOpenAIImageTaskResult struct {
	index   int
	account *service.Account
	result  *service.OpenAIForwardResult
	body    []byte
	header  http.Header
	status  int
	err     error
}

// handleOpenAIImagesNSplit fans out n>1 non-streaming image requests when the
// selected upstream path cannot be relied on to implement native batching.
// Every child request still goes through account scheduling and the account
// concurrency slot, which preserves pool-mode capacity accounting.
func (h *OpenAIGatewayHandler) handleOpenAIImagesNSplit(
	c *gin.Context,
	apiKey *service.APIKey,
	parsed *service.OpenAIImagesRequest,
	body []byte,
	requestCtx context.Context,
	sessionHash string,
	routingModel string,
	requestModel string,
	channelMapping service.ChannelMappingResult,
	subscription *service.UserSubscription,
	reqLog *zap.Logger,
) bool {
	if parsed == nil || parsed.N <= 1 || parsed.Stream {
		return false
	}

	// The public request remains one user operation; child tasks only consume
	// account slots. Use independent session seeds so multiple accounts can be
	// selected, while pool-mode accounts still receive their own sticky hash.
	tasks := make([]splitOpenAIImageTaskResult, parsed.N)
	var wg sync.WaitGroup
	for index := range tasks {
		tasks[index].index = index
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			result := h.runOpenAIImagesNSplitTask(
				c,
				apiKey.GroupID,
				parsed,
				body,
				requestCtx,
				splitImageSessionHash(sessionHash, index),
				routingModel,
				requestModel,
				channelMapping.MappedModel,
				effectiveAPIKeyPlatform(c, apiKey),
				reqLog,
			)
			result.index = index
			tasks[index] = result
		}(index)
	}
	wg.Wait()

	parts := make([][]byte, 0, len(tasks))
	successes := make([]splitOpenAIImageTaskResult, 0, len(tasks))
	for _, task := range tasks {
		if task.result == nil || task.result.ImageCount <= 0 || len(task.body) == 0 {
			continue
		}
		parts = append(parts, task.body)
		successes = append(successes, task)
	}
	if len(successes) == 0 {
		for _, task := range tasks {
			var leoImageErr *service.LeoImageRequestError
			if errors.As(task.err, &leoImageErr) {
				h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", leoImageErr.Error())
				return true
			}
			var imageErr *service.OpenAIImagesUpstreamError
			if errors.As(task.err, &imageErr) {
				status := imageErr.StatusCode
				if status <= 0 {
					status = http.StatusBadGateway
				}
				errType := imageErr.ErrorType
				if strings.TrimSpace(errType) == "" {
					errType = "upstream_error"
				}
				message := imageErr.ClientMessage()
				if strings.TrimSpace(message) == "" {
					message = "Upstream request failed"
				}
				h.errorResponse(c, status, errType, message)
				return true
			}
		}
		h.handleStreamingAwareError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed", false)
		return true
	}

	merged, err := service.MergeOpenAIImageResponses(parts)
	if err != nil {
		reqLog.Warn("openai.images.n_split_response_merge_failed", zap.Error(err))
		h.handleStreamingAwareError(c, http.StatusBadGateway, "upstream_error", "Upstream response could not be combined", false)
		return true
	}
	if len(successes) < len(tasks) {
		reqLog.Warn("openai.images.n_split_partial_success",
			zap.Int("requested_count", len(tasks)),
			zap.Int("success_count", len(successes)),
		)
	}

	first := successes[0]
	for key, values := range first.header {
		if strings.EqualFold(key, "Content-Length") {
			continue
		}
		c.Writer.Header()[key] = append([]string(nil), values...)
	}
	status := first.status
	if status <= 0 {
		status = http.StatusOK
	}
	c.Data(status, "application/json", merged)

	for _, task := range successes {
		h.recordSplitOpenAIImageUsage(splitImageUsageContext(c.Request.Context(), task.index), c, apiKey, subscription, task.account, task.result, requestModel, body, channelMapping, reqLog)
	}
	return true
}

func (h *OpenAIGatewayHandler) runOpenAIImagesNSplitTask(
	parent *gin.Context,
	groupID *int64,
	parsed *service.OpenAIImagesRequest,
	body []byte,
	requestCtx context.Context,
	sessionHash string,
	routingModel string,
	requestModel string,
	channelMappedModel string,
	platform string,
	reqLog *zap.Logger,
) splitOpenAIImageTaskResult {
	child, recorder := newBufferedGatewayContext(parent)
	childParsed := *parsed
	childParsed.N = 1
	failedAccountIDs := make(map[int64]struct{})
	sameAccountRetryCount := make(map[int64]int)
	maxAccountSwitches := h.maxAccountSwitches
	if maxAccountSwitches <= 0 {
		maxAccountSwitches = 3
	}
	switchCount := 0
	var lastFailoverErr *service.UpstreamFailoverError
	var oauth429FailoverState service.OpenAIOAuth429FailoverState

	for {
		if err := requestCtx.Err(); err != nil {
			return splitOpenAIImageTaskResult{err: err}
		}
		selection, _, err := h.gatewayService.SelectAccountWithSchedulerForImages(
			requestCtx,
			groupID,
			sessionHash,
			routingModel,
			failedAccountIDs,
			childParsed.RequiredCapability,
			platform,
		)
		if err != nil || selection == nil || selection.Account == nil {
			if lastFailoverErr != nil {
				return splitOpenAIImageTaskResult{err: lastFailoverErr}
			}
			if err == nil {
				err = service.ErrNoAvailableAccounts
			}
			return splitOpenAIImageTaskResult{err: err}
		}

		account := selection.Account
		sessionHash = h.bindSplitImagePoolSession(requestCtx, groupID, sessionHash, account, reqLog)
		childStreamStarted := false
		release, slotResult := h.acquireResponsesAccountSlot(child, groupID, sessionHash, selection, false, &childStreamStarted, reqLog)
		if slotResult == openAISlotAcquireProfitVetoed {
			failedAccountIDs[account.ID] = struct{}{}
			continue
		}
		if slotResult != openAISlotAcquireOK {
			return splitOpenAIImageTaskResult{err: fmt.Errorf("image account slot unavailable")}
		}

		result, forwardErr := func() (*service.OpenAIForwardResult, error) {
			defer func() {
				if release != nil {
					release()
				}
			}()
			return h.gatewayService.ForwardImages(requestCtx, child, account, body, &childParsed, channelMappedModel)
		}()
		if result != nil && result.ImageCount > 0 && recorder.Body.Len() > 0 {
			if account.Type == service.AccountTypeOAuth && !account.IsShadow() {
				h.gatewayService.UpdateCodexUsageSnapshotFromHeaders(child.Request.Context(), account.ID, result.ResponseHeaders)
			}
			h.gatewayService.ReportOpenAIAccountScheduleResult(account.ID, account.GetMappedModel(requestModel), true, result.FirstTokenMs)
			return splitOpenAIImageTaskResult{
				account: account,
				result:  result,
				body:    append([]byte(nil), recorder.Body.Bytes()...),
				header:  recorder.Header().Clone(),
				status:  recorder.Code,
			}
		}
		if forwardErr == nil {
			forwardErr = fmt.Errorf("upstream returned no image output")
		}

		var failoverErr *service.UpstreamFailoverError
		if !errors.As(forwardErr, &failoverErr) {
			h.gatewayService.ReportOpenAIAccountScheduleResult(account.ID, account.GetMappedModel(requestModel), false, nil)
			return splitOpenAIImageTaskResult{err: forwardErr}
		}
		h.gatewayService.ReportOpenAIAccountScheduleResult(account.ID, account.GetMappedModel(requestModel), false, nil)
		if failoverErr.RetryableOnSameAccount {
			retryLimit := account.GetPoolModeRetryCount()
			if sameAccountRetryCount[account.ID] < retryLimit {
				sameAccountRetryCount[account.ID]++
				select {
				case <-requestCtx.Done():
					return splitOpenAIImageTaskResult{err: requestCtx.Err()}
				case <-time.After(sameAccountRetryDelay):
				}
				continue
			}
		}
		h.gatewayService.RecordOpenAIAccountSwitch()
		failedAccountIDs[account.ID] = struct{}{}
		lastFailoverErr = failoverErr
		if switchCount >= maxAccountSwitches {
			return splitOpenAIImageTaskResult{err: lastFailoverErr}
		}
		switchCount++
		if h.gatewayService.ShouldStopOpenAIOAuth429Failover(account, failoverErr.StatusCode, switchCount, &oauth429FailoverState) {
			return splitOpenAIImageTaskResult{err: lastFailoverErr}
		}
	}
}

func newBufferedGatewayContext(parent *gin.Context) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	child, _ := gin.CreateTestContext(recorder)
	if parent == nil {
		return child, recorder
	}
	copyContext := parent.Copy()
	child.Request = copyContext.Request
	child.Params = copyContext.Params
	child.Keys = copyContext.Keys
	return child, recorder
}

func splitImageSessionHash(base string, index int) string {
	if strings.TrimSpace(base) == "" {
		return ""
	}
	return service.DeriveSessionHashFromSeed(fmt.Sprintf("%s:image:%d", base, index))
}

func (h *OpenAIGatewayHandler) bindSplitImagePoolSession(
	ctx context.Context,
	groupID *int64,
	sessionHash string,
	account *service.Account,
	reqLog *zap.Logger,
) string {
	sessionHash = ensureOpenAIPoolModeSessionHash(sessionHash, account)
	if account == nil || !account.IsPoolMode() || sessionHash == "" {
		return sessionHash
	}
	if err := h.gatewayService.BindStickySession(ctx, groupID, sessionHash, account.ID); err != nil {
		reqLog.Warn("images.n_split_pool_session_bind_failed", zap.Int64("account_id", account.ID), zap.Error(err))
	}
	return sessionHash
}

func splitImageUsageContext(parent context.Context, index int) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	suffix := fmt.Sprintf(":image:%d", index+1)
	if value, _ := parent.Value(ctxkey.ClientRequestID).(string); strings.TrimSpace(value) != "" {
		parent = context.WithValue(parent, ctxkey.ClientRequestID, strings.TrimSpace(value)+suffix)
	}
	if value, _ := parent.Value(ctxkey.RequestID).(string); strings.TrimSpace(value) != "" {
		parent = context.WithValue(parent, ctxkey.RequestID, strings.TrimSpace(value)+suffix)
	}
	return parent
}

func (h *OpenAIGatewayHandler) recordSplitOpenAIImageUsage(
	usageParent context.Context,
	c *gin.Context,
	apiKey *service.APIKey,
	subscription *service.UserSubscription,
	account *service.Account,
	result *service.OpenAIForwardResult,
	requestModel string,
	body []byte,
	channelMapping service.ChannelMappingResult,
	reqLog *zap.Logger,
) {
	if account == nil || result == nil {
		return
	}
	requestPayloadHash := service.HashUsageRequestPayload(body)
	inboundEndpoint := GetInboundEndpoint(c)
	upstreamEndpoint := GetUpstreamEndpoint(c, account.Platform)
	quotaPlatform := service.QuotaPlatform(c.Request.Context(), apiKey)
	upstreamModel := result.UpstreamModel
	channelUsageFields := clientRequestedUsageFields(c, channelMapping, requestModel, upstreamModel)
	userAgent := c.GetHeader("User-Agent")
	clientIP := ip.GetClientIP(c)
	sessionID := service.ExtractClientSessionID(c)
	h.submitMandatoryUsageRecordTask(usageParent, func(ctx context.Context) {
		if err := h.gatewayService.RecordUsage(ctx, &service.OpenAIRecordUsageInput{
			Result:             result,
			APIKey:             apiKey,
			User:               apiKey.User,
			Account:            account,
			Subscription:       subscription,
			InboundEndpoint:    inboundEndpoint,
			UpstreamEndpoint:   upstreamEndpoint,
			UserAgent:          userAgent,
			IPAddress:          clientIP,
			RequestPayloadHash: requestPayloadHash,
			APIKeyService:      h.apiKeyService,
			QuotaPlatform:      quotaPlatform,
			SessionID:          sessionID,
			ChannelUsageFields: channelUsageFields,
		}); err != nil {
			reqLog.Warn("openai.images.n_split_record_usage_failed", zap.Int64("account_id", account.ID), zap.Error(err))
		}
	})
}
