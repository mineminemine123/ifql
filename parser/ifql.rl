package parser

import (
	"fmt"
	"log"

	"github.com/influxdata/ifql/ast"
)

var (
	blockMarker = (*ast.BlockStatement)(nil)
)

%%{
	machine ifql;
	action start {
		m.ts = m.p
		log.Println("start", m.ts)
	}
	action finish {
		m.te = m.p
		log.Println("finish", m.ts, m.te)
	}

	action identifier {
		log.Println("action identifier", m.stack)

		m.push(&ast.Identifier{
			Name: string(m.current()),
		})
	}

	action variableStatement {
		log.Println("action variableStatement", m.stack)

		init := m.pop().(ast.Expression)
		id := m.pop().(*ast.Identifier)
		s := &ast.VariableDeclaration{
			Declarations: []*ast.VariableDeclarator{{
				ID: id,
				Init: init,
			}},
		}
		m.push(s)
	}
	action expressionStatement {
		log.Println("action expressionStatement", m.stack)

		expr := m.pop().(ast.Expression)
		s := &ast.ExpressionStatement{
			Expression: expr,
		}
		m.push(s)
	}
	action returnStatement {
		log.Println("action returnStatement", m.stack)

		expr := m.pop().(ast.Expression)
		s := &ast.ReturnStatement{
			Argument: expr,
		}
		m.push(s)
	}

	action beginBlockStatement {
		log.Println("action beginBlockStatement", m.stack)
		m.push(blockMarker)
	}
	action endBlockStatement {
		log.Println("action endBlockStatement", m.stack)

		s := &ast.BlockStatement{}
		for {
			n := m.pop()
			if n == blockMarker {
				break
			}
			stmt := n.(ast.Statement)
			s.Body = append(s.Body, stmt)
		}
		m.push(s)
	}

	action program {
		log.Println("action program", m.stack)

		p := &ast.Program{}
		for m.len() > 0 {
			stmt := m.pop().(ast.Statement)
			p.Body = append(p.Body, stmt)
		}
		m.push(p)
	}

	identifier = /[A-Za-z_][0-9A-Za-z_]*/ >start %finish %identifier;

	# Expressions
	expression = identifier;

	# Statements
	variableStatement = identifier space '=' space expression %variableStatement;
	expressionStatement = expression %expressionStatement;
	returnStatement = 'return' expression %returnStatement;
	blockStatement = '{' expressionStatement* '}' >beginBlockStatement %endBlockStatement;

	statement = |*
		returnStatement;
		blockStatement;
		variableStatement;
		expressionStatement;
	*|;

	program =  statement+ %program;

	main := program;

	write data;
}%%

type Machine struct {
	data       []byte
	cs         int
	ts, te     int
	p, pe, eof int

	stack []ast.Node
}

func NewMachine(data []byte) Machine {
	m := Machine{
		data: data,
		pe: len(data),
		eof: len(data),
	}

	%% access m.;
	%% variable p m.p;
	%% variable pe m.pe;
	%% variable eof m.eof;
	%% variable data m.data;
	%% write init;

	return m
}

func (m *Machine) Scan() (error) {
	if m.p >= m.pe {
		return fmt.Errorf("EOF")
	}

	%% write exec;

	if m.len() != 1 {
		log.Println(m.stack)
		return fmt.Errorf("failed to parse")
	}

	return nil
}

func (m *Machine) push(n ast.Node) {
	m.stack = append(m.stack, n)
}

func (m *Machine) len() int {
	return len(m.stack)
}

func (m *Machine) pop() ast.Node {
	e := len(m.stack) - 1
	n := m.stack[e]
	m.stack = m.stack[:e]
	return n
}

func (m *Machine) current() []byte {
	return m.data[m.ts:m.te]
}
