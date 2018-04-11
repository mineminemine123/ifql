package influxql

import (
	"context"
	"regexp"

	"github.com/influxdata/ifql/query"
)

type Transpiler struct {
	schema Schema
}

func (t *Transpiler) Transpile(ctx context.Context, txt string) (*query.Spec, error) {
	panic("not implemented")
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
