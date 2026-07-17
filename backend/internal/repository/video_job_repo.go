package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/videojob"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type videoJobRepository struct {
	client *dbent.Client
}

func NewVideoJobRepository(client *dbent.Client) service.VideoJobRepository {
	return &videoJobRepository{client: client}
}

func (r *videoJobRepository) CreateVideoJob(ctx context.Context, job *service.VideoJob) error {
	client := clientFromContext(ctx, r.client)
	builder := client.VideoJob.Create().
		SetJobID(job.JobID).
		SetUserID(job.UserID).
		SetAPIKeyID(job.APIKeyID).
		SetGroupID(job.GroupID).
		SetAccountID(job.AccountID).
		SetStatus(job.Status).
		SetRequestedModel(job.RequestedModel).
		SetUpstreamModel(job.UpstreamModel).
		SetPrompt(job.Prompt).
		SetResolution(job.Resolution).
		SetDurationSeconds(job.DurationSeconds).
		SetAspectRatio(job.AspectRatio).
		SetAudio(job.Audio).
		SetImageSource(job.ImageSource).
		SetRequestHash(job.RequestHash)
	if job.UpstreamJobID != 0 {
		builder.SetUpstreamJobID(job.UpstreamJobID)
	}
	if job.ImageURL != "" {
		builder.SetImageURL(job.ImageURL)
	}
	if job.LocalInputName != "" {
		builder.SetLocalInputName(job.LocalInputName)
	}
	if len(job.Result) != 0 {
		builder.SetResult(job.Result)
	}
	if job.ErrorMessage != "" {
		builder.SetErrorMessage(job.ErrorMessage)
	}
	if job.HoldAmount != nil {
		builder.SetHoldAmount(*job.HoldAmount)
	}
	if job.ActualCost != nil {
		builder.SetActualCost(*job.ActualCost)
	}
	if len(job.BillingSnapshot) != 0 {
		builder.SetBillingSnapshot(job.BillingSnapshot)
	}
	if job.SettledAt != nil {
		builder.SetSettledAt(*job.SettledAt)
	}
	if job.StartedAt != nil {
		builder.SetStartedAt(*job.StartedAt)
	}
	if job.FinishedAt != nil {
		builder.SetFinishedAt(*job.FinishedAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	applyVideoJobEntity(job, created)
	return nil
}

func (r *videoJobRepository) GetVideoJob(ctx context.Context, jobID string) (*service.VideoJob, error) {
	row, err := r.client.VideoJob.Query().Where(videojob.JobIDEQ(jobID)).Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrVideoJobNotFound
		}
		return nil, err
	}
	return videoJobEntityToService(row), nil
}

func (r *videoJobRepository) GetVideoJobForAPIKey(ctx context.Context, jobID string, apiKeyID int64) (*service.VideoJob, error) {
	row, err := r.client.VideoJob.Query().
		Where(videojob.JobIDEQ(jobID), videojob.APIKeyIDEQ(apiKeyID)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrVideoJobNotFound
		}
		return nil, err
	}
	return videoJobEntityToService(row), nil
}

func (r *videoJobRepository) ListVideoJobsForAPIKey(ctx context.Context, apiKeyID int64, limit int, status string) ([]*service.VideoJob, error) {
	limit = normalizeVideoJobLimit(limit)
	query := r.client.VideoJob.Query().Where(videojob.APIKeyIDEQ(apiKeyID))
	if status != "" {
		query.Where(videojob.StatusEQ(status))
	}
	rows, err := query.
		Order(dbent.Desc(videojob.FieldCreatedAt), dbent.Desc(videojob.FieldID)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return videoJobEntitiesToService(rows), nil
}

func (r *videoJobRepository) ListActiveVideoJobs(ctx context.Context, limit int) ([]*service.VideoJob, error) {
	rows, err := r.client.VideoJob.Query().
		Where(videojob.StatusIn(service.VideoJobPending, service.VideoJobRunning, service.VideoJobSettling)).
		Order(dbent.Asc(videojob.FieldUpdatedAt), dbent.Asc(videojob.FieldID)).
		Limit(normalizeVideoJobLimit(limit)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return videoJobEntitiesToService(rows), nil
}

func (r *videoJobRepository) TransitionVideoJob(ctx context.Context, jobID string, allowedStatuses []string, status string, transition service.VideoJobTransition) error {
	if len(allowedStatuses) == 0 {
		return service.ErrVideoJobTransitionConflict
	}
	client := clientFromContext(ctx, r.client)
	update := client.VideoJob.Update().
		Where(videojob.JobIDEQ(jobID), videojob.StatusIn(allowedStatuses...)).
		SetStatus(status)
	if transition.AccountID != nil {
		update.SetAccountID(*transition.AccountID)
	}
	if transition.UpstreamJobID != nil {
		update.SetUpstreamJobID(*transition.UpstreamJobID)
	}
	if len(transition.Result) != 0 {
		update.SetResult(transition.Result)
	}
	if transition.ErrorMessage != nil {
		update.SetErrorMessage(*transition.ErrorMessage)
	}
	if transition.HoldAmount != nil {
		update.SetHoldAmount(*transition.HoldAmount)
	}
	if transition.ActualCost != nil {
		update.SetActualCost(*transition.ActualCost)
	}
	if transition.SettledAt != nil {
		update.SetSettledAt(*transition.SettledAt)
	}
	if transition.StartedAt != nil {
		update.SetStartedAt(*transition.StartedAt)
	}
	if transition.FinishedAt != nil {
		update.SetFinishedAt(*transition.FinishedAt)
	}

	affected, err := update.Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return service.ErrVideoJobTransitionConflict
	}
	return nil
}

func normalizeVideoJobLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func applyVideoJobEntity(dst *service.VideoJob, src *dbent.VideoJob) {
	if dst == nil || src == nil {
		return
	}
	*dst = *videoJobEntityToService(src)
}

func videoJobEntityToService(row *dbent.VideoJob) *service.VideoJob {
	if row == nil {
		return nil
	}
	job := &service.VideoJob{
		ID: row.ID, JobID: row.JobID, UserID: row.UserID, APIKeyID: row.APIKeyID,
		GroupID: row.GroupID, AccountID: row.AccountID, Status: row.Status,
		RequestedModel: row.RequestedModel, UpstreamModel: row.UpstreamModel,
		Prompt: row.Prompt, Resolution: row.Resolution, DurationSeconds: row.DurationSeconds,
		AspectRatio: row.AspectRatio, Audio: row.Audio, ImageSource: row.ImageSource,
		Result: row.Result, HoldAmount: row.HoldAmount, ActualCost: row.ActualCost,
		BillingSnapshot: row.BillingSnapshot, RequestHash: row.RequestHash,
		SettledAt: row.SettledAt, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		StartedAt: row.StartedAt, FinishedAt: row.FinishedAt,
	}
	if row.UpstreamJobID != nil {
		job.UpstreamJobID = *row.UpstreamJobID
	}
	if row.ImageURL != nil {
		job.ImageURL = *row.ImageURL
	}
	if row.LocalInputName != nil {
		job.LocalInputName = *row.LocalInputName
	}
	if row.ErrorMessage != nil {
		job.ErrorMessage = *row.ErrorMessage
	}
	return job
}

func videoJobEntitiesToService(rows []*dbent.VideoJob) []*service.VideoJob {
	jobs := make([]*service.VideoJob, 0, len(rows))
	for _, row := range rows {
		jobs = append(jobs, videoJobEntityToService(row))
	}
	return jobs
}
