package query_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/ifql/query"
)

func TestOperation_JSON(t *testing.T) {
	testCases := map[string]struct {
		json string
		op   *query.Operation
	}{
		"select": {
			json: `{
				"id": "select",
				"kind": "select",
				"spec": {
					"database":"mydb"
				}
			}`,
			op: &query.Operation{
				OperationID: "select",
				Spec: &query.SelectOpSpec{
					Database: "mydb",
				},
			},
		},
		"range": {
			json: `{
				"id": "range",
				"kind": "range",
				"spec": {
					"start": "-4h",
					"stop": "now"
				}
			}`,
			op: &query.Operation{
				OperationID: "range",
				Spec: &query.RangeOpSpec{
					Start: query.Time{
						Relative: -4 * time.Hour,
					},
					Stop: query.Time{},
				},
			},
		},
		"clear": {
			json: `{
				"id": "clear",
				"kind": "clear"
			}`,
			op: &query.Operation{
				OperationID: "clear",
				Spec:        &query.ClearOpSpec{},
			},
		},
		"window": {
			json: `{
				"id": "window",
				"kind": "window",
				"spec":{
					"every":"1m",
					"period":"1h",
					"every_count": 100,
					"period_count": 200,
					"start": "2017-08-01T00:00:00Z"
				}
			}`,
			op: &query.Operation{
				OperationID: "window",
				Spec: &query.WindowOpSpec{
					Every:       query.Duration(time.Minute),
					Period:      query.Duration(time.Hour),
					EveryCount:  100,
					PeriodCount: 200,
					Start: query.Time{
						Absolute: time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		//TODO implement full spec unmarshalling for all OpSpecs below
		"merge": {
			json: `{
				"id": "merge",
				"kind": "merge"
			}`,
			op: &query.Operation{
				OperationID: "merge",
				Spec:        &query.MergeOpSpec{},
			},
		},
		"keys": {
			json: `{
				"id": "keys",
				"kind": "keys"
			}`,
			op: &query.Operation{
				OperationID: "keys",
				Spec:        &query.KeysOpSpec{},
			},
		},
		"values": {
			json: `{
				"id": "values",
				"kind": "values"
			}`,
			op: &query.Operation{
				OperationID: "values",
				Spec:        &query.ValuesOpSpec{},
			},
		},
		"cardinality": {
			json: `{
				"id": "cardinality",
				"kind": "cardinality"
			}`,
			op: &query.Operation{
				OperationID: "cardinality",
				Spec:        &query.CardinalityOpSpec{},
			},
		},
		"limit": {
			json: `{
				"id": "limit",
				"kind": "limit"
			}`,
			op: &query.Operation{
				OperationID: "limit",
				Spec:        &query.LimitOpSpec{},
			},
		},
		"shift": {
			json: `{
				"id": "shift",
				"kind": "shift"
			}`,
			op: &query.Operation{
				OperationID: "shift",
				Spec:        &query.ShiftOpSpec{},
			},
		},
		"interpolate": {
			json: `{
				"id": "interpolate",
				"kind": "interpolate"
			}`,
			op: &query.Operation{
				OperationID: "interpolate",
				Spec:        &query.InterpolateOpSpec{},
			},
		},
		"join": {
			json: `{
				"id": "join",
				"kind": "join"
			}`,
			op: &query.Operation{
				OperationID: "join",
				Spec:        &query.JoinOpSpec{},
			},
		},
		"union": {
			json: `{
				"id": "union",
				"kind": "union"
			}`,
			op: &query.Operation{
				OperationID: "union",
				Spec:        &query.UnionOpSpec{},
			},
		},
		"filter": {
			json: `{
				"id": "filter",
				"kind": "filter"
			}`,
			op: &query.Operation{
				OperationID: "filter",
				Spec:        &query.FilterOpSpec{},
			},
		},
		"sort": {
			json: `{
				"id": "sort",
				"kind": "sort"
			}`,
			op: &query.Operation{
				OperationID: "sort",
				Spec:        &query.SortOpSpec{},
			},
		},
		"rate": {
			json: `{
				"id": "rate",
				"kind": "rate"
			}`,
			op: &query.Operation{
				OperationID: "rate",
				Spec:        &query.RateOpSpec{},
			},
		},
		"count": {
			json: `{
				"id": "count",
				"kind": "count"
			}`,
			op: &query.Operation{
				OperationID: "count",
				Spec:        &query.CountOpSpec{},
			},
		},
		"sum": {
			json: `{
				"id": "sum",
				"kind": "sum"
			}`,
			op: &query.Operation{
				OperationID: "sum",
				Spec:        &query.SumOpSpec{},
			},
		},
		"mean": {
			json: `{
				"id": "mean",
				"kind": "mean"
			}`,
			op: &query.Operation{
				OperationID: "mean",
				Spec:        &query.MeanOpSpec{},
			},
		},
		"percentile": {
			json: `{
				"id": "percentile",
				"kind": "percentile"
			}`,
			op: &query.Operation{
				OperationID: "percentile",
				Spec:        &query.PercentileOpSpec{},
			},
		},
		"stddev": {
			json: `{
				"id": "stddev",
				"kind": "stddev"
			}`,
			op: &query.Operation{
				OperationID: "stddev",
				Spec:        &query.StddevOpSpec{},
			},
		},
		"min": {
			json: `{
				"id": "min",
				"kind": "min"
			}`,
			op: &query.Operation{
				OperationID: "min",
				Spec:        &query.MinOpSpec{},
			},
		},
		"max": {
			json: `{
				"id": "max",
				"kind": "max"
			}`,
			op: &query.Operation{
				OperationID: "max",
				Spec:        &query.MaxOpSpec{},
			},
		},
		"top": {
			json: `{
				"id": "top",
				"kind": "top"
			}`,
			op: &query.Operation{
				OperationID: "top",
				Spec:        &query.TopOpSpec{},
			},
		},
		"difference": {
			json: `{
				"id": "difference",
				"kind": "difference"
			}`,
			op: &query.Operation{
				OperationID: "difference",
				Spec:        &query.DifferenceOpSpec{},
			},
		},
	}
	if got, exp := len(testCases), query.NumberOfKinds; got != exp {
		t.Fatalf("unexpected number of test cases, have %d test cases, there are %d kinds", got, exp)
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			// Ensure we can properly unmarshal a spec
			gotO := &query.Operation{}
			if err := json.Unmarshal([]byte(tc.json), &gotO); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(gotO, tc.op) {
				t.Errorf("unexpected operation:\n%s", cmp.Diff(gotO, tc.op))
			}

			// Marshal the spec and ensure we can unmarshal it again.
			data, err := json.Marshal(tc.op)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(data, &gotO); err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(gotO, tc.op) {
				t.Errorf("unexpected operation after marshalling:\n%s", cmp.Diff(gotO, tc.op))
			}
		})
	}
}
