package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

type VideoJobRuntime struct {
	Repo         VideoJobRepository
	Accounts     VideoJobAccountSelector
	Client       VideoJobAsyncClient
	Billing      *VideoJobBillingService
	PollInterval time.Duration

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func (r *VideoJobRuntime) Start(ctx context.Context) {
	if r == nil {
		return
	}
	r.mu.Lock()
	if r.done != nil {
		r.mu.Unlock()
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	r.done = make(chan struct{})
	done := r.done
	interval := r.PollInterval
	if interval <= 0 {
		interval = 2 * time.Second
	}
	r.mu.Unlock()

	go func() {
		defer close(done)
		_ = r.RunOnce(runCtx)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-runCtx.Done():
				return
			case <-ticker.C:
				_ = r.RunOnce(runCtx)
			}
		}
	}()
}

func (r *VideoJobRuntime) Stop() {
	if r == nil {
		return
	}
	r.mu.Lock()
	cancel := r.cancel
	done := r.done
	r.mu.Unlock()
	if cancel == nil || done == nil {
		return
	}
	cancel()
	<-done
	r.mu.Lock()
	r.cancel = nil
	r.done = nil
	r.mu.Unlock()
}

func (r *VideoJobRuntime) RunOnce(ctx context.Context) error {
	if r == nil || r.Repo == nil || r.Accounts == nil || r.Client == nil || r.Billing == nil {
		return errors.New("video job runtime is not configured")
	}
	jobs, err := r.Repo.ListActiveVideoJobs(ctx, 50)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if err := r.processJob(ctx, job); err != nil && ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

func (r *VideoJobRuntime) processJob(ctx context.Context, job *VideoJob) error {
	if job == nil || IsTerminalVideoJobStatus(job.Status) {
		return nil
	}
	if job.Status == VideoJobSettling {
		if len(job.Result) == 0 {
			return r.failJob(ctx, job, errors.New("settling video job has no result"))
		}
		return r.settleCompleted(ctx, job, job.Result)
	}
	if job.UpstreamJobID <= 0 {
		return r.failJob(ctx, job, errors.New("video job has no upstream job ID"))
	}
	account, err := r.Accounts.GetByID(ctx, job.AccountID)
	if err != nil {
		return err
	}
	upstream, err := r.Client.GetLeoAsyncVideo(ctx, account, job.UpstreamJobID)
	if err != nil {
		var upstreamErr *LeoAsyncUpstreamError
		if errors.As(err, &upstreamErr) && upstreamErr.Retryable {
			return nil
		}
		return r.failJob(ctx, job, err)
	}
	if upstream == nil {
		return nil
	}
	switch strings.ToLower(strings.TrimSpace(upstream.Status)) {
	case VideoJobPending:
		return nil
	case VideoJobRunning:
		if job.Status == VideoJobPending {
			return r.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending}, VideoJobRunning, VideoJobTransition{})
		}
		return nil
	case VideoJobCompleted:
		result := append(json.RawMessage(nil), upstream.Result...)
		if err := r.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending, VideoJobRunning}, VideoJobSettling, VideoJobTransition{Result: result}); err != nil {
			return err
		}
		job.Status = VideoJobSettling
		job.Result = result
		return r.settleCompleted(ctx, job, result)
	case VideoJobFailed, VideoJobCanceled:
		message := ""
		if upstream.Error != nil {
			message = strings.TrimSpace(upstream.Error.Message)
		}
		if message == "" {
			message = "LeoStudio video job failed"
		}
		if apiKey := account.GetLeoAPIKey(); apiKey != "" {
			message = strings.ReplaceAll(message, apiKey, "***")
		}
		if err := r.Billing.SettleWithoutCharge(ctx, job); err != nil {
			return err
		}
		finished := time.Now()
		status := VideoJobFailed
		if strings.EqualFold(upstream.Status, VideoJobCanceled) {
			status = VideoJobCanceled
		}
		return r.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending, VideoJobRunning}, status, VideoJobTransition{ErrorMessage: &message, FinishedAt: &finished, SettledAt: job.SettledAt})
	default:
		return nil
	}
}

func (r *VideoJobRuntime) settleCompleted(ctx context.Context, job *VideoJob, result json.RawMessage) error {
	if err := r.Billing.SettleCompleted(ctx, job, result); err != nil {
		return nil
	}
	finished := time.Now()
	transition := VideoJobTransition{Result: result, ActualCost: job.ActualCost, FinishedAt: &finished, SettledAt: job.SettledAt}
	if err := r.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobSettling}, VideoJobCompleted, transition); err != nil {
		return err
	}
	return nil
}

func (r *VideoJobRuntime) failJob(ctx context.Context, job *VideoJob, err error) error {
	if releaseErr := r.Billing.SettleWithoutCharge(ctx, job); releaseErr != nil {
		return releaseErr
	}
	message := videoJobFailureMessage(err)
	finished := time.Now()
	return r.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending, VideoJobRunning, VideoJobSettling}, VideoJobFailed, VideoJobTransition{ErrorMessage: &message, FinishedAt: &finished, SettledAt: job.SettledAt})
}
