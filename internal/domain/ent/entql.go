// Code generated by ent, DO NOT EDIT.

package ent

import (
	"QuickBrick/internal/domain/ent/pipelineexecutionlog"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/entql"
	"entgo.io/ent/schema/field"
)

// schemaGraph holds a representation of ent/schema at runtime.
var schemaGraph = func() *sqlgraph.Schema {
	graph := &sqlgraph.Schema{Nodes: make([]*sqlgraph.Node, 1)}
	graph.Nodes[0] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   pipelineexecutionlog.Table,
			Columns: pipelineexecutionlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: pipelineexecutionlog.FieldID,
			},
		},
		Type: "PipelineExecutionLog",
		Fields: map[string]*sqlgraph.FieldSpec{
			pipelineexecutionlog.FieldEnv:           {Type: field.TypeString, Column: pipelineexecutionlog.FieldEnv},
			pipelineexecutionlog.FieldType:          {Type: field.TypeString, Column: pipelineexecutionlog.FieldType},
			pipelineexecutionlog.FieldEventType:     {Type: field.TypeString, Column: pipelineexecutionlog.FieldEventType},
			pipelineexecutionlog.FieldPipelineName:  {Type: field.TypeString, Column: pipelineexecutionlog.FieldPipelineName},
			pipelineexecutionlog.FieldUsernameEmail: {Type: field.TypeString, Column: pipelineexecutionlog.FieldUsernameEmail},
			pipelineexecutionlog.FieldCommitID:      {Type: field.TypeString, Column: pipelineexecutionlog.FieldCommitID},
			pipelineexecutionlog.FieldProjectURL:    {Type: field.TypeString, Column: pipelineexecutionlog.FieldProjectURL},
			pipelineexecutionlog.FieldStatus:        {Type: field.TypeString, Column: pipelineexecutionlog.FieldStatus},
			pipelineexecutionlog.FieldCreatedAt:     {Type: field.TypeTime, Column: pipelineexecutionlog.FieldCreatedAt},
		},
	}
	return graph
}()

// predicateAdder wraps the addPredicate method.
// All update, update-one and query builders implement this interface.
type predicateAdder interface {
	addPredicate(func(s *sql.Selector))
}

// addPredicate implements the predicateAdder interface.
func (pelq *PipelineExecutionLogQuery) addPredicate(pred func(s *sql.Selector)) {
	pelq.predicates = append(pelq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the PipelineExecutionLogQuery builder.
func (pelq *PipelineExecutionLogQuery) Filter() *PipelineExecutionLogFilter {
	return &PipelineExecutionLogFilter{config: pelq.config, predicateAdder: pelq}
}

// addPredicate implements the predicateAdder interface.
func (m *PipelineExecutionLogMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the PipelineExecutionLogMutation builder.
func (m *PipelineExecutionLogMutation) Filter() *PipelineExecutionLogFilter {
	return &PipelineExecutionLogFilter{config: m.config, predicateAdder: m}
}

// PipelineExecutionLogFilter provides a generic filtering capability at runtime for PipelineExecutionLogQuery.
type PipelineExecutionLogFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *PipelineExecutionLogFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[0].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int64 predicate on the id field.
func (f *PipelineExecutionLogFilter) WhereID(p entql.Int64P) {
	f.Where(p.Field(pipelineexecutionlog.FieldID))
}

// WhereEnv applies the entql string predicate on the env field.
func (f *PipelineExecutionLogFilter) WhereEnv(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldEnv))
}

// WhereType applies the entql string predicate on the type field.
func (f *PipelineExecutionLogFilter) WhereType(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldType))
}

// WhereEventType applies the entql string predicate on the event_type field.
func (f *PipelineExecutionLogFilter) WhereEventType(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldEventType))
}

// WherePipelineName applies the entql string predicate on the pipeline_name field.
func (f *PipelineExecutionLogFilter) WherePipelineName(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldPipelineName))
}

// WhereUsernameEmail applies the entql string predicate on the username_email field.
func (f *PipelineExecutionLogFilter) WhereUsernameEmail(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldUsernameEmail))
}

// WhereCommitID applies the entql string predicate on the commit_id field.
func (f *PipelineExecutionLogFilter) WhereCommitID(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldCommitID))
}

// WhereProjectURL applies the entql string predicate on the project_url field.
func (f *PipelineExecutionLogFilter) WhereProjectURL(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldProjectURL))
}

// WhereStatus applies the entql string predicate on the status field.
func (f *PipelineExecutionLogFilter) WhereStatus(p entql.StringP) {
	f.Where(p.Field(pipelineexecutionlog.FieldStatus))
}

// WhereCreatedAt applies the entql time.Time predicate on the created_at field.
func (f *PipelineExecutionLogFilter) WhereCreatedAt(p entql.TimeP) {
	f.Where(p.Field(pipelineexecutionlog.FieldCreatedAt))
}
