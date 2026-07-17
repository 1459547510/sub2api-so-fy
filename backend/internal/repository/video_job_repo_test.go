package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newVideoJobRepositoryTestRepo(t *testing.T) service.VideoJobRepository {
	t.Helper()
	dsn := fmt.Sprintf("file:video_job_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	driver := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(driver)))
	t.Cleanup(func() { _ = client.Close() })
	return NewVideoJobRepository(client)
}

func TestVideoJobRepositoryCreateListAndTransition(t *testing.T) {
	ctx := context.Background()
	repo := newVideoJobRepositoryTestRepo(t)

	first := &service.VideoJob{
		JobID:           "vidjob_test_1",
		UserID:          1,
		APIKeyID:        2,
		GroupID:         3,
		AccountID:       9,
		UpstreamJobID:   42,
		Status:          service.VideoJobPending,
		RequestedModel:  "seedance-2.0",
		UpstreamModel:   "seedance-2.0",
		Prompt:          "waves",
		Resolution:      "720p",
		DurationSeconds: 8,
		AspectRatio:     "16:9",
		BillingSnapshot: []byte(`{"video_720p_per_second":0.1}`),
		RequestHash:     "request-hash-1",
	}
	require.NoError(t, repo.CreateVideoJob(ctx, first))

	got, err := repo.GetVideoJobForAPIKey(ctx, first.JobID, 2)
	require.NoError(t, err)
	require.Equal(t, int64(42), got.UpstreamJobID)
	require.JSONEq(t, string(first.BillingSnapshot), string(got.BillingSnapshot))

	second := *first
	second.ID = 0
	second.JobID = "vidjob_test_2"
	second.UpstreamJobID = 43
	second.RequestHash = "request-hash-2"
	require.NoError(t, repo.CreateVideoJob(ctx, &second))

	jobs, err := repo.ListVideoJobsForAPIKey(ctx, 2, 20, "")
	require.NoError(t, err)
	require.Len(t, jobs, 2)
	require.Equal(t, second.JobID, jobs[0].JobID)

	startedAt := time.Now().UTC().Truncate(time.Millisecond)
	require.NoError(t, repo.TransitionVideoJob(ctx, first.JobID, []string{service.VideoJobPending}, service.VideoJobRunning, service.VideoJobTransition{StartedAt: &startedAt}))
	updated, err := repo.GetVideoJob(ctx, first.JobID)
	require.NoError(t, err)
	require.Equal(t, service.VideoJobRunning, updated.Status)
	require.NotNil(t, updated.StartedAt)
	require.WithinDuration(t, startedAt, *updated.StartedAt, time.Millisecond)

	err = repo.TransitionVideoJob(ctx, first.JobID, []string{service.VideoJobPending}, service.VideoJobCompleted, service.VideoJobTransition{})
	require.ErrorIs(t, err, service.ErrVideoJobTransitionConflict)

	_, err = repo.GetVideoJobForAPIKey(ctx, first.JobID, 999)
	require.ErrorIs(t, err, service.ErrVideoJobNotFound)
}
