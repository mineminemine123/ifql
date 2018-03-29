package functions

import (
	"fmt"

	"github.com/influxdata/ifql/query"
	"github.com/influxdata/ifql/query/execute"
	"github.com/influxdata/ifql/query/plan"
)

const CumulativeSumKind = "cumulativeSum"

type CumulativeSumOpSpec struct{}

var cumulativeSumSignature = query.DefaultFunctionSignature()

func init() {
	query.RegisterFunction(CumulativeSumKind, createCumulativeSumOpSpec, cumulativeSumSignature)
	query.RegisterOpSpec(CumulativeSumKind, newCumulativeSumOp)
	plan.RegisterProcedureSpec(CumulativeSumKind, newCumulativeSumProcedure, CumulativeSumKind)
	execute.RegisterTransformation(CumulativeSumKind, createCumulativeSumTransformation)
}

func createCumulativeSumOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(CumulativeSumOpSpec)
	return spec, nil
}

func newCumulativeSumOp() query.OperationSpec {
	return new(CumulativeSumOpSpec)
}

func (s *CumulativeSumOpSpec) Kind() query.OperationKind {
	return CumulativeSumKind
}

type CumulativeSumProcedureSpec struct{}

func newCumulativeSumProcedure(qs query.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	_, ok := qs.(*CumulativeSumOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &CumulativeSumProcedureSpec{}, nil
}

func (s *CumulativeSumProcedureSpec) Kind() plan.ProcedureKind {
	return CumulativeSumKind
}
func (s *CumulativeSumProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(CumulativeSumProcedureSpec)
	*ns = *s
	return ns
}

func createCumulativeSumTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*CumulativeSumProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewBlockBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewCumulativeSumTransformation(d, cache, s)
	return t, d, nil
}

type cumulativeSumTransformation struct {
	d     execute.Dataset
	cache execute.BlockBuilderCache
}

func NewCumulativeSumTransformation(d execute.Dataset, cache execute.BlockBuilderCache, spec *CumulativeSumProcedureSpec) *cumulativeSumTransformation {
	return &cumulativeSumTransformation{
		d:     d,
		cache: cache,
	}
}

func (t *cumulativeSumTransformation) RetractBlock(id execute.DatasetID, meta execute.BlockMetadata) error {
	return t.d.RetractBlock(execute.ToBlockKey(meta))
}

func (t *cumulativeSumTransformation) Process(id execute.DatasetID, b execute.Block) error {
	builder, new := t.cache.BlockBuilder(b)
	if new {
		execute.AddBlockCols(b, builder)
	}
	//timeIdx := execute.TimeIdx(b.Cols())
	//valueIdx := execute.ValueIdx(b.Cols())
	values, err := b.Values()
	if err != nil {
		return err
	}
	values.DoFloat(func(vs []float64, rr execute.RowReader) {
		//builder.AppendFloat(int, float64)
	})
	return nil
}

func (t *cumulativeSumTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}
func (t *cumulativeSumTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}
func (t *cumulativeSumTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
