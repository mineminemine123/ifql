package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/ifql/ast"
	"github.com/influxdata/ifql/ast/asttest"
	"github.com/influxdata/ifql/parser"
)

func TestRagel(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    *ast.Program
		wantErr bool
	}{
		{
			name: "var decl",
			raw:  `x = a`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{{
							ID:   &ast.Identifier{Name: "x"},
							Init: &ast.Identifier{Name: "a"},
						}},
					},
				},
			},
		},
		{
			name: "return statement",
			raw:  `return a`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.ReturnStatement{
						Argument: &ast.Identifier{Name: "a"},
					},
				},
			},
		},
		{
			name: "block statement",
			raw:  `{ a }`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.BlockStatement{
						Body: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Identifier{Name: "a"},
							},
						},
					},
				},
			},
		},
		{
			name: "identifier",
			raw:  `abc`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.Identifier{Name: "abc"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ifql.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !cmp.Equal(tt.want, got, asttest.CompareOptions...) {
				t.Errorf("ifql.Parse() = -want/+got %s", cmp.Diff(tt.want, got, asttest.CompareOptions...))
			}
		})
	}
}
