package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

const (
	VideoJobPending   = "pending"
	VideoJobRunning   = "running"
	VideoJobSettling  = "settling"
	VideoJobCompleted = "completed"
	VideoJobFailed    = "failed"
	VideoJobCanceled  = "canceled"
)

var (
	ErrVideoJobNotFound           = errors.New("video job not found")
	ErrVideoJobTransitionConflict = errors.New("video job transition conflict")
	ErrVideoJobCancelConflict     = errors.New("video job cancel conflict")
)

type VideoJob struct {
	ID              int64
	JobID           string
	UserID          int64
	APIKeyID        int64
	GroupID         int64
	AccountID       int64
	UpstreamJobID   int64
	Status          string
	RequestedModel  string
	UpstreamModel   string
	Prompt          string
	Resolution      string
	DurationSeconds int
	AspectRatio     string
	Audio           bool
	ImageSource     string
	ImageURL        string
	LocalInputName  string
	Result          json.RawMessage
	ErrorMessage    string
	HoldAmount      *float64
	ActualCost      *float64
	BillingSnapshot json.RawMessage
	RequestHash     string
	SettledAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
}

type VideoJobTransition struct {
	AccountID     *int64
	UpstreamJobID *int64
	Result        json.RawMessage
	ErrorMessage  *string
	HoldAmount    *float64
	ActualCost    *float64
	SettledAt     *time.Time
	StartedAt     *time.Time
	FinishedAt    *time.Time
}

type VideoJobRepository interface {
	CreateVideoJob(ctx context.Context, job *VideoJob) error
	GetVideoJob(ctx context.Context, jobID string) (*VideoJob, error)
	GetVideoJobForAPIKey(ctx context.Context, jobID string, apiKeyID int64) (*VideoJob, error)
	ListVideoJobsForAPIKey(ctx context.Context, apiKeyID int64, limit, offset int, status string) ([]*VideoJob, int, error)
	ListActiveVideoJobs(ctx context.Context, limit int) ([]*VideoJob, error)
	TransitionVideoJob(ctx context.Context, jobID string, allowedStatuses []string, status string, transition VideoJobTransition) error
}

func IsTerminalVideoJobStatus(status string) bool {
	switch status {
	case VideoJobCompleted, VideoJobFailed, VideoJobCanceled:
		return true
	default:
		return false
	}
}

func NewVideoJobID() (string, error) {
	var raw [24]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return "vidjob_" + base64.RawURLEncoding.EncodeToString(raw[:]), nil
}
