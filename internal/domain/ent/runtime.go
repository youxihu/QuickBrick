// Code generated by ent, DO NOT EDIT.

package ent

import (
	"QuickBrick/internal/domain/ent/pipelineexecutionlog"
	"QuickBrick/internal/domain/ent/schema"
	"time"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	pipelineexecutionlogFields := schema.PipelineExecutionLog{}.Fields()
	_ = pipelineexecutionlogFields
	// pipelineexecutionlogDescStatus is the schema descriptor for status field.
	pipelineexecutionlogDescStatus := pipelineexecutionlogFields[8].Descriptor()
	// pipelineexecutionlog.DefaultStatus holds the default value on creation for the status field.
	pipelineexecutionlog.DefaultStatus = pipelineexecutionlogDescStatus.Default.(string)
	// pipelineexecutionlogDescCreatedAt is the schema descriptor for created_at field.
	pipelineexecutionlogDescCreatedAt := pipelineexecutionlogFields[9].Descriptor()
	// pipelineexecutionlog.DefaultCreatedAt holds the default value on creation for the created_at field.
	pipelineexecutionlog.DefaultCreatedAt = pipelineexecutionlogDescCreatedAt.Default.(func() time.Time)
}
