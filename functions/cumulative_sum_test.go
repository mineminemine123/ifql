package functions_test

import (
	"testing"

	"github.com/influxdata/ifql/functions"
	"github.com/influxdata/ifql/query"
	"github.com/influxdata/ifql/query/execute"
	"github.com/influxdata/ifql/query/execute/executetest"
	"github.com/influxdata/ifql/query/querytest"
)

func TestCumulativeSumOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"cumulativeSum","kind":"cumulativeSum","spec":{}}`)
	op := &query.Operation{
		ID:   "cumulativeSum",
		Spec: &functions.CumulativeSumOpSpec{},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestCumulativeSum_PassThrough(t *testing.T) {
	executetest.TransformationPassThroughTestHelper(t, func(d execute.Dataset, c execute.BlockBuilderCache) execute.Transformation {
		s := functions.NewCumulativeSumTransformation(
			d,
			c,
			&functions.CumulativeSumProcedureSpec{},
		)
		return s
	})
}

func TestCumulativeSum_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *functions.CumulativeSumProcedureSpec
		data []execute.Block
		want []*executetest.Block
	}{
		{
			name: "float",
			spec: &functions.CumulativeSumProcedureSpec{},
			data: []execute.Block{&executetest.Block{
				Bnds: execute.Bounds{
					Start: 0,
					Stop:  10,
				},
				ColMeta: []execute.ColMeta{
					{Label: "_time", Type: execute.TTime, Kind: execute.TimeColKind},
					{Label: "_value", Type: execute.TFloat, Kind: execute.ValueColKind},
				},
				Data: [][]interface{}{
					{execute.Time(0), 2.0},
					{execute.Time(1), 1.0},
					{execute.Time(2), 3.0},
					{execute.Time(3), 4.0},
					{execute.Time(4), 2.0},
					{execute.Time(5), 6.0},
					{execute.Time(6), 2.0},
					{execute.Time(7), 7.0},
					{execute.Time(8), 3.0},
					{execute.Time(9), 8.0},
				},
			}},
			want: []*executetest.Block{{
				Bnds: execute.Bounds{
					Start: 0,
					Stop:  10,
				},
				ColMeta: []execute.ColMeta{
					{Label: "_time", Type: execute.TTime, Kind: execute.TimeColKind},
					{Label: "_value", Type: execute.TFloat, Kind: execute.ValueColKind},
				},
				Data: [][]interface{}{
					// TODO: Update _value for expected cumulative sum
					{execute.Time(0), 0},
					{execute.Time(1), 0},
					{execute.Time(2), 0},
					{execute.Time(3), 0},
					{execute.Time(4), 0},
					{execute.Time(5), 0},
					{execute.Time(6), 0},
					{execute.Time(7), 0},
					{execute.Time(8), 0},
					{execute.Time(9), 0},
				},
			}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.ProcessTestHelper(
				t,
				tc.data,
				tc.want,
				func(d execute.Dataset, c execute.BlockBuilderCache) execute.Transformation {
					return functions.NewCumulativeSumTransformation(d, c, tc.spec)
				},
			)
		})
	}
}
