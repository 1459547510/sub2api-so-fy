package schema

import (
	"testing"

	"entgo.io/ent/entc/load"
	"github.com/stretchr/testify/require"
)

func TestVideoJobSchema(t *testing.T) {
	spec, err := (&load.Config{Path: "."}).Load()
	require.NoError(t, err)

	schemas := map[string]*load.Schema{}
	for _, item := range spec.Schemas {
		schemas[item.Name] = item
	}

	videoJob := requireSchema(t, schemas, "VideoJob")
	requireSchemaFields(t, videoJob,
		"job_id", "user_id", "api_key_id", "group_id", "account_id", "upstream_job_id",
		"status", "requested_model", "upstream_model", "prompt", "resolution", "duration_seconds",
		"aspect_ratio", "audio", "image_source", "image_url", "local_input_name", "result",
		"error_message", "hold_amount", "actual_cost", "billing_snapshot", "request_hash",
		"settled_at", "created_at", "updated_at", "started_at", "finished_at",
	)
	requireHasUniqueIndex(t, videoJob, "job_id")
	requireHasUniqueIndex(t, videoJob, "account_id", "upstream_job_id")
	requireSchemaIndex(t, videoJob, "api_key_id", "created_at")
	requireSchemaIndex(t, videoJob, "status", "updated_at")
}

func requireSchemaIndex(t *testing.T, schema *load.Schema, fields ...string) {
	t.Helper()
	for _, index := range schema.Indexes {
		if len(index.Fields) != len(fields) {
			continue
		}
		matches := true
		for i := range fields {
			if index.Fields[i] != fields[i] {
				matches = false
				break
			}
		}
		if matches {
			return
		}
	}
	require.Failf(t, "missing index", "schema %s should include index on %v", schema.Name, fields)
}
