package influxql

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/influxdata/ifql/ast"
	"github.com/influxdata/ifql/functions"
	"github.com/influxdata/ifql/query"
	"github.com/influxdata/ifql/query/execute"
	"github.com/influxdata/ifql/semantic"
	"github.com/influxdata/influxql"
)

type Transpiler struct {
	schema Schema
}

func (t *Transpiler) Transpile(ctx context.Context, txt string) (*query.Spec, error) {
	// Parse the text of the query.
	q, err := influxql.ParseQuery(txt)
	if err != nil {
		return nil, err
	}

	if len(q.Statements) != 1 {
		// TODO(jsternberg): Handle queries with multiple statements.
		return nil, errors.New("only one statement is allowed")
	}

	stmt, ok := q.Statements[0].(*influxql.SelectStatement)
	if !ok {
		// TODO(jsternberg): Support meta queries.
		return nil, errors.New("only supports select statements")
	}

	spec := &query.Spec{}
	if err := t.processFrom(spec, stmt); err != nil {
		return nil, err
	}
	if err := t.processCondition(spec, stmt); err != nil {
		return nil, err
	}
	if err := t.processGroupBy(spec, stmt); err != nil {
		return nil, err
	}

	// TODO(jsternberg): Use the query compiler to compile and prepare the query.
	// This will likely involve breaking out the compiler and preparer from influxdb
	// and giving it a public interface.
	// For now, let's just determine if this is an aggregate and go from there.
	if len(stmt.Fields) != 1 {
		return nil, errors.New("unsupported select statement")
	} else if _, ok := stmt.Fields[0].Expr.(*influxql.Call); !ok {
		return nil, errors.New("unsupported select statement")
	}

	if err := t.processFields(spec, stmt); err != nil {
		return nil, err
	}

	return spec, nil
}

func (t *Transpiler) processFrom(spec *query.Spec, stmt *influxql.SelectStatement) error {
	if len(stmt.Sources) != 1 {
		// TODO(jsternberg): Support multiple sources.
		return errors.New("only one source is allowed")
	}

	// Only support a direct measurement. Subqueries are not supported yet.
	mm, ok := stmt.Sources[0].(*influxql.Measurement)
	if !ok {
		return errors.New("source must be a measurement")
	}

	// TODO(jsternberg): Verify the retention policy is the default one so we avoid
	// unexpected behavior.

	// Create the from spec and add it to the list of operations.
	from := functions.FromOpSpec{Database: mm.Database}
	spec.Operations = append(spec.Operations, &query.Operation{
		ID:   "from0",
		Spec: &from,
	})

	// Add a filter for the measurement name.
	measurement := functions.FilterOpSpec{
		Fn: &semantic.FunctionExpression{
			Params: []*semantic.FunctionParam{
				{
					Key: &semantic.Identifier{Name: "r"},
				},
			},
			Body: &semantic.BinaryExpression{
				Operator: ast.EqualOperator,
				Left: &semantic.MemberExpression{
					Object: &semantic.IdentifierExpression{
						Name: "r",
					},
					Property: "_measurement",
				},
				Right: &semantic.StringLiteral{
					Value: mm.Name,
				},
			},
		},
	}
	spec.Operations = append(spec.Operations, &query.Operation{
		ID:   "measurement0",
		Spec: &measurement,
	})
	spec.Edges = append(spec.Edges, query.Edge{
		Parent: "range0",
		Child:  "measurement0",
	})
	return nil
}

func (t *Transpiler) processCondition(spec *query.Spec, stmt *influxql.SelectStatement) error {
	valuer := influxql.NowValuer{
		Now: time.Now(),
	}
	cond, tr, err := influxql.ConditionExpr(stmt.Condition, &valuer)
	if err != nil {
		return fmt.Errorf("could not process condition: %s", err)
	}

	if cond != nil {
		// TODO(jsternberg): Support adding a filter on tags from the condition.
		return errors.New("conditions are unsupported")
	}

	// Add a range spec from the time range. The range will always be there and always be explicit.
	rangeOp := functions.RangeOpSpec{
		Start: query.Time{Absolute: tr.MinTime()},
		Stop:  query.Time{Absolute: tr.MaxTime()},
	}
	spec.Operations = append(spec.Operations, &query.Operation{
		ID:   "range0",
		Spec: &rangeOp,
	})
	spec.Edges = append(spec.Edges, query.Edge{
		Parent: "from0",
		Child:  "range0",
	})
	return nil
}

func (t *Transpiler) processGroupBy(spec *query.Spec, stmt *influxql.SelectStatement) error {
	group := functions.GroupOpSpec{}
	for _, d := range stmt.Dimensions {
		switch expr := d.Expr.(type) {
		case *influxql.VarRef:
			group.By = append(group.By, expr.Val)
		case *influxql.Call:
			if expr.Name == "time" {
				return errors.New("group by time() is unsupported")
			}
		default:
			return errors.New("unsupported")
		}
	}

	spec.Operations = append(spec.Operations, &query.Operation{
		ID:   "group0",
		Spec: &group,
	})
	spec.Edges = append(spec.Edges, query.Edge{
		Parent: "field0",
		Child:  "group0",
	})
	return nil
}

func (t *Transpiler) processFields(spec *query.Spec, stmt *influxql.SelectStatement) error {
	for i, f := range stmt.Fields {
		switch expr := f.Expr.(type) {
		case *influxql.Call:
			if len(expr.Args) != 1 {
				return errors.New("unsupported function with more than one argument")
			}

			ref, ok := expr.Args[0].(*influxql.VarRef)
			if !ok {
				return errors.New("first argument to function must be a variable")
			}

			spec.Operations = append(spec.Operations, &query.Operation{
				ID: query.OperationID(fmt.Sprintf("field%d", i)),
				Spec: &functions.FilterOpSpec{
					Fn: &semantic.FunctionExpression{
						Params: []*semantic.FunctionParam{
							{
								Key: &semantic.Identifier{Name: "r"},
							},
						},
						Body: &semantic.BinaryExpression{
							Operator: ast.EqualOperator,
							Left: &semantic.MemberExpression{
								Object: &semantic.IdentifierExpression{
									Name: "r",
								},
								Property: "_measurement",
							},
							Right: &semantic.StringLiteral{
								Value: ref.Val,
							},
						},
					},
				},
			})
			spec.Edges = append(spec.Edges, query.Edge{
				Parent: "measurement0",
				Child:  query.OperationID(fmt.Sprintf("field%d", i)),
			})

			id := query.OperationID(fmt.Sprintf("%s%d", expr.Name, i))
			switch expr.Name {
			case "max":
				spec.Operations = append(spec.Operations, &query.Operation{
					ID: id,
					Spec: &functions.MaxOpSpec{
						SelectorConfig: execute.SelectorConfig{
							UseRowTime: true,
						},
					},
				})
			case "mean":
				spec.Operations = append(spec.Operations, &query.Operation{
					ID: id,
					Spec: &functions.MeanOpSpec{
						AggregateConfig: execute.AggregateConfig{
							UseStartTime: true,
						},
					},
				})
			default:
				return fmt.Errorf("unsupported function: %s", expr.Name)
			}

			spec.Edges = append(spec.Edges, query.Edge{
				Parent: "group0",
				Child:  id,
			})
		}
	}
	return nil
}

type Type int

const (
	Field Type = iota
	Tag
)

type Schema interface {
	MeasurementsForPattern(pattern *regexp.Regexp) []string
	Type(measurement, id string) Type
}
