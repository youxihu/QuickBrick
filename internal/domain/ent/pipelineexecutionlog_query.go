// Code generated by ent, DO NOT EDIT.

package ent

import (
	"QuickBrick/internal/domain/ent/pipelineexecutionlog"
	"QuickBrick/internal/domain/ent/predicate"
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PipelineExecutionLogQuery is the builder for querying PipelineExecutionLog entities.
type PipelineExecutionLogQuery struct {
	config
	ctx        *QueryContext
	order      []pipelineexecutionlog.OrderOption
	inters     []Interceptor
	predicates []predicate.PipelineExecutionLog
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PipelineExecutionLogQuery builder.
func (pelq *PipelineExecutionLogQuery) Where(ps ...predicate.PipelineExecutionLog) *PipelineExecutionLogQuery {
	pelq.predicates = append(pelq.predicates, ps...)
	return pelq
}

// Limit the number of records to be returned by this query.
func (pelq *PipelineExecutionLogQuery) Limit(limit int) *PipelineExecutionLogQuery {
	pelq.ctx.Limit = &limit
	return pelq
}

// Offset to start from.
func (pelq *PipelineExecutionLogQuery) Offset(offset int) *PipelineExecutionLogQuery {
	pelq.ctx.Offset = &offset
	return pelq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pelq *PipelineExecutionLogQuery) Unique(unique bool) *PipelineExecutionLogQuery {
	pelq.ctx.Unique = &unique
	return pelq
}

// Order specifies how the records should be ordered.
func (pelq *PipelineExecutionLogQuery) Order(o ...pipelineexecutionlog.OrderOption) *PipelineExecutionLogQuery {
	pelq.order = append(pelq.order, o...)
	return pelq
}

// First returns the first PipelineExecutionLog entity from the query.
// Returns a *NotFoundError when no PipelineExecutionLog was found.
func (pelq *PipelineExecutionLogQuery) First(ctx context.Context) (*PipelineExecutionLog, error) {
	nodes, err := pelq.Limit(1).All(setContextOp(ctx, pelq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{pipelineexecutionlog.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) FirstX(ctx context.Context) *PipelineExecutionLog {
	node, err := pelq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PipelineExecutionLog ID from the query.
// Returns a *NotFoundError when no PipelineExecutionLog ID was found.
func (pelq *PipelineExecutionLogQuery) FirstID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = pelq.Limit(1).IDs(setContextOp(ctx, pelq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{pipelineexecutionlog.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) FirstIDX(ctx context.Context) int64 {
	id, err := pelq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PipelineExecutionLog entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PipelineExecutionLog entity is found.
// Returns a *NotFoundError when no PipelineExecutionLog entities are found.
func (pelq *PipelineExecutionLogQuery) Only(ctx context.Context) (*PipelineExecutionLog, error) {
	nodes, err := pelq.Limit(2).All(setContextOp(ctx, pelq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{pipelineexecutionlog.Label}
	default:
		return nil, &NotSingularError{pipelineexecutionlog.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) OnlyX(ctx context.Context) *PipelineExecutionLog {
	node, err := pelq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PipelineExecutionLog ID in the query.
// Returns a *NotSingularError when more than one PipelineExecutionLog ID is found.
// Returns a *NotFoundError when no entities are found.
func (pelq *PipelineExecutionLogQuery) OnlyID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = pelq.Limit(2).IDs(setContextOp(ctx, pelq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{pipelineexecutionlog.Label}
	default:
		err = &NotSingularError{pipelineexecutionlog.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) OnlyIDX(ctx context.Context) int64 {
	id, err := pelq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PipelineExecutionLogs.
func (pelq *PipelineExecutionLogQuery) All(ctx context.Context) ([]*PipelineExecutionLog, error) {
	ctx = setContextOp(ctx, pelq.ctx, ent.OpQueryAll)
	if err := pelq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PipelineExecutionLog, *PipelineExecutionLogQuery]()
	return withInterceptors[[]*PipelineExecutionLog](ctx, pelq, qr, pelq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) AllX(ctx context.Context) []*PipelineExecutionLog {
	nodes, err := pelq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PipelineExecutionLog IDs.
func (pelq *PipelineExecutionLogQuery) IDs(ctx context.Context) (ids []int64, err error) {
	if pelq.ctx.Unique == nil && pelq.path != nil {
		pelq.Unique(true)
	}
	ctx = setContextOp(ctx, pelq.ctx, ent.OpQueryIDs)
	if err = pelq.Select(pipelineexecutionlog.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) IDsX(ctx context.Context) []int64 {
	ids, err := pelq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pelq *PipelineExecutionLogQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pelq.ctx, ent.OpQueryCount)
	if err := pelq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pelq, querierCount[*PipelineExecutionLogQuery](), pelq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) CountX(ctx context.Context) int {
	count, err := pelq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pelq *PipelineExecutionLogQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, pelq.ctx, ent.OpQueryExist)
	switch _, err := pelq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (pelq *PipelineExecutionLogQuery) ExistX(ctx context.Context) bool {
	exist, err := pelq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PipelineExecutionLogQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pelq *PipelineExecutionLogQuery) Clone() *PipelineExecutionLogQuery {
	if pelq == nil {
		return nil
	}
	return &PipelineExecutionLogQuery{
		config:     pelq.config,
		ctx:        pelq.ctx.Clone(),
		order:      append([]pipelineexecutionlog.OrderOption{}, pelq.order...),
		inters:     append([]Interceptor{}, pelq.inters...),
		predicates: append([]predicate.PipelineExecutionLog{}, pelq.predicates...),
		// clone intermediate query.
		sql:       pelq.sql.Clone(),
		path:      pelq.path,
		modifiers: append([]func(*sql.Selector){}, pelq.modifiers...),
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Env string `json:"env,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.PipelineExecutionLog.Query().
//		GroupBy(pipelineexecutionlog.FieldEnv).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (pelq *PipelineExecutionLogQuery) GroupBy(field string, fields ...string) *PipelineExecutionLogGroupBy {
	pelq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PipelineExecutionLogGroupBy{build: pelq}
	grbuild.flds = &pelq.ctx.Fields
	grbuild.label = pipelineexecutionlog.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Env string `json:"env,omitempty"`
//	}
//
//	client.PipelineExecutionLog.Query().
//		Select(pipelineexecutionlog.FieldEnv).
//		Scan(ctx, &v)
func (pelq *PipelineExecutionLogQuery) Select(fields ...string) *PipelineExecutionLogSelect {
	pelq.ctx.Fields = append(pelq.ctx.Fields, fields...)
	sbuild := &PipelineExecutionLogSelect{PipelineExecutionLogQuery: pelq}
	sbuild.label = pipelineexecutionlog.Label
	sbuild.flds, sbuild.scan = &pelq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PipelineExecutionLogSelect configured with the given aggregations.
func (pelq *PipelineExecutionLogQuery) Aggregate(fns ...AggregateFunc) *PipelineExecutionLogSelect {
	return pelq.Select().Aggregate(fns...)
}

func (pelq *PipelineExecutionLogQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range pelq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, pelq); err != nil {
				return err
			}
		}
	}
	for _, f := range pelq.ctx.Fields {
		if !pipelineexecutionlog.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if pelq.path != nil {
		prev, err := pelq.path(ctx)
		if err != nil {
			return err
		}
		pelq.sql = prev
	}
	return nil
}

func (pelq *PipelineExecutionLogQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PipelineExecutionLog, error) {
	var (
		nodes = []*PipelineExecutionLog{}
		_spec = pelq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PipelineExecutionLog).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PipelineExecutionLog{config: pelq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(pelq.modifiers) > 0 {
		_spec.Modifiers = pelq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pelq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (pelq *PipelineExecutionLogQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pelq.querySpec()
	if len(pelq.modifiers) > 0 {
		_spec.Modifiers = pelq.modifiers
	}
	_spec.Node.Columns = pelq.ctx.Fields
	if len(pelq.ctx.Fields) > 0 {
		_spec.Unique = pelq.ctx.Unique != nil && *pelq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, pelq.driver, _spec)
}

func (pelq *PipelineExecutionLogQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(pipelineexecutionlog.Table, pipelineexecutionlog.Columns, sqlgraph.NewFieldSpec(pipelineexecutionlog.FieldID, field.TypeInt64))
	_spec.From = pelq.sql
	if unique := pelq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pelq.path != nil {
		_spec.Unique = true
	}
	if fields := pelq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pipelineexecutionlog.FieldID)
		for i := range fields {
			if fields[i] != pipelineexecutionlog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pelq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pelq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pelq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pelq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pelq *PipelineExecutionLogQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pelq.driver.Dialect())
	t1 := builder.Table(pipelineexecutionlog.Table)
	columns := pelq.ctx.Fields
	if len(columns) == 0 {
		columns = pipelineexecutionlog.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pelq.sql != nil {
		selector = pelq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pelq.ctx.Unique != nil && *pelq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range pelq.modifiers {
		m(selector)
	}
	for _, p := range pelq.predicates {
		p(selector)
	}
	for _, p := range pelq.order {
		p(selector)
	}
	if offset := pelq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pelq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pelq *PipelineExecutionLogQuery) Modify(modifiers ...func(s *sql.Selector)) *PipelineExecutionLogSelect {
	pelq.modifiers = append(pelq.modifiers, modifiers...)
	return pelq.Select()
}

// PipelineExecutionLogGroupBy is the group-by builder for PipelineExecutionLog entities.
type PipelineExecutionLogGroupBy struct {
	selector
	build *PipelineExecutionLogQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pelgb *PipelineExecutionLogGroupBy) Aggregate(fns ...AggregateFunc) *PipelineExecutionLogGroupBy {
	pelgb.fns = append(pelgb.fns, fns...)
	return pelgb
}

// Scan applies the selector query and scans the result into the given value.
func (pelgb *PipelineExecutionLogGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pelgb.build.ctx, ent.OpQueryGroupBy)
	if err := pelgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PipelineExecutionLogQuery, *PipelineExecutionLogGroupBy](ctx, pelgb.build, pelgb, pelgb.build.inters, v)
}

func (pelgb *PipelineExecutionLogGroupBy) sqlScan(ctx context.Context, root *PipelineExecutionLogQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pelgb.fns))
	for _, fn := range pelgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pelgb.flds)+len(pelgb.fns))
		for _, f := range *pelgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pelgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pelgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PipelineExecutionLogSelect is the builder for selecting fields of PipelineExecutionLog entities.
type PipelineExecutionLogSelect struct {
	*PipelineExecutionLogQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pels *PipelineExecutionLogSelect) Aggregate(fns ...AggregateFunc) *PipelineExecutionLogSelect {
	pels.fns = append(pels.fns, fns...)
	return pels
}

// Scan applies the selector query and scans the result into the given value.
func (pels *PipelineExecutionLogSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pels.ctx, ent.OpQuerySelect)
	if err := pels.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PipelineExecutionLogQuery, *PipelineExecutionLogSelect](ctx, pels.PipelineExecutionLogQuery, pels, pels.inters, v)
}

func (pels *PipelineExecutionLogSelect) sqlScan(ctx context.Context, root *PipelineExecutionLogQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pels.fns))
	for _, fn := range pels.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pels.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pels.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pels *PipelineExecutionLogSelect) Modify(modifiers ...func(s *sql.Selector)) *PipelineExecutionLogSelect {
	pels.modifiers = append(pels.modifiers, modifiers...)
	return pels
}
