package query

import (
	"fmt"

	"github.com/influxdata/ifql/interpreter"
	"github.com/influxdata/ifql/parser"
	"github.com/influxdata/ifql/semantic"
)

type REPL struct {
	Scope        *interpreter.Scope
	Declarations semantic.DeclarationScope
}

func NewREPL() *REPL {
	return &REPL{
		Scope:        builtinScope.Nest(),
		Declarations: builtinDeclarations.Copy(),
	}
}

func (r *REPL) Input(t string) {
	if t == "" {
		return
	}
	astProg, err := parser.NewAST(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	semProg, err := semantic.New(astProg, r.Declarations)
	if err != nil {
		fmt.Println(err)
		return
	}

	d := new(queryDomain)
	if err := interpreter.Eval(semProg, r.Scope, d); err != nil {
		fmt.Println(err)
		return
	}

	v := r.Scope.Return()
	if v.Type() != semantic.Invalid {
		fmt.Println(v.Value())
	}
}
