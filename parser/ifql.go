
//line ifql.rl:1
package parser

import (
	"fmt"
	"log"

	"github.com/influxdata/ifql/ast"
)

var (
	blockMarker = (*ast.BlockStatement)(nil)
)


//line ifql.go:18
var _ifql_actions []byte = []byte{
	0, 1, 0, 1, 6, 1, 7, 2, 1, 
	2, 2, 7, 0, 2, 7, 8, 3, 
	1, 2, 3, 3, 1, 2, 4, 4, 
	0, 1, 2, 4, 4, 1, 2, 3, 
	0, 4, 1, 2, 3, 4, 4, 1, 
	2, 3, 8, 4, 1, 2, 4, 0, 
	4, 1, 2, 4, 6, 4, 1, 2, 
	4, 8, 4, 1, 2, 5, 4, 5, 
	0, 1, 2, 3, 4, 5, 0, 1, 
	2, 5, 4, 5, 1, 2, 3, 4, 
	0, 5, 1, 2, 3, 4, 8, 5, 
	1, 2, 5, 3, 4, 5, 1, 2, 
	5, 4, 0, 5, 1, 2, 5, 4, 
	8, 6, 0, 1, 2, 5, 3, 4, 
	6, 1, 2, 5, 3, 4, 0, 6, 
	1, 2, 5, 3, 4, 8, 
}

var _ifql_key_offsets []int16 = []int16{
	0, 0, 7, 8, 11, 16, 22, 30, 
	42, 51, 63, 76, 89, 102, 114, 128, 
	140, 152, 165, 178, 191, 203, 217, 229, 
	236, 249, 262, 275, 287, 301, 313, 325, 
	338, 351, 364, 376, 390, 
}

var _ifql_trans_keys []byte = []byte{
	95, 114, 123, 65, 90, 97, 122, 61, 
	32, 9, 13, 95, 65, 90, 97, 122, 
	95, 125, 65, 90, 97, 122, 95, 125, 
	48, 57, 65, 90, 97, 122, 32, 95, 
	114, 123, 9, 13, 48, 57, 65, 90, 
	97, 122, 95, 114, 123, 48, 57, 65, 
	90, 97, 122, 32, 95, 114, 123, 9, 
	13, 48, 57, 65, 90, 97, 122, 32, 
	95, 101, 114, 123, 9, 13, 48, 57, 
	65, 90, 97, 122, 32, 95, 114, 116, 
	123, 9, 13, 48, 57, 65, 90, 97, 
	122, 32, 95, 114, 117, 123, 9, 13, 
	48, 57, 65, 90, 97, 122, 32, 95, 
	114, 123, 9, 13, 48, 57, 65, 90, 
	97, 122, 32, 95, 101, 110, 114, 123, 
	9, 13, 48, 57, 65, 90, 97, 122, 
	32, 95, 114, 123, 9, 13, 48, 57, 
	65, 90, 97, 122, 32, 95, 114, 123, 
	9, 13, 48, 57, 65, 90, 97, 122, 
	32, 95, 101, 114, 123, 9, 13, 48, 
	57, 65, 90, 97, 122, 32, 95, 114, 
	116, 123, 9, 13, 48, 57, 65, 90, 
	97, 122, 32, 95, 114, 117, 123, 9, 
	13, 48, 57, 65, 90, 97, 122, 32, 
	95, 114, 123, 9, 13, 48, 57, 65, 
	90, 97, 122, 32, 95, 101, 110, 114, 
	123, 9, 13, 48, 57, 65, 90, 97, 
	122, 32, 95, 114, 123, 9, 13, 48, 
	57, 65, 90, 97, 122, 95, 114, 123, 
	65, 90, 97, 122, 32, 95, 101, 114, 
	123, 9, 13, 48, 57, 65, 90, 97, 
	122, 32, 95, 114, 116, 123, 9, 13, 
	48, 57, 65, 90, 97, 122, 32, 95, 
	114, 117, 123, 9, 13, 48, 57, 65, 
	90, 97, 122, 32, 95, 114, 123, 9, 
	13, 48, 57, 65, 90, 97, 122, 32, 
	95, 101, 110, 114, 123, 9, 13, 48, 
	57, 65, 90, 97, 122, 32, 95, 114, 
	123, 9, 13, 48, 57, 65, 90, 97, 
	122, 32, 95, 114, 123, 9, 13, 48, 
	57, 65, 90, 97, 122, 32, 95, 101, 
	114, 123, 9, 13, 48, 57, 65, 90, 
	97, 122, 32, 95, 114, 116, 123, 9, 
	13, 48, 57, 65, 90, 97, 122, 32, 
	95, 114, 117, 123, 9, 13, 48, 57, 
	65, 90, 97, 122, 32, 95, 114, 123, 
	9, 13, 48, 57, 65, 90, 97, 122, 
	32, 95, 101, 110, 114, 123, 9, 13, 
	48, 57, 65, 90, 97, 122, 32, 95, 
	114, 123, 9, 13, 48, 57, 65, 90, 
	97, 122, 
}

var _ifql_single_lengths []byte = []byte{
	0, 3, 1, 1, 1, 2, 2, 4, 
	3, 4, 5, 5, 5, 4, 6, 4, 
	4, 5, 5, 5, 4, 6, 4, 3, 
	5, 5, 5, 4, 6, 4, 4, 5, 
	5, 5, 4, 6, 4, 
}

var _ifql_range_lengths []byte = []byte{
	0, 2, 0, 1, 2, 2, 3, 4, 
	3, 4, 4, 4, 4, 4, 4, 4, 
	4, 4, 4, 4, 4, 4, 4, 2, 
	4, 4, 4, 4, 4, 4, 4, 4, 
	4, 4, 4, 4, 4, 
}

var _ifql_index_offsets []int16 = []int16{
	0, 0, 6, 8, 11, 15, 20, 26, 
	35, 42, 51, 61, 71, 81, 90, 101, 
	110, 119, 129, 139, 149, 158, 169, 178, 
	184, 194, 204, 214, 223, 234, 243, 252, 
	262, 272, 282, 291, 302, 
}

var _ifql_indicies []byte = []byte{
	0, 2, 3, 0, 0, 1, 4, 1, 
	5, 5, 1, 6, 6, 6, 1, 7, 
	8, 7, 7, 1, 10, 11, 9, 10, 
	10, 1, 12, 14, 15, 16, 12, 13, 
	14, 14, 1, 18, 19, 20, 17, 18, 
	18, 1, 12, 22, 23, 24, 12, 21, 
	22, 22, 1, 12, 22, 25, 23, 24, 
	12, 21, 22, 22, 1, 12, 22, 23, 
	26, 24, 12, 21, 22, 22, 1, 12, 
	22, 23, 27, 24, 12, 21, 22, 22, 
	1, 12, 22, 28, 24, 12, 21, 22, 
	22, 1, 12, 22, 25, 29, 23, 24, 
	12, 21, 22, 22, 1, 12, 30, 31, 
	24, 12, 21, 30, 30, 1, 12, 33, 
	34, 35, 12, 32, 33, 33, 1, 12, 
	33, 36, 34, 35, 12, 32, 33, 33, 
	1, 12, 33, 34, 37, 35, 12, 32, 
	33, 33, 1, 12, 33, 34, 38, 35, 
	12, 32, 33, 33, 1, 12, 33, 39, 
	35, 12, 32, 33, 33, 1, 12, 33, 
	36, 40, 34, 35, 12, 32, 33, 33, 
	1, 12, 41, 42, 35, 12, 32, 41, 
	41, 1, 43, 44, 45, 43, 43, 1, 
	12, 14, 46, 15, 16, 12, 13, 14, 
	14, 1, 12, 14, 15, 47, 16, 12, 
	13, 14, 14, 1, 12, 14, 15, 48, 
	16, 12, 13, 14, 14, 1, 12, 14, 
	49, 16, 12, 13, 14, 14, 1, 12, 
	14, 46, 50, 15, 16, 12, 13, 14, 
	14, 1, 12, 51, 52, 16, 12, 13, 
	51, 51, 1, 12, 54, 55, 56, 12, 
	53, 54, 54, 1, 12, 54, 57, 55, 
	56, 12, 53, 54, 54, 1, 12, 54, 
	55, 58, 56, 12, 53, 54, 54, 1, 
	12, 54, 55, 59, 56, 12, 53, 54, 
	54, 1, 12, 54, 60, 56, 12, 53, 
	54, 54, 1, 12, 54, 57, 61, 55, 
	56, 12, 53, 54, 54, 1, 12, 62, 
	63, 56, 12, 53, 62, 62, 1, 
}

var _ifql_trans_targs []byte = []byte{
	7, 0, 24, 5, 3, 4, 8, 6, 
	23, 6, 6, 23, 2, 7, 7, 24, 
	5, 8, 9, 10, 5, 9, 9, 10, 
	5, 11, 12, 13, 14, 15, 16, 17, 
	16, 16, 17, 5, 18, 19, 20, 21, 
	22, 16, 17, 7, 24, 5, 25, 26, 
	27, 28, 29, 30, 31, 30, 30, 31, 
	5, 32, 33, 34, 35, 36, 30, 31, 
}

var _ifql_trans_actions []byte = []byte{
	1, 0, 1, 0, 0, 0, 1, 1, 
	3, 0, 44, 49, 7, 0, 44, 44, 
	20, 0, 29, 29, 16, 0, 76, 76, 
	34, 76, 76, 76, 76, 76, 64, 64, 
	0, 113, 113, 88, 113, 113, 113, 113, 
	113, 106, 106, 10, 10, 5, 44, 44, 
	44, 44, 44, 24, 24, 0, 94, 94, 
	59, 94, 94, 94, 94, 94, 70, 70, 
}

var _ifql_eof_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 54, 
	39, 82, 82, 82, 82, 82, 82, 82, 
	120, 120, 120, 120, 120, 120, 120, 13, 
	54, 54, 54, 54, 54, 54, 100, 100, 
	100, 100, 100, 100, 100, 
}

const ifql_start int = 1
const ifql_first_final int = 7
const ifql_error int = 0

const ifql_en_main int = 1


//line ifql.rl:115


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

	
//line ifql.rl:134
	
//line ifql.rl:135
	
//line ifql.rl:136
	
//line ifql.rl:137
	
//line ifql.rl:138
	
//line ifql.go:233
	{
	 m.cs = ifql_start
	}

//line ifql.rl:139

	return m
}

func (m *Machine) Scan() (error) {
	if m.p >= m.pe {
		return fmt.Errorf("EOF")
	}

	
//line ifql.go:249
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if ( m.p) == ( m.pe) {
		goto _test_eof
	}
	if  m.cs == 0 {
		goto _out
	}
_resume:
	_keys = int(_ifql_key_offsets[ m.cs])
	_trans = int(_ifql_index_offsets[ m.cs])

	_klen = int(_ifql_single_lengths[ m.cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case ( m.data)[( m.p)] < _ifql_trans_keys[_mid]:
				_upper = _mid - 1
			case ( m.data)[( m.p)] > _ifql_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_ifql_range_lengths[ m.cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case ( m.data)[( m.p)] < _ifql_trans_keys[_mid]:
				_upper = _mid - 2
			case ( m.data)[( m.p)] > _ifql_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	_trans = int(_ifql_indicies[_trans])
	 m.cs = int(_ifql_trans_targs[_trans])

	if _ifql_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_ifql_trans_actions[_trans])
	_nacts = uint(_ifql_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _ifql_actions[_acts-1] {
		case 0:
//line ifql.rl:16

		m.ts = m.p
		log.Println("start", m.ts)
	
		case 1:
//line ifql.rl:20

		m.te = m.p
		log.Println("finish", m.ts, m.te)
	
		case 2:
//line ifql.rl:25

		log.Println("action identifier", m.stack)

		m.push(&ast.Identifier{
			Name: string(m.current()),
		})
	
		case 3:
//line ifql.rl:33

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
	
		case 4:
//line ifql.rl:46

		log.Println("action expressionStatement", m.stack)

		expr := m.pop().(ast.Expression)
		s := &ast.ExpressionStatement{
			Expression: expr,
		}
		m.push(s)
	
		case 5:
//line ifql.rl:55

		log.Println("action returnStatement", m.stack)

		expr := m.pop().(ast.Expression)
		s := &ast.ReturnStatement{
			Argument: expr,
		}
		m.push(s)
	
		case 6:
//line ifql.rl:65

		log.Println("action beginBlockStatement", m.stack)
		m.push(blockMarker)
	
		case 7:
//line ifql.rl:69

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
	
//line ifql.go:408
		}
	}

_again:
	if  m.cs == 0 {
		goto _out
	}
	( m.p)++
	if ( m.p) != ( m.pe) {
		goto _resume
	}
	_test_eof: {}
	if ( m.p) == ( m.eof) {
		__acts := _ifql_eof_actions[ m.cs]
		__nacts := uint(_ifql_actions[__acts]); __acts++
		for ; __nacts > 0; __nacts-- {
			__acts++
			switch _ifql_actions[__acts-1] {
			case 1:
//line ifql.rl:20

		m.te = m.p
		log.Println("finish", m.ts, m.te)
	
			case 2:
//line ifql.rl:25

		log.Println("action identifier", m.stack)

		m.push(&ast.Identifier{
			Name: string(m.current()),
		})
	
			case 3:
//line ifql.rl:33

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
	
			case 4:
//line ifql.rl:46

		log.Println("action expressionStatement", m.stack)

		expr := m.pop().(ast.Expression)
		s := &ast.ExpressionStatement{
			Expression: expr,
		}
		m.push(s)
	
			case 5:
//line ifql.rl:55

		log.Println("action returnStatement", m.stack)

		expr := m.pop().(ast.Expression)
		s := &ast.ReturnStatement{
			Argument: expr,
		}
		m.push(s)
	
			case 7:
//line ifql.rl:69

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
	
			case 8:
//line ifql.rl:84

		log.Println("action program", m.stack)

		p := &ast.Program{}
		for m.len() > 0 {
			stmt := m.pop().(ast.Statement)
			p.Body = append(p.Body, stmt)
		}
		m.push(p)
	
//line ifql.go:507
			}
		}
	}

	_out: {}
	}

//line ifql.rl:149

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
