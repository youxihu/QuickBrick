package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"time"
)

// PipelineExecutionLog holds the schema definition for the PipelineExecutionLog entity.
type PipelineExecutionLog struct {
	ent.Schema
}

// Fields of the PipelineExecutionLog.
func (PipelineExecutionLog) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").SchemaType(map[string]string{
			dialect.MySQL: "bigint", // Override MySQL.
		}).Unique(),

		field.String("env").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}),

		field.String("type").SchemaType(map[string]string{
			dialect.MySQL: "varchar(50)", // Override MySQL.
		}),

		field.String("event_type").SchemaType(map[string]string{
			dialect.MySQL: "varchar(50)", // Override MySQL.
		}),

		field.String("pipeline_name").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}),

		field.String("username_email").SchemaType(map[string]string{
			dialect.MySQL: "varchar(512)", // Override MySQL.
		}),

		field.String("commit_id").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}),

		field.String("project_url").SchemaType(map[string]string{
			dialect.MySQL: "varchar(512)", // Override MySQL.
		}),

		field.String("status").SchemaType(map[string]string{
			dialect.MySQL: "varchar(20)", // Override MySQL.
		}).Default("unknown"),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "datetime", // Override MySQL.
		}).Default(time.Now),
	}

}

// Edges of the PipelineExecutionLog.
func (PipelineExecutionLog) Edges() []ent.Edge {
	return nil
}

func (PipelineExecutionLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pipeline_execution_log"},
	}
}
