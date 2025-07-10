package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"time"
)

// RetryHistory holds the schema definition for the RetryHistory entity.
type RetryHistory struct {
	ent.Schema
}

// Fields of the RetryHistory.
func (RetryHistory) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Unique(),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "datetime", // Override MySQL.
		}).Default(time.Now),

		field.String("env").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("构建环境，如 fe-beta"),

		field.String("project").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("项目名称，如 bbx"),

		field.String("project_url").SchemaType(map[string]string{
			dialect.MySQL: "varchar(512)", // Override MySQL.
		}).Comment("项目地址，如 http://xxx/bbx"),

		field.String("ref").SchemaType(map[string]string{
			dialect.MySQL: "varchar(512)", // Override MySQL.
		}).Comment("分支/Tag 名称，如 refs/heads/master"),

		field.String("event_type").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("事件类型，如 push"),

		field.String("commit_id").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("提交 ID"),

		field.String("committer").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("提交者（名字+邮箱）"),

		field.Text("commit_message").SchemaType(map[string]string{
			dialect.MySQL: "text", // Override MySQL.
		}).Comment("提交信息"),

		field.String("commit_url").SchemaType(map[string]string{
			dialect.MySQL: "varchar(512)", // Override MySQL.
		}).Comment("提交链接"),

		field.String("pipeline_name").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("Pipeline 名称"),

		field.String("pipeline_type").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("Pipeline 类型"),
	}

}

// Edges of the RetryHistory.
func (RetryHistory) Edges() []ent.Edge {
	return nil
}

func (RetryHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "retry_history"},
	}
}
