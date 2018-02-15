// +build !parser_debug

package parser

// //go:generate pigeon -optimize-parser -optimize-grammar -o ifql.go ifql.peg
//go:generate ragel -Z ifql.rl

import (
	"github.com/influxdata/ifql/ast"
)

// NewAST parses ifql query and produces an ast.Program
func NewAST(ifql string) (*ast.Program, error) {
	return nil, nil
}

func Parse(ifql string) (*ast.Program, error) {
	data := []byte(ifql)
	p := &parser{
		m: NewMachine(data),
	}
	return p.parse()
}

type parser struct {
	m Machine
}

func (p *parser) parse() (*ast.Program, error) {
	err := p.m.Scan()
	if err != nil {
		return nil, err
	}

	prog := p.m.pop().(*ast.Program)
	return prog, nil
}
