package schema

import (
	"encoding/json"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// VideoJob stores the durable mapping between a Sub2API video job and LeoStudio.
type VideoJob struct {
	ent.Schema
}

func (VideoJob) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "video_jobs"}}
}

func (VideoJob) Fields() []ent.Field {
	return []ent.Field{
		field.String("job_id").MaxLen(64).Immutable(),
		field.Int64("user_id").Immutable(),
		field.Int64("api_key_id").Immutable(),
		field.Int64("group_id").Immutable(),
		field.Int64("account_id"),
		field.Int64("upstream_job_id").Optional().Nillable(),
		field.String("status").MaxLen(32),
		field.String("requested_model").MaxLen(128).Immutable(),
		field.String("upstream_model").MaxLen(128).Immutable(),
		field.String("prompt").SchemaType(map[string]string{dialect.Postgres: "text"}).Immutable(),
		field.String("resolution").MaxLen(32).Immutable(),
		field.Int("duration_seconds").Immutable(),
		field.String("aspect_ratio").MaxLen(32).Immutable(),
		field.Bool("audio").Default(false).Immutable(),
		field.String("image_source").MaxLen(16).Default("none").Immutable(),
		field.String("image_url").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}).Immutable(),
		field.String("local_input_name").Optional().Nillable().MaxLen(255).Immutable(),
		field.JSON("result", json.RawMessage{}).Optional().SchemaType(map[string]string{dialect.Postgres: "jsonb"}),
		field.String("error_message").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Float("hold_amount").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}),
		field.Float("actual_cost").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}),
		field.JSON("billing_snapshot", json.RawMessage{}).Optional().SchemaType(map[string]string{dialect.Postgres: "jsonb"}).Immutable(),
		field.String("request_hash").MaxLen(128).Immutable(),
		field.Time("settled_at").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("created_at").Immutable().Default(time.Now).SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("started_at").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("finished_at").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (VideoJob) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("job_id").Unique(),
		index.Fields("user_id", "created_at"),
		index.Fields("api_key_id", "created_at"),
		index.Fields("status", "updated_at"),
		index.Fields("account_id", "upstream_job_id").Unique().Annotations(entsql.IndexWhere("upstream_job_id IS NOT NULL")),
	}
}
