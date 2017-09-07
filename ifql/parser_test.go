package ifql

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/ifql/ast"
)

func TestNewAST(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    interface{}
		wantErr bool
	}{
		{
			name: "select",
			raw:  `select()`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.Identifier{
								Name: "select",
							},
						},
					},
				},
			},
		},
		{
			name: "use variable to declare something",
			raw: `var howdy = 1
			select()`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{
							&ast.VariableDeclarator{
								ID:   &ast.Identifier{Name: "howdy"},
								Init: &ast.NumberLiteral{Value: 1},
							},
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.Identifier{
								Name: "select",
							},
						},
					},
				},
			},
		},
		{
			name: "variable is select statement",
			raw: `var howdy = select()
			howdy.count()`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{
							&ast.VariableDeclarator{
								ID: &ast.Identifier{
									Name: "howdy",
								},
								Init: &ast.CallExpression{
									Callee: &ast.Identifier{
										Name: "select",
									},
								},
							},
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.MemberExpression{
								Object: &ast.Identifier{
									Name: "howdy",
								},
								Property: &ast.Identifier{
									Name: "count",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "two variables for two selects",
			raw: `var howdy = select()
			var doody = select()
			howdy.count()
			doody.sum()`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{
							&ast.VariableDeclarator{
								ID: &ast.Identifier{
									Name: "howdy",
								},
								Init: &ast.CallExpression{
									Callee: &ast.Identifier{
										Name: "select",
									},
								},
							},
						},
					},
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{
							&ast.VariableDeclarator{
								ID: &ast.Identifier{
									Name: "doody",
								},
								Init: &ast.CallExpression{
									Callee: &ast.Identifier{
										Name: "select",
									},
								},
							},
						},
					},

					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.MemberExpression{
								Object: &ast.Identifier{
									Name: "howdy",
								},
								Property: &ast.Identifier{
									Name: "count",
								},
							},
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.MemberExpression{
								Object: &ast.Identifier{
									Name: "doody",
								},
								Property: &ast.Identifier{
									Name: "sum",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "select with database",
			raw:  `select(database:"telegraf")`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.Identifier{
								Name: "select",
							},
							Arguments: []ast.Expression{
								&ast.ObjectExpression{
									Properties: []*ast.Property{
										&ast.Property{
											Key: &ast.Identifier{
												Name: "database",
											},
											Value: &ast.StringLiteral{
												Value: "telegraf",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "select with where with no parens",
			raw:  `select(database:"telegraf").where(exp:{"other"=="mem" and "this"=="that" or "these"!="those"})`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.MemberExpression{
								Property: &ast.Identifier{Name: "where"},
								Object: &ast.CallExpression{
									Callee: &ast.Identifier{
										Name: "select",
									},
									Arguments: []ast.Expression{
										&ast.ObjectExpression{
											Properties: []*ast.Property{
												&ast.Property{
													Key:   &ast.Identifier{Name: "database"},
													Value: &ast.StringLiteral{Value: "telegraf"},
												},
											},
										},
									},
								},
							},
							Arguments: []ast.Expression{
								&ast.ObjectExpression{
									Properties: []*ast.Property{
										&ast.Property{
											Key: &ast.Identifier{Name: "exp"},
											Value: &ast.LogicalExpression{
												Operator: ast.OrOperator,
												Left: &ast.LogicalExpression{
													Operator: ast.AndOperator,
													Left: &ast.BinaryExpression{
														Operator: ast.EqualOperator,
														Left:     &ast.StringLiteral{Value: "other"},
														Right:    &ast.StringLiteral{Value: "mem"},
													},
													Right: &ast.BinaryExpression{
														Operator: ast.EqualOperator,
														Left:     &ast.StringLiteral{Value: "this"},
														Right:    &ast.StringLiteral{Value: "that"},
													},
												},
												Right: &ast.BinaryExpression{
													Operator: ast.NotEqualOperator,
													Left:     &ast.StringLiteral{Value: "these"},
													Right:    &ast.StringLiteral{Value: "those"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "select with range",
			raw:  `select(database:"telegraf").range(start:-1h, end:10m)`,
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Callee: &ast.MemberExpression{
								Object: &ast.CallExpression{
									Callee: &ast.Identifier{Name: "select"},
									Arguments: []ast.Expression{
										&ast.ObjectExpression{
											Properties: []*ast.Property{
												&ast.Property{
													Key:   &ast.Identifier{Name: "database"},
													Value: &ast.StringLiteral{Value: "telegraf"},
												},
											},
										},
									},
								},
								Property: &ast.Identifier{Name: "range"},
							},
							Arguments: []ast.Expression{
								&ast.ObjectExpression{
									Properties: []*ast.Property{
										&ast.Property{
											Key:   &ast.Identifier{Name: "start"},
											Value: &ast.DurationLiteral{Value: -time.Hour},
										},
										&ast.Property{
											Key:   &ast.Identifier{Name: "end"},
											Value: &ast.DurationLiteral{Value: 10 * time.Minute},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Parse("", []byte(tt.raw))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !cmp.Equal(tt.want, got) {
				t.Errorf("Parse() = -got/+want %s", cmp.Diff(tt.want, got))
			}
		})
	}
}
