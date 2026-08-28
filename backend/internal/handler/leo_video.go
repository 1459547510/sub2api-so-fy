package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkghttputil "github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func leoVideoRequestPlatform(c *gin.Context) string {
	if c != nil && c.Request != nil {
		if platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context()); ok && service.IsMediaPlatform(platform) {
			return platform
		}
	}
	return service.PlatformLeo
}

func (h *OpenAIGatewayHandler) LeoVideoGeneration(c *gin.Context) {
	if service.PrefersLeoRespondAsync(c.Request.Header) {
		h.LeoVideoAsyncGeneration(c)
		return
	}
	streamStarted := false
	defer h.recoverResponsesPanic(c, &streamStarted)

	requestStart := time.Now()
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
	}

	body, err := pkghttputil.ReadRequestBodyWithPrealloc(c.Request)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.errorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
		}
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
	}
	if len(body) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
	}
	if contentType := strings.ToLower(strings.TrimSpace(c.GetHeader("Content-Type"))); contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Content-Type must be application/json")
		return
	}
	requestInfo, err := service.ValidateLeoVideoRequest(body)
	if err != nil {
		message := err.Error()
		if !json.Valid(body) {
			message = "Request body must be valid JSON"
		}
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", message)
		return
	}
	requestModel := strings.TrimSpace(requestInfo.Model)
	requestPlatform := leoVideoRequestPlatform(c)
	if requestModel == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
	}
	if strings.TrimSpace(requestInfo.Prompt) == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "prompt is required")
		return
	}
	if !service.GroupAllowsImageGeneration(apiKey.Group) {
		h.errorResponse(c, http.StatusForbidden, "permission_error", service.ImageGenerationPermissionMessage())
		return
	}

	reqLog := requestLogger(
		c,
		"handler.openai_gateway.leo_video",
		zap.Int64("user_id", subject.UserID),
		zap.Int64("api_key_id", apiKey.ID),
		zap.Any("group_id", apiKey.GroupID),
		zap.String("model", requestModel),
	)
	if !h.ensureResponsesDependencies(c, reqLog) {
		return
	}

	setOpsRequestContext(c, requestModel, false)
	setOpsEndpointContext(c, "/v1/videos/generations", int16(service.RequestTypeSync))
	if moderationBody := leoVideoModerationBody(requestInfo); len(moderationBody) > 0 {
		decision := h.checkSecurityAudit(c, reqLog, apiKey, subject, service.ContentModerationProtocolOpenAIImages, requestModel, moderationBody)
		if decision != nil && !decision.AllowNextStage {
			h.openAISecurityAuditError(c, decision)
			return
		}
	}
	imageReleaseFunc, acquired := h.acquireImageGenerationSlot(c, streamStarted)
	if !acquired {
		return
	}
	if imageReleaseFunc != nil {
		defer imageReleaseFunc()
	}

	if h.errorPassthroughService != nil {
		service.BindErrorPassthroughService(c, h.errorPassthroughService)
	}
	subscription, _ := middleware2.GetSubscriptionFromContext(c)
	service.SetOpsLatencyMs(c, service.OpsAuthLatencyMsKey, time.Since(requestStart).Milliseconds())

	userReleaseFunc, acquired := h.acquireResponsesUserSlot(c, subject.UserID, subject.Concurrency, false, &streamStarted, reqLog)
	if !acquired {
		return
	}
	if userReleaseFunc != nil {
		defer userReleaseFunc()
	}

	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription, service.QuotaPlatform(c.Request.Context(), apiKey)); err != nil {
		status, code, message, retryAfter := billingErrorDetails(err)
		if retryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
		}
		h.errorResponse(c, status, code, message)
		return
	}

	syncHold, err := h.holdSyncVideoJob(c, apiKey, subject, subscription, body)
	if err != nil {
		h.leoVideoCreateErrorResponse(c, err)
		return
	}
	holdSettledWithUsage := false
	defer func() {
		if !holdSettledWithUsage {
			h.releaseSyncVideoHold(syncHold)
		}
	}()

	requestCtx := c.Request.Context()
	sessionHash := h.gatewayService.GenerateExplicitSessionHash(c, body)
	failedAccountIDs := make(map[int64]struct{})
	var lastFailoverErr *service.UpstreamFailoverError
	var lastAccount *service.Account
	switchCount := 0
	maxAccountSwitches := h.maxAccountSwitches
	if maxAccountSwitches <= 0 {
		maxAccountSwitches = 3
	}
	routingStart := time.Now()

	for {
		selection, scheduleDecision, err := h.gatewayService.SelectAccountWithSchedulerForCapability(
			requestCtx,
			apiKey.GroupID,
			"",
			sessionHash,
			requestModel,
			failedAccountIDs,
			service.OpenAIUpstreamTransportHTTPSSE,
			"",
			false,
			false,
			false,
			requestPlatform,
		)
		if err != nil || selection == nil || selection.Account == nil {
			if len(failedAccountIDs) == 0 {
				cls := classifyNoAccountErrorFromGin(c, h.gatewayService, apiKey, requestModel, requestModel, requestPlatform)
				if !cls.ModelNotFound {
					markOpsRoutingCapacityLimitedIfNoAvailable(c, err)
				}
				h.errorResponse(c, cls.Status, cls.ErrType, cls.Message)
				return
			}
			if lastFailoverErr != nil {
				h.handleFailoverExhausted(c, lastFailoverErr, false, requestPlatform)
				h.persistSyncVideoJob(c, apiKey, subject, lastAccount, body, service.VideoJobFailed, nil, lastFailoverErr)
			} else {
				h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Upstream request failed")
				h.persistSyncVideoJob(c, apiKey, subject, lastAccount, body, service.VideoJobFailed, nil, errors.New("upstream request failed"))
			}
			return
		}

		reqLog.Debug("leo_video.account_schedule_decision",
			zap.String("layer", scheduleDecision.Layer),
			zap.Bool("sticky_session_hit", scheduleDecision.StickySessionHit),
			zap.Int("candidate_count", scheduleDecision.CandidateCount),
			zap.Int64("latency_ms", scheduleDecision.LatencyMs),
		)
		account := selection.Account
		lastAccount = account
		setOpsSelectedAccount(c, account.ID, account.Platform)
		accountReleaseFunc, slotResult := h.acquireResponsesAccountSlot(c, apiKey.GroupID, sessionHash, selection, false, &streamStarted, reqLog, false)
		if slotResult == openAISlotAcquireProfitVetoed {
			failedAccountIDs[account.ID] = struct{}{}
			continue
		}
		if slotResult != openAISlotAcquireOK {
			return
		}

		service.SetOpsLatencyMs(c, service.OpsRoutingLatencyMsKey, time.Since(routingStart).Milliseconds())
		forwardStart := time.Now()
		writerSizeBeforeForward := c.Writer.Size()
		result, err := func() (*service.OpenAIForwardResult, error) {
			defer func() {
				if accountReleaseFunc != nil {
					accountReleaseFunc()
				}
			}()
			return h.gatewayService.ForwardLeoVideo(requestCtx, c, account, body)
		}()
		forwardDurationMs := time.Since(forwardStart).Milliseconds()
		upstreamLatencyMs, _ := getContextInt64(c, service.OpsUpstreamLatencyMsKey)
		responseLatencyMs := forwardDurationMs
		if upstreamLatencyMs > 0 && forwardDurationMs > upstreamLatencyMs {
			responseLatencyMs = forwardDurationMs - upstreamLatencyMs
		}
		service.SetOpsLatencyMs(c, service.OpsResponseLatencyMsKey, responseLatencyMs)

		if err != nil {
			var validationErr *service.LeoVideoValidationError
			if errors.As(err, &validationErr) {
				h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", validationErr.Error())
				return
			}
			var rejectedErr *service.LeoVideoRejectedError
			if errors.As(err, &rejectedErr) {
				h.gatewayService.ReportOpenAIAccountScheduleResult(account, account.GetMappedModel(requestModel), false, nil)
				h.persistSyncVideoJob(c, apiKey, subject, account, body, service.VideoJobFailed, nil, rejectedErr)
				return
			}
			var failoverErr *service.UpstreamFailoverError
			if errors.As(err, &failoverErr) {
				h.gatewayService.ReportOpenAIAccountScheduleResult(account, account.GetMappedModel(requestModel), false, nil)
				if c.Writer.Size() != writerSizeBeforeForward {
					h.handleFailoverExhausted(c, failoverErr, true, requestPlatform)
					h.persistSyncVideoJob(c, apiKey, subject, account, body, service.VideoJobFailed, nil, failoverErr)
					return
				}
				h.gatewayService.RecordOpenAIAccountSwitch()
				failedAccountIDs[account.ID] = struct{}{}
				lastFailoverErr = failoverErr
				if switchCount >= maxAccountSwitches {
					h.handleFailoverExhausted(c, failoverErr, false, requestPlatform)
					h.persistSyncVideoJob(c, apiKey, subject, account, body, service.VideoJobFailed, nil, failoverErr)
					return
				}
				switchCount++
				reqLog.Warn("leo_video.upstream_failover_switching",
					zap.Int64("account_id", account.ID),
					zap.Int("upstream_status", failoverErr.StatusCode),
					zap.Int("switch_count", switchCount),
					zap.Int("max_switches", maxAccountSwitches),
				)
				continue
			}
			h.gatewayService.ReportOpenAIAccountScheduleResult(account, account.GetMappedModel(requestModel), false, nil)
			if c.Writer.Size() == writerSizeBeforeForward {
				h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Upstream request failed")
			}
			reqLog.Warn("leo_video.forward_failed", zap.Int64("account_id", account.ID), zap.Error(err))
			h.persistSyncVideoJob(c, apiKey, subject, account, body, service.VideoJobFailed, nil, err)
			return
		}

		h.gatewayService.ReportOpenAIAccountScheduleResult(account, account.GetMappedModel(requestModel), true, nil)
		if result == nil {
			return
		}
		holdSettledWithUsage = true
		recordLeoVideoUsage(c, h, reqLog, apiKey, subject, subscription, account, result, requestModel, body, func() {
			h.releaseSyncVideoHold(syncHold)
		})
		h.persistSyncVideoJob(c, apiKey, subject, account, body, service.VideoJobCompleted, result.VideoResult, nil)
		return
	}
}

func leoVideoModerationBody(info service.LeoVideoRequestInfo) []byte {
	payload := map[string]any{"prompt": info.Prompt}
	if imageURL := strings.TrimSpace(info.ImageURL); imageURL != "" {
		payload["images"] = []map[string]string{{"image_url": imageURL}}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return body
}

func (h *OpenAIGatewayHandler) persistSyncVideoJob(
	c *gin.Context,
	apiKey *service.APIKey,
	subject middleware2.AuthSubject,
	account *service.Account,
	body []byte,
	status string,
	result json.RawMessage,
	cause error,
) {
	if h == nil || h.videoJobService == nil || apiKey == nil {
		return
	}
	user := apiKey.User
	if user == nil {
		user = &service.User{ID: subject.UserID}
	}
	localInputName := ""
	if h.videoInputHandler != nil && h.videoInputHandler.store != nil {
		localInputName = strings.Join(h.videoInputHandler.store.TokensFromVideoRequest(body), ",")
	}
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = context.WithoutCancel(c.Request.Context())
	}
	_, err := h.videoJobService.RecordSync(ctx, service.RecordSyncVideoJobInput{
		APIKey: apiKey, User: user, Account: account, Body: body, LocalInputName: localInputName,
		Status: status, Result: result, ErrorMessage: syncVideoJobErrorMessage(cause),
		OutputStore: h.videoOutputStore,
	})
	if err == nil {
		return
	}
	logger.L().With(
		zap.String("component", "handler.openai_gateway.leo_video"),
		zap.Int64("user_id", subject.UserID),
		zap.Int64("api_key_id", apiKey.ID),
	).Warn("leo_video.persist_sync_job_failed", zap.Error(err))
}

func syncVideoJobErrorMessage(err error) string {
	var rejected *service.LeoVideoRejectedError
	if errors.As(err, &rejected) && rejected != nil && strings.TrimSpace(rejected.Message) != "" {
		return rejected.Message
	}
	var failoverErr *service.UpstreamFailoverError
	if errors.As(err, &failoverErr) && failoverErr != nil {
		message := service.PublicVideoErrorMessage(service.SanitizeVideoProviderMessage(service.ExtractUpstreamErrorMessage(failoverErr.ResponseBody)))
		if message != "" {
			return message
		}
	}
	if err != nil && strings.TrimSpace(err.Error()) != "" {
		return service.SanitizeVideoProviderMessage(service.PublicVideoErrorMessage(err.Error()))
	}
	return "Video service request failed"
}

func (h *OpenAIGatewayHandler) holdSyncVideoJob(
	c *gin.Context,
	apiKey *service.APIKey,
	subject middleware2.AuthSubject,
	subscription *service.UserSubscription,
	body []byte,
) (*service.VideoJob, error) {
	if h == nil || h.videoJobService == nil || h.videoJobService.Billing == nil {
		return nil, nil
	}
	user := apiKey.User
	if user == nil {
		user = &service.User{ID: subject.UserID}
	}
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	return h.videoJobService.HoldSync(ctx, service.CreateVideoJobInput{
		APIKey: apiKey, User: user, Subscription: subscription, Body: body,
	})
}

func (h *OpenAIGatewayHandler) releaseSyncVideoHold(job *service.VideoJob) {
	if h == nil || h.videoJobService == nil || h.videoJobService.Billing == nil || job == nil {
		return
	}
	if err := h.videoJobService.Billing.SettleWithoutCharge(context.Background(), job); err != nil {
		logger.L().With(
			zap.String("component", "handler.openai_gateway.leo_video"),
			zap.String("job_id", job.JobID),
		).Warn("leo_video.release_sync_hold_failed", zap.Error(err))
	}
}

func recordLeoVideoUsage(
	c *gin.Context,
	h *OpenAIGatewayHandler,
	reqLog *zap.Logger,
	apiKey *service.APIKey,
	subject middleware2.AuthSubject,
	subscription *service.UserSubscription,
	account *service.Account,
	result *service.OpenAIForwardResult,
	requestModel string,
	body []byte,
	after func(),
) {
	channelUsageFields := service.ChannelUsageFields{OriginalModel: requestModel, ChannelMappedModel: requestModel}
	if apiKey.GroupID != nil {
		mapping := h.gatewayService.ResolveChannelMapping(c.Request.Context(), *apiKey.GroupID, requestModel)
		channelUsageFields = mapping.ToUsageFields(requestModel, result.UpstreamModel)
	}
	inboundEndpoint := GetInboundEndpoint(c)
	userAgent := c.GetHeader("User-Agent")
	clientIP := ip.GetClientIP(c)
	requestPayloadHash := service.HashUsageRequestPayload(body)
	quotaPlatform := service.QuotaPlatform(c.Request.Context(), apiKey)
	h.submitOpenAIUsageRecordTask(c.Request.Context(), result, func(ctx context.Context) {
		defer func() {
			if after != nil {
				after()
			}
		}()
		if err := h.gatewayService.RecordUsage(ctx, &service.OpenAIRecordUsageInput{
			Result:             result,
			APIKey:             apiKey,
			User:               apiKey.User,
			Account:            account,
			Subscription:       subscription,
			InboundEndpoint:    inboundEndpoint,
			UpstreamEndpoint:   "/v1/videos/generations",
			UserAgent:          userAgent,
			IPAddress:          clientIP,
			RequestPayloadHash: requestPayloadHash,
			APIKeyService:      h.apiKeyService,
			QuotaPlatform:      quotaPlatform,
			ChannelUsageFields: channelUsageFields,
		}); err != nil {
			logger.L().With(
				zap.String("component", "handler.openai_gateway.leo_video"),
				zap.Int64("user_id", subject.UserID),
				zap.Int64("api_key_id", apiKey.ID),
				zap.Int64("account_id", account.ID),
			).Error("leo_video.record_usage_failed", zap.Error(err))
			reqLog.Debug("leo_video.record_usage_failed", zap.Error(err))
		}
	})
}
