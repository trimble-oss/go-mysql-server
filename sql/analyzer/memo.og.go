// Code generated by optgen; DO NOT EDIT.

package analyzer

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

type crossJoin struct {
	*joinBase
}

var _ relExpr = (*crossJoin)(nil)
var _ joinRel = (*crossJoin)(nil)

func (r *crossJoin) String() string {
	return formatRelExpr(r)
}

func (r *crossJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type innerJoin struct {
	*joinBase
}

var _ relExpr = (*innerJoin)(nil)
var _ joinRel = (*innerJoin)(nil)

func (r *innerJoin) String() string {
	return formatRelExpr(r)
}

func (r *innerJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type leftJoin struct {
	*joinBase
}

var _ relExpr = (*leftJoin)(nil)
var _ joinRel = (*leftJoin)(nil)

func (r *leftJoin) String() string {
	return formatRelExpr(r)
}

func (r *leftJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type semiJoin struct {
	*joinBase
}

var _ relExpr = (*semiJoin)(nil)
var _ joinRel = (*semiJoin)(nil)

func (r *semiJoin) String() string {
	return formatRelExpr(r)
}

func (r *semiJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type antiJoin struct {
	*joinBase
}

var _ relExpr = (*antiJoin)(nil)
var _ joinRel = (*antiJoin)(nil)

func (r *antiJoin) String() string {
	return formatRelExpr(r)
}

func (r *antiJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type lookupJoin struct {
	*joinBase
	lookup *lookup
}

var _ relExpr = (*lookupJoin)(nil)
var _ joinRel = (*lookupJoin)(nil)

func (r *lookupJoin) String() string {
	return formatRelExpr(r)
}

func (r *lookupJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type concatJoin struct {
	*joinBase
	concat []*lookup
}

var _ relExpr = (*concatJoin)(nil)
var _ joinRel = (*concatJoin)(nil)

func (r *concatJoin) String() string {
	return formatRelExpr(r)
}

func (r *concatJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type hashJoin struct {
	*joinBase
	innerAttrs []sql.Expression
	outerAttrs []sql.Expression
}

var _ relExpr = (*hashJoin)(nil)
var _ joinRel = (*hashJoin)(nil)

func (r *hashJoin) String() string {
	return formatRelExpr(r)
}

func (r *hashJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type fullOuterJoin struct {
	*joinBase
}

var _ relExpr = (*fullOuterJoin)(nil)
var _ joinRel = (*fullOuterJoin)(nil)

func (r *fullOuterJoin) String() string {
	return formatRelExpr(r)
}

func (r *fullOuterJoin) joinPrivate() *joinBase {
	return r.joinBase
}

type tableScan struct {
	*relBase
	table *plan.ResolvedTable
}

var _ relExpr = (*tableScan)(nil)
var _ sourceRel = (*tableScan)(nil)

func (r *tableScan) String() string {
	return formatRelExpr(r)
}

func (r *tableScan) name() string {
	return strings.ToLower(r.table.Name())
}

func (r *tableScan) tableId() TableId {
	return tableIdForSource(r.g.id)
}

func (r *tableScan) children() []*exprGroup {
	return nil
}

func (r *tableScan) outputCols() sql.Schema {
	return r.table.Schema()
}

type values struct {
	*relBase
	table *plan.ValueDerivedTable
}

var _ relExpr = (*values)(nil)
var _ sourceRel = (*values)(nil)

func (r *values) String() string {
	return formatRelExpr(r)
}

func (r *values) name() string {
	return strings.ToLower(r.table.Name())
}

func (r *values) tableId() TableId {
	return tableIdForSource(r.g.id)
}

func (r *values) children() []*exprGroup {
	return nil
}

func (r *values) outputCols() sql.Schema {
	return r.table.Schema()
}

type tableAlias struct {
	*relBase
	table *plan.TableAlias
}

var _ relExpr = (*tableAlias)(nil)
var _ sourceRel = (*tableAlias)(nil)

func (r *tableAlias) String() string {
	return formatRelExpr(r)
}

func (r *tableAlias) name() string {
	return strings.ToLower(r.table.Name())
}

func (r *tableAlias) tableId() TableId {
	return tableIdForSource(r.g.id)
}

func (r *tableAlias) children() []*exprGroup {
	return nil
}

func (r *tableAlias) outputCols() sql.Schema {
	return r.table.Schema()
}

type recursiveTable struct {
	*relBase
	table *plan.RecursiveTable
}

var _ relExpr = (*recursiveTable)(nil)
var _ sourceRel = (*recursiveTable)(nil)

func (r *recursiveTable) String() string {
	return formatRelExpr(r)
}

func (r *recursiveTable) name() string {
	return strings.ToLower(r.table.Name())
}

func (r *recursiveTable) tableId() TableId {
	return tableIdForSource(r.g.id)
}

func (r *recursiveTable) children() []*exprGroup {
	return nil
}

func (r *recursiveTable) outputCols() sql.Schema {
	return r.table.Schema()
}

type recursiveCte struct {
	*relBase
	table *plan.RecursiveCte
}

var _ relExpr = (*recursiveCte)(nil)
var _ sourceRel = (*recursiveCte)(nil)

func (r *recursiveCte) String() string {
	return formatRelExpr(r)
}

func (r *recursiveCte) name() string {
	return strings.ToLower(r.table.Name())
}

func (r *recursiveCte) tableId() TableId {
	return tableIdForSource(r.g.id)
}

func (r *recursiveCte) children() []*exprGroup {
	return nil
}

func (r *recursiveCte) outputCols() sql.Schema {
	return r.table.Schema()
}

type subqueryAlias struct {
	*relBase
	table *plan.SubqueryAlias
}

var _ relExpr = (*subqueryAlias)(nil)
var _ sourceRel = (*subqueryAlias)(nil)

func (r *subqueryAlias) String() string {
	return formatRelExpr(r)
}

func (r *subqueryAlias) name() string {
	return strings.ToLower(r.table.Name())
}

func (r *subqueryAlias) tableId() TableId {
	return tableIdForSource(r.g.id)
}

func (r *subqueryAlias) children() []*exprGroup {
	return nil
}

func (r *subqueryAlias) outputCols() sql.Schema {
	return r.table.Schema()
}

func formatRelExpr(r relExpr) string {
	switch r := r.(type) {
	case *crossJoin:
		return fmt.Sprintf("crossJoin %d %d", r.left.id, r.right.id)
	case *innerJoin:
		return fmt.Sprintf("innerJoin %d %d", r.left.id, r.right.id)
	case *leftJoin:
		return fmt.Sprintf("leftJoin %d %d", r.left.id, r.right.id)
	case *semiJoin:
		return fmt.Sprintf("semiJoin %d %d", r.left.id, r.right.id)
	case *antiJoin:
		return fmt.Sprintf("antiJoin %d %d", r.left.id, r.right.id)
	case *lookupJoin:
		return fmt.Sprintf("lookupJoin %d %d", r.left.id, r.right.id)
	case *concatJoin:
		return fmt.Sprintf("concatJoin %d %d", r.left.id, r.right.id)
	case *hashJoin:
		return fmt.Sprintf("hashJoin %d %d", r.left.id, r.right.id)
	case *fullOuterJoin:
		return fmt.Sprintf("fullOuterJoin %d %d", r.left.id, r.right.id)
	case *tableScan:
		return fmt.Sprintf("tableScan: %s", r.name())
	case *values:
		return fmt.Sprintf("values: %s", r.name())
	case *tableAlias:
		return fmt.Sprintf("tableAlias: %s", r.name())
	case *recursiveTable:
		return fmt.Sprintf("recursiveTable: %s", r.name())
	case *recursiveCte:
		return fmt.Sprintf("recursiveCte: %s", r.name())
	case *subqueryAlias:
		return fmt.Sprintf("subqueryAlias: %s", r.name())
	default:
		panic(fmt.Sprintf("unknown relExpr type: %T", r))
	}
}

func buildRelExpr(b *ExecBuilder, r relExpr, input sql.Schema, children ...sql.Node) (sql.Node, error) {
	switch r := r.(type) {
	case *crossJoin:
		return b.buildCrossJoin(r, input, children...)
	case *innerJoin:
		return b.buildInnerJoin(r, input, children...)
	case *leftJoin:
		return b.buildLeftJoin(r, input, children...)
	case *semiJoin:
		return b.buildSemiJoin(r, input, children...)
	case *antiJoin:
		return b.buildAntiJoin(r, input, children...)
	case *lookupJoin:
		return b.buildLookupJoin(r, input, children...)
	case *concatJoin:
		return b.buildConcatJoin(r, input, children...)
	case *hashJoin:
		return b.buildHashJoin(r, input, children...)
	case *fullOuterJoin:
		return b.buildFullOuterJoin(r, input, children...)
	case *tableScan:
		return b.buildTableScan(r, input, children...)
	case *values:
		return b.buildValues(r, input, children...)
	case *tableAlias:
		return b.buildTableAlias(r, input, children...)
	case *recursiveTable:
		return b.buildRecursiveTable(r, input, children...)
	case *recursiveCte:
		return b.buildRecursiveCte(r, input, children...)
	case *subqueryAlias:
		return b.buildSubqueryAlias(r, input, children...)
	default:
		panic(fmt.Sprintf("unknown relExpr type: %T", r))
	}
}
