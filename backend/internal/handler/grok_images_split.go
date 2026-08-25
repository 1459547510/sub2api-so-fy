package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type splitGrokImageTaskResult struct {
	index   int
	account *service.Account
	result  *service.OpenAIForwardResult
	body    []byte
	header  http.Header
	status  int
	err     error
}

func (h *OpenAIGatewayHandler) handleGrokImagesNSplit(
	c *gin.Context,
	apiKey *service.APIKey,
	subject middleware.AuthSubject,
	subscription *service.UserSubscription,
	endpoint service.GrokMediaEndpoint,
	requestModel string,
	routingModel string,
	requestInfo service.GrokMediaRequestInfo,
	body []byte,
	contentType string,
	requestCtx context.Context,
	sessionHash string,
	reqLog *zap.Logger,
) bool {
	if requestInfo.N <= 1 {
		return false
	}

	tasks := make([]splitGrokImageTaskResult, requestInfo.N)
	var wg sync.WaitGroup
	for index := range tasks {
		tasks[index].index = index
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			result := h.runGrokImagesNSplitTask(
				c,
				apiKey.GroupID,
				endpoint,
				requestModel,
				routingModel,
				body,
				contentType,
				requestCtx,
				splitImageSessionHash(sessionHash, index),
				reqLog,
			)
			result.index = index
			tasks[index] = result
		}(index)
	}
	wg.Wait()

	parts := make([][]byte, 0, len(tasks))
	successes := make([]splitGrokImageTaskResult, 0, len(tasks))
	for _, task := range tasks {
		if task.result == nil || task.result.ImageCount <= 0 || len(task.body) == 0 {
			continue
		}
		parts = append(parts, task.body)
		successes = append(successes, task)
	}
	if len(successes) == 0 {
		for _, task := range tasks {
			if task.err != nil {
				if failover, ok := task.err.(*service.UpstreamFailoverError); ok {
					h.handleFailoverExhausted(c, failover, false, service.PlatformGrok)
					return true
				}
			}
		}
		h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Upstream request failed")
		return true
	}

	merged, err := service.MergeOpenAIImageResponses(parts)
	if err != nil {
		reqLog.Warn("grok_media.n_split_response_merge_failed", zap.Error(err))
		h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Upstream response could not be combined")
		return true
	}
	if len(successes) < len(tasks) {
		reqLog.Warn("grok_media.n_split_partial_success",
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
		recordGrokMediaUsage(splitImageUsageContext(c.Request.Context(), task.index), c, h, reqLog, apiKey, subject, subscription, task.account, task.result, requestModel, body, "")
	}
	return true
}

func (h *OpenAIGatewayHandler) runGrokImagesNSplitTask(
	parent *gin.Context,
	groupID *int64,
	endpoint service.GrokMediaEndpoint,
	requestModel string,
	routingModel string,
	body []byte,
	contentType string,
	requestCtx context.Context,
	sessionHash string,
	reqLog *zap.Logger,
) splitGrokImageTaskResult {
	child, recorder := newBufferedGatewayContext(parent)
	childBody, childContentType, err := service.RewriteOpenAIImagesN(body, contentType, 1)
	if err != nil {
		return splitGrokImageTaskResult{err: err}
	}
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
			return splitGrokImageTaskResult{err: err}
		}
		selection, _, selectErr := h.gatewayService.SelectAccountWithSchedulerForCapability(
			requestCtx,
			groupID,
			"",
			sessionHash,
			routingModel,
			failedAccountIDs,
			service.OpenAIUpstreamTransportHTTPSSE,
			service.OpenAIEndpointCapabilityGrokMediaGeneration,
			false,
			false,
			false,
			service.PlatformGrok,
		)
		if selectErr != nil || selection == nil || selection.Account == nil {
			if lastFailoverErr != nil {
				return splitGrokImageTaskResult{err: lastFailoverErr}
			}
			if selectErr == nil {
				selectErr = service.ErrNoAvailableAccounts
			}
			return splitGrokImageTaskResult{err: selectErr}
		}

		account := selection.Account
		eligible, _, eligibilityErr := h.ensureGrokMediaAccountEligibility(requestCtx, account)
		if !eligible {
			failedAccountIDs[account.ID] = struct{}{}
			if eligibilityErr != nil {
				reqLog.Warn("grok_media.n_split_account_eligibility_probe_failed", zap.Int64("account_id", account.ID), zap.Error(eligibilityErr))
			}
			if switchCount >= maxAccountSwitches {
				return splitGrokImageTaskResult{err: fmt.Errorf("no eligible media account")}
			}
			switchCount++
			continue
		}
		sessionHash = h.bindSplitImagePoolSession(requestCtx, groupID, sessionHash, account, reqLog)
		childStreamStarted := false
		release, slotResult := h.acquireResponsesAccountSlot(child, groupID, sessionHash, selection, false, &childStreamStarted, reqLog)
		if slotResult == openAISlotAcquireProfitVetoed {
			failedAccountIDs[account.ID] = struct{}{}
			continue
		}
		if slotResult != openAISlotAcquireOK {
			return splitGrokImageTaskResult{err: fmt.Errorf("media account slot unavailable")}
		}

		result, forwardErr := func() (*service.OpenAIForwardResult, error) {
			defer func() {
				if release != nil {
					release()
				}
			}()
			return h.gatewayService.ForwardGrokMedia(requestCtx, child, account, endpoint, "", childBody, childContentType)
		}()
		if result != nil && result.ImageCount > 0 && recorder.Body.Len() > 0 {
			h.gatewayService.ReportOpenAIAccountScheduleResult(account, grokMediaScheduleModel(account, routingModel, result), true, nil)
			return splitGrokImageTaskResult{
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
			h.gatewayService.ReportOpenAIAccountScheduleResult(account, grokMediaScheduleModel(account, routingModel, nil), false, nil)
			return splitGrokImageTaskResult{err: forwardErr}
		}
		h.gatewayService.ReportOpenAIAccountScheduleResult(account, grokMediaScheduleModel(account, routingModel, nil), false, nil)
		if failoverErr.RetryableOnSameAccount {
			retryLimit := account.GetPoolModeRetryCount()
			if sameAccountRetryCount[account.ID] < retryLimit {
				sameAccountRetryCount[account.ID]++
				select {
				case <-requestCtx.Done():
					return splitGrokImageTaskResult{err: requestCtx.Err()}
				case <-time.After(sameAccountRetryDelay):
				}
				continue
			}
		}
		h.gatewayService.RecordOpenAIAccountSwitch()
		failedAccountIDs[account.ID] = struct{}{}
		lastFailoverErr = failoverErr
		if switchCount >= maxAccountSwitches {
			return splitGrokImageTaskResult{err: lastFailoverErr}
		}
		switchCount++
		if h.gatewayService.ShouldStopOpenAIOAuth429Failover(account, failoverErr.StatusCode, switchCount, &oauth429FailoverState) {
			return splitGrokImageTaskResult{err: lastFailoverErr}
		}
	}
}
