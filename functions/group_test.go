package functions_test

import (
	"testing"

	"github.com/influxdata/ifql/functions"
	"github.com/influxdata/ifql/query"
	"github.com/influxdata/ifql/query/execute"
	"github.com/influxdata/ifql/query/execute/executetest"
	"github.com/influxdata/ifql/query/plan"
	"github.com/influxdata/ifql/query/plan/plantest"
	"github.com/influxdata/ifql/query/querytest"
)

func TestGroupOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"group","kind":"group","spec":{"by":["t1","t2"],"keep":["t3","t4"]}}`)
	op := &query.Operation{
		ID: "group",
		Spec: &functions.GroupOpSpec{
			By:   []string{"t1", "t2"},
			Keep: []string{"t3", "t4"},
		},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestGroup_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *functions.GroupProcedureSpec
		data []execute.Block
		want []*executetest.Block
	}{
		{
			name: "fan in",
			spec: &functions.GroupProcedureSpec{
				By: []string{"t1"},
			},
			data: []execute.Block{
				&executetest.Block{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a", "x"},
						{execute.Time(2), 1.0, "a", "y"},
					},
				},
				&executetest.Block{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 4.0, "b", "x"},
						{execute.Time(2), 7.0, "b", "y"},
					},
				},
			},
			want: []*executetest.Block{
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a"},
						{execute.Time(2), 1.0, "a"},
					},
				},
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 4.0, "b"},
						{execute.Time(2), 7.0, "b"},
					},
				},
			},
		},
		{
			name: "fan in ignoring",
			spec: &functions.GroupProcedureSpec{
				Except: []string{"t2"},
			},
			data: []execute.Block{
				&executetest.Block{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a", "m", "x"},
						{execute.Time(2), 1.0, "a", "n", "x"},
					},
				},
				&executetest.Block{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 4.0, "b", "m", "x"},
						{execute.Time(2), 7.0, "b", "n", "x"},
					},
				},
			},
			want: []*executetest.Block{
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a", "x"},
						{execute.Time(2), 1.0, "a", "x"},
					},
				},
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 4.0, "b", "x"},
						{execute.Time(2), 7.0, "b", "x"},
					},
				},
			},
		},
		{
			name: "fan in ignoring with keep",
			spec: &functions.GroupProcedureSpec{
				Except: []string{"t2"},
				Keep:   []string{"t2"},
			},
			data: []execute.Block{
				&executetest.Block{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a", "m", "x"},
						{execute.Time(2), 1.0, "a", "n", "x"},
					},
				},
				&executetest.Block{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 4.0, "b", "m", "x"},
						{execute.Time(2), 7.0, "b", "n", "x"},
					},
				},
			},
			want: []*executetest.Block{
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a", "x", "m"},
						{execute.Time(2), 1.0, "a", "x", "n"},
					},
				},
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
					},
					Data: [][]interface{}{
						{execute.Time(1), 4.0, "b", "x", "m"},
						{execute.Time(2), 7.0, "b", "x", "n"},
					},
				},
			},
		},
		{
			name: "fan out",
			spec: &functions.GroupProcedureSpec{
				By: []string{"t1"},
			},
			data: []execute.Block{&executetest.Block{
				Bnds: execute.Bounds{
					Start: 1,
					Stop:  3,
				},
				ColMeta: []execute.ColMeta{
					{Label: "time", Type: execute.TTime},
					{Label: "value", Type: execute.TFloat},
					{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: false},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0, "a"},
					{execute.Time(2), 1.0, "b"},
				},
			}},
			want: []*executetest.Block{
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a"},
					},
				},
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(2), 1.0, "b"},
					},
				},
			},
		},
		{
			name: "fan out ignoring",
			spec: &functions.GroupProcedureSpec{
				Except: []string{"t2"},
			},
			data: []execute.Block{&executetest.Block{
				Bnds: execute.Bounds{
					Start: 1,
					Stop:  3,
				},
				ColMeta: []execute.ColMeta{
					{Label: "time", Type: execute.TTime},
					{Label: "value", Type: execute.TFloat},
					{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
					{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
					{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: false},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0, "a", "m", "x"},
					{execute.Time(2), 1.0, "a", "n", "y"},
				},
			}},
			want: []*executetest.Block{
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 2.0, "a", "x"},
					},
				},
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(2), 1.0, "a", "y"},
					},
				},
			},
		},
		{
			name: "fan out ignoring with keep",
			spec: &functions.GroupProcedureSpec{
				Except: []string{"t2"},
				Keep:   []string{"t2"},
			},
			data: []execute.Block{&executetest.Block{
				Bnds: execute.Bounds{
					Start: 1,
					Stop:  3,
				},
				ColMeta: []execute.ColMeta{
					{Label: "time", Type: execute.TTime},
					{Label: "value", Type: execute.TFloat},
					{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
					{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
					{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: false},
				},
				Data: [][]interface{}{
					{execute.Time(1), 3.0, "a", "m", "x"},
					{execute.Time(2), 2.0, "a", "n", "x"},
					{execute.Time(3), 1.0, "a", "m", "y"},
					{execute.Time(4), 0.0, "a", "n", "y"},
				},
			}},
			want: []*executetest.Block{
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(1), 3.0, "a", "m", "x"},
						{execute.Time(2), 2.0, "a", "n", "x"},
					},
				},
				{
					Bnds: execute.Bounds{
						Start: 1,
						Stop:  3,
					},
					ColMeta: []execute.ColMeta{
						{Label: "time", Type: execute.TTime},
						{Label: "value", Type: execute.TFloat},
						{Label: "t1", Type: execute.TString, IsTag: true, IsCommon: true},
						{Label: "t2", Type: execute.TString, IsTag: true, IsCommon: false},
						{Label: "t3", Type: execute.TString, IsTag: true, IsCommon: true},
					},
					Data: [][]interface{}{
						{execute.Time(3), 1.0, "a", "m", "y"},
						{execute.Time(4), 0.0, "a", "n", "y"},
					},
				},
			},
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
					return functions.NewGroupTransformation(d, c, tc.spec)
				},
			)
		})
	}
}
func TestGroup_PushDown_Single(t *testing.T) {
	lp := &plan.LogicalPlanSpec{
		Procedures: map[plan.ProcedureID]*plan.Procedure{
			plan.ProcedureIDFromOperationID("from"): {
				ID: plan.ProcedureIDFromOperationID("from"),
				Spec: &functions.FromProcedureSpec{
					Database: "mydb",
				},
				Parents:  nil,
				Children: []plan.ProcedureID{plan.ProcedureIDFromOperationID("range")},
			},
			plan.ProcedureIDFromOperationID("range"): {
				ID: plan.ProcedureIDFromOperationID("range"),
				Spec: &functions.RangeProcedureSpec{
					Bounds: plan.BoundsSpec{
						Stop: query.Now,
					},
				},
				Parents:  []plan.ProcedureID{plan.ProcedureIDFromOperationID("from")},
				Children: []plan.ProcedureID{plan.ProcedureIDFromOperationID("group")},
			},
			plan.ProcedureIDFromOperationID("group"): {
				ID: plan.ProcedureIDFromOperationID("group"),
				Spec: &functions.GroupProcedureSpec{
					By:     []string{"a", "b"},
					Keep:   []string{"c", "d"},
					Except: []string{"e", "f"},
				},
				Parents:  []plan.ProcedureID{plan.ProcedureIDFromOperationID("range")},
				Children: nil,
			},
		},
		Order: []plan.ProcedureID{
			plan.ProcedureIDFromOperationID("from"),
			plan.ProcedureIDFromOperationID("range"),
			plan.ProcedureIDFromOperationID("group"),
		},
	}

	want := &plan.PlanSpec{
		Bounds: plan.BoundsSpec{
			Stop: query.Now,
		},
		Procedures: map[plan.ProcedureID]*plan.Procedure{
			plan.ProcedureIDFromOperationID("from"): {
				ID: plan.ProcedureIDFromOperationID("from"),
				Spec: &functions.FromProcedureSpec{
					Database:  "mydb",
					BoundsSet: true,
					Bounds: plan.BoundsSpec{
						Stop: query.Now,
					},
					GroupingSet: true,
					GroupKeys:   []string{"a", "b"},
					GroupKeep:   []string{"c", "d"},
					GroupExcept: []string{"e", "f"},
				},
				Children: []plan.ProcedureID{},
			},
		},
		Results: []plan.ProcedureID{
			(plan.ProcedureIDFromOperationID("from")),
		},
		Order: []plan.ProcedureID{
			plan.ProcedureIDFromOperationID("from"),
		},
	}

	plantest.PhysicalPlanTestHelper(t, lp, want)
}

func TestGroup_PushDown_Branch(t *testing.T) {
	lp := &plan.LogicalPlanSpec{
		Procedures: map[plan.ProcedureID]*plan.Procedure{
			plan.ProcedureIDFromOperationID("from"): {
				ID: plan.ProcedureIDFromOperationID("from"),
				Spec: &functions.FromProcedureSpec{
					Database: "mydb",
				},
				Parents: nil,
				Children: []plan.ProcedureID{
					plan.ProcedureIDFromOperationID("range"),
				},
			},
			plan.ProcedureIDFromOperationID("range"): {
				ID: plan.ProcedureIDFromOperationID("range"),
				Spec: &functions.RangeProcedureSpec{
					Bounds: plan.BoundsSpec{
						Stop: query.Now,
					},
				},
				Parents: []plan.ProcedureID{plan.ProcedureIDFromOperationID("from")},
				Children: []plan.ProcedureID{
					plan.ProcedureIDFromOperationID("groupA"),
					plan.ProcedureIDFromOperationID("groupB"),
				},
			},
			plan.ProcedureIDFromOperationID("groupA"): {
				ID: plan.ProcedureIDFromOperationID("groupA"),
				Spec: &functions.GroupProcedureSpec{
					By:     []string{"a", "b"},
					Keep:   []string{"c", "d"},
					Except: []string{"e", "f"},
				},
				Parents: []plan.ProcedureID{
					plan.ProcedureIDFromOperationID("range"),
				},
				Children: nil,
			},
			plan.ProcedureIDFromOperationID("groupB"): {
				ID: plan.ProcedureIDFromOperationID("groupB"),
				Spec: &functions.GroupProcedureSpec{
					Keep: []string{"C", "D"},
				},
				Parents: []plan.ProcedureID{
					plan.ProcedureIDFromOperationID("range"),
				},
				Children: nil,
			},
		},
		Order: []plan.ProcedureID{
			plan.ProcedureIDFromOperationID("from"),
			plan.ProcedureIDFromOperationID("range"),
			plan.ProcedureIDFromOperationID("groupA"),
			plan.ProcedureIDFromOperationID("groupB"), // groupB is last so it will be duplicated
		},
	}

	fromID := plan.ProcedureIDFromOperationID("from")
	fromIDDup := plan.ProcedureIDForDuplicate(fromID)
	want := &plan.PlanSpec{
		Bounds: plan.BoundsSpec{
			Stop: query.Now,
		},
		Procedures: map[plan.ProcedureID]*plan.Procedure{
			fromID: {
				ID: fromID,
				Spec: &functions.FromProcedureSpec{
					Database:  "mydb",
					BoundsSet: true,
					Bounds: plan.BoundsSpec{
						Stop: query.Now,
					},
					GroupingSet: true,
					GroupKeys:   []string{"a", "b"},
					GroupKeep:   []string{"c", "d"},
					GroupExcept: []string{"e", "f"},
				},
				Children: []plan.ProcedureID{},
			},
			fromIDDup: {
				ID: fromIDDup,
				Spec: &functions.FromProcedureSpec{
					Database:  "mydb",
					BoundsSet: true,
					Bounds: plan.BoundsSpec{
						Stop: query.Now,
					},
					GroupingSet: true,
					MergeAll:    true,
					GroupKeys:   []string{},
					GroupKeep:   []string{"C", "D"},
					GroupExcept: []string{},
				},
				Parents:  []plan.ProcedureID{},
				Children: []plan.ProcedureID{},
			},
		},
		Results: []plan.ProcedureID{
			fromID,
			fromIDDup,
		},
		Order: []plan.ProcedureID{
			fromID,
			fromIDDup,
		},
	}

	plantest.PhysicalPlanTestHelper(t, lp, want)
}
