package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type CreateVideoJobInput struct {
	APIKey         *APIKey
	User           *User
	Subscription   *UserSubscription
	Body           []byte
	LocalInputName string
}

// RecordSyncVideoJobInput writes a terminal workbench log for a synchronous
// generation. Billing stays on the sync usage path; this must not hold or settle.
type RecordSyncVideoJobInput struct {
	APIKey         *APIKey
	User           *User
	Account        *Account
	Body           []byte
	LocalInputName string
	Status         string
	Result         json.RawMessage
	ErrorMessage   string
	OutputStore    *VideoOutputStore
}

type VideoJobAccountSelector interface {
	Select(ctx context.Context, groupID int64, model string, excluded map[int64]struct{}) (*Account, error)
	GetByID(ctx context.Context, id int64) (*Account, error)
}

type VideoJobAsyncClient interface {
	CreateLeoAsyncVideo(ctx context.Context, account *Account, body []byte) (*LeoAsyncAccepted, error)
	GetLeoAsyncVideo(ctx context.Context, account *Account, upstreamJobID int64) (*LeoAsyncJob, error)
	CancelLeoAsyncVideo(ctx context.Context, account *Account, upstreamJobID int64) (*LeoAsyncJob, error)
}

type VideoJobService struct {
	Repo     VideoJobRepository
	Selector VideoJobAccountSelector
	Client   VideoJobAsyncClient
	Billing  *VideoJobBillingService
	Inputs   *VideoInputStore
}

func NewVideoJobService(repo VideoJobRepository, selector VideoJobAccountSelector, client VideoJobAsyncClient, billing *VideoJobBillingService) *VideoJobService {
	return &VideoJobService{Repo: repo, Selector: selector, Client: client, Billing: billing}
}

func (s *VideoJobService) SetVideoInputStore(inputs *VideoInputStore) {
	if s != nil {
		s.Inputs = inputs
	}
}

func (s *VideoJobService) Create(ctx context.Context, in CreateVideoJobInput) (*VideoJob, error) {
	if s == nil || s.Repo == nil || s.Selector == nil || s.Client == nil || s.Billing == nil {
		return nil, errors.New("video job service is not configured")
	}
	if in.APIKey == nil || in.User == nil || in.APIKey.Group == nil {
		return nil, errors.New("video job request context is incomplete")
	}
	info, err := parseVideoJobCreateBody(in.Body)
	if err != nil {
		return nil, err
	}
	groupID := in.APIKey.Group.ID
	if in.APIKey.GroupID != nil {
		groupID = *in.APIKey.GroupID
	}
	account, err := s.Selector.Select(ctx, groupID, info.Model, nil)
	if err != nil || account == nil {
		if err == nil {
			err = errors.New("no leo account is available")
		}
		return nil, err
	}
	info, err = ValidateLeoVideoRequestForModel(in.Body, accountMappedModel(account, info.Model))
	if err != nil {
		return nil, err
	}
	jobID, err := NewVideoJobID()
	if err != nil {
		return nil, err
	}
	imageSource := "none"
	if strings.TrimSpace(in.LocalInputName) != "" {
		imageSource = "local"
	} else if info.ImageURL != "" {
		imageSource = "url"
	}
	job := &VideoJob{
		JobID: jobID, UserID: in.User.ID, APIKeyID: in.APIKey.ID, GroupID: groupID, AccountID: account.ID,
		Status: VideoJobPending, RequestedModel: info.Model, UpstreamModel: accountMappedModel(account, info.Model),
		Prompt: info.Prompt, Resolution: info.Resolution, DurationSeconds: info.DurationSeconds,
		AspectRatio: info.AspectRatio, Audio: info.Audio, ImageSource: imageSource, ImageURL: info.ImageURL,
		LocalInputName: strings.TrimSpace(in.LocalInputName), RequestHash: HashUsageRequestPayload(in.Body),
	}
	if err := s.Billing.Prepare(ctx, job, in.APIKey, in.User, in.Subscription); err != nil {
		return nil, err
	}
	if err := s.Repo.CreateVideoJob(ctx, job); err != nil {
		_ = s.Billing.SettleWithoutCharge(ctx, job)
		return nil, err
	}

	excluded := map[int64]struct{}{}
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			excluded[account.ID] = struct{}{}
			next, selectErr := s.Selector.Select(ctx, groupID, info.Model, excluded)
			if selectErr != nil || next == nil {
				if selectErr != nil {
					err = selectErr
				} else {
					err = errors.New("no leo failover account is available")
				}
				break
			}
			account = next
			job.AccountID = account.ID
			accountID := account.ID
			if err := s.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending}, VideoJobPending, VideoJobTransition{AccountID: &accountID}); err != nil {
				return nil, err
			}
		}
		accepted, createErr := s.Client.CreateLeoAsyncVideo(ctx, account, in.Body)
		if createErr == nil && accepted != nil {
			upstreamID := accepted.JobID
			status := accepted.Status
			if status != VideoJobRunning {
				status = VideoJobPending
			}
			transition := VideoJobTransition{UpstreamJobID: &upstreamID}
			if status == VideoJobRunning {
				now := time.Now()
				transition.StartedAt = &now
			}
			if err := s.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending}, status, transition); err != nil {
				return nil, err
			}
			job.UpstreamJobID = upstreamID
			job.UpstreamModel = accountMappedModel(account, info.Model)
			job.Status = status
			if transition.StartedAt != nil {
				job.StartedAt = transition.StartedAt
			}
			return job, nil
		}
		err = createErr
		var upstreamErr *LeoAsyncUpstreamError
		if !errors.As(createErr, &upstreamErr) || !upstreamErr.Retryable || upstreamErr.Ambiguous {
			break
		}
	}

	message := videoJobFailureMessage(err)
	if releaseErr := s.Billing.SettleWithoutCharge(ctx, job); releaseErr != nil && err == nil {
		err = releaseErr
	}
	finished := time.Now()
	transition := VideoJobTransition{ErrorMessage: &message, FinishedAt: &finished}
	if job.SettledAt != nil {
		transition.SettledAt = job.SettledAt
	}
	_ = s.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending}, VideoJobFailed, transition)
	if err == nil {
		err = errors.New(message)
	}
	return nil, err
}

// HoldSync reserves the estimated video cost for a synchronous generation.
// It does not persist a job or call upstream. The caller must release the
// hold after usage settlement or on any failure before forwarding completes.
func (s *VideoJobService) HoldSync(ctx context.Context, in CreateVideoJobInput) (*VideoJob, error) {
	if s == nil || s.Billing == nil {
		return nil, errors.New("video job service is not configured")
	}
	if in.APIKey == nil || in.User == nil || in.APIKey.Group == nil {
		return nil, errors.New("video job request context is incomplete")
	}
	info, err := parseVideoJobCreateBody(in.Body)
	if err != nil {
		return nil, err
	}
	groupID := in.APIKey.Group.ID
	if in.APIKey.GroupID != nil {
		groupID = *in.APIKey.GroupID
	}
	jobID, err := NewVideoJobID()
	if err != nil {
		return nil, err
	}
	job := &VideoJob{
		JobID: jobID, UserID: in.User.ID, APIKeyID: in.APIKey.ID, GroupID: groupID,
		Status: VideoJobPending, RequestedModel: info.Model, UpstreamModel: info.Model,
		Prompt: info.Prompt, Resolution: info.Resolution, DurationSeconds: info.DurationSeconds,
		AspectRatio: info.AspectRatio, Audio: info.Audio, RequestHash: HashUsageRequestPayload(in.Body),
	}
	if err := s.Billing.Prepare(ctx, job, in.APIKey, in.User, in.Subscription); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *VideoJobService) RecordSync(ctx context.Context, in RecordSyncVideoJobInput) (*VideoJob, error) {
	if s == nil || s.Repo == nil {
		return nil, errors.New("video job service is not configured")
	}
	if in.APIKey == nil || in.User == nil {
		return nil, errors.New("video job request context is incomplete")
	}
	status := strings.TrimSpace(in.Status)
	if status != VideoJobCompleted && status != VideoJobFailed {
		return nil, errors.New("sync video job status is invalid")
	}
	info, err := parseVideoJobCreateBody(in.Body)
	if err != nil {
		return nil, err
	}
	groupID := int64(0)
	if in.APIKey.Group != nil {
		groupID = in.APIKey.Group.ID
	}
	if in.APIKey.GroupID != nil {
		groupID = *in.APIKey.GroupID
	}
	accountID := int64(0)
	upstreamModel := info.Model
	if in.Account != nil {
		accountID = in.Account.ID
		upstreamModel = accountMappedModel(in.Account, info.Model)
	}
	jobID, err := NewVideoJobID()
	if err != nil {
		return nil, err
	}
	imageSource := "none"
	if strings.TrimSpace(in.LocalInputName) != "" {
		imageSource = "local"
	} else if info.ImageURL != "" {
		imageSource = "url"
	}
	now := time.Now()
	result := json.RawMessage(nil)
	if status == VideoJobCompleted && len(in.Result) != 0 {
		if public := PublicLeoVideoResult(in.Result); len(public) != 0 {
			result = public
		} else {
			result = append(json.RawMessage(nil), in.Result...)
		}
	}
	errorMessage := ""
	if status == VideoJobFailed {
		errorMessage = SanitizeVideoProviderMessage(PublicVideoErrorMessage(strings.TrimSpace(in.ErrorMessage)))
		if errorMessage == "" {
			errorMessage = "Video service request failed"
		}
	}
	job := &VideoJob{
		JobID: jobID, UserID: in.User.ID, APIKeyID: in.APIKey.ID, GroupID: groupID, AccountID: accountID,
		Status: status, RequestedModel: info.Model, UpstreamModel: upstreamModel,
		Prompt: info.Prompt, Resolution: info.Resolution, DurationSeconds: info.DurationSeconds,
		AspectRatio: info.AspectRatio, Audio: info.Audio, ImageSource: imageSource, ImageURL: info.ImageURL,
		LocalInputName: strings.TrimSpace(in.LocalInputName), Result: result, ErrorMessage: errorMessage,
		RequestHash: HashUsageRequestPayload(in.Body), SettledAt: &now, StartedAt: &now, FinishedAt: &now,
	}
	if err := s.Repo.CreateVideoJob(ctx, job); err != nil {
		return nil, err
	}
	_ = MarkVideoInputTerminal(s.Inputs, job.LocalInputName, now)
	if status == VideoJobCompleted && in.OutputStore != nil && len(in.Result) != 0 {
		raw := append(json.RawMessage(nil), in.Result...)
		go s.saveSyncVideoOutput(job.JobID, raw, in.OutputStore)
	}
	return job, nil
}

func (s *VideoJobService) saveSyncVideoOutput(jobID string, raw json.RawMessage, store *VideoOutputStore) {
	if s == nil || s.Repo == nil || store == nil || strings.TrimSpace(jobID) == "" || len(raw) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	saved, err := store.Save(ctx, jobID, raw)
	if err != nil || len(saved) == 0 {
		return
	}
	_ = s.Repo.TransitionVideoJob(ctx, jobID, []string{VideoJobCompleted}, VideoJobCompleted, VideoJobTransition{Result: saved})
}

func (s *VideoJobService) List(ctx context.Context, apiKeyID int64, limit, offset int, status string) ([]*VideoJob, int, error) {
	if s == nil || s.Repo == nil {
		return nil, 0, errors.New("video job service is not configured")
	}
	return s.Repo.ListVideoJobsForAPIKey(ctx, apiKeyID, limit, offset, status)
}

func (s *VideoJobService) Get(ctx context.Context, jobID string, apiKeyID int64) (*VideoJob, error) {
	if s == nil || s.Repo == nil {
		return nil, errors.New("video job service is not configured")
	}
	return s.Repo.GetVideoJobForAPIKey(ctx, jobID, apiKeyID)
}

func (s *VideoJobService) Cancel(ctx context.Context, jobID string, apiKeyID int64) (*VideoJob, error) {
	if s == nil || s.Repo == nil || s.Selector == nil || s.Client == nil || s.Billing == nil {
		return nil, errors.New("video job service is not configured")
	}
	job, err := s.Repo.GetVideoJobForAPIKey(ctx, jobID, apiKeyID)
	if err != nil {
		return nil, err
	}
	if job.Status != VideoJobPending || job.UpstreamJobID <= 0 {
		return nil, ErrVideoJobCancelConflict
	}
	account, err := s.Selector.GetByID(ctx, job.AccountID)
	if err != nil {
		return nil, err
	}
	canceled, err := s.Client.CancelLeoAsyncVideo(ctx, account, job.UpstreamJobID)
	if err != nil {
		return nil, err
	}
	if canceled == nil || canceled.Status != VideoJobCanceled {
		return nil, ErrVideoJobCancelConflict
	}
	if err := s.Billing.SettleWithoutCharge(ctx, job); err != nil {
		return nil, err
	}
	finished := time.Now()
	transition := VideoJobTransition{FinishedAt: &finished, SettledAt: job.SettledAt}
	if err := s.Repo.TransitionVideoJob(ctx, job.JobID, []string{VideoJobPending}, VideoJobCanceled, transition); err != nil {
		return nil, err
	}
	job.Status = VideoJobCanceled
	job.FinishedAt = &finished
	if err := MarkVideoInputTerminal(s.Inputs, job.LocalInputName, finished); err != nil {
		return nil, err
	}
	return job, nil
}

func parseVideoJobCreateBody(body []byte) (LeoVideoRequestInfo, error) {
	if !gjson.ValidBytes(body) {
		return LeoVideoRequestInfo{}, errors.New("request body must be valid JSON")
	}
	info, err := ValidateLeoVideoRequest(body)
	if err != nil {
		return LeoVideoRequestInfo{}, err
	}
	return info, nil
}

func accountMappedModel(account *Account, requested string) string {
	if account == nil {
		return requested
	}
	upstream, _ := account.ResolveMappedModel(requested)
	if strings.TrimSpace(upstream) == "" {
		return requested
	}
	return upstream
}

func videoJobFailureMessage(err error) string {
	var upstreamErr *LeoAsyncUpstreamError
	if errors.As(err, &upstreamErr) && strings.TrimSpace(upstreamErr.Message) != "" {
		return SanitizeVideoProviderMessage(upstreamErr.Message)
	}
	if err != nil && strings.TrimSpace(err.Error()) != "" {
		return SanitizeVideoProviderMessage(err.Error())
	}
	return "Video provider request failed"
}

// OpenAIGatewayService implements account selection for the durable video service.
func (s *OpenAIGatewayService) Select(ctx context.Context, groupID int64, model string, excluded map[int64]struct{}) (*Account, error) {
	if s == nil {
		return nil, errors.New("gateway service is nil")
	}
	platform := PlatformLeo
	if resolved, ok := ResolvedTargetPlatformFromContext(ctx); ok && IsMediaPlatform(resolved) {
		platform = resolved
	}
	selection, _, err := s.SelectAccountWithSchedulerForCapability(ctx, int64Pointer(groupID), "", "", model, excluded, OpenAIUpstreamTransportHTTPSSE, "", false, false, false, platform)
	if err != nil || selection == nil || selection.Account == nil {
		if err == nil {
			err = errors.New("no leo account is available")
		}
		return nil, err
	}
	return selection.Account, nil
}

func (s *OpenAIGatewayService) GetByID(ctx context.Context, id int64) (*Account, error) {
	if s == nil || s.accountRepo == nil {
		return nil, ErrAccountNotFound
	}
	return s.accountRepo.GetByID(ctx, id)
}
