//line ifql.rl:1
package parser

import "fmt"

//line ifql.go:9
var _ifql_actions []byte = []byte{
	0, 2, 0, 1,
}

var _ifql_key_offsets []byte = []byte{
	0, 0, 1,
}

var _ifql_trans_keys []byte = []byte{
	119,
}

var _ifql_single_lengths []byte = []byte{
	0, 1, 0,
}

var _ifql_range_lengths []byte = []byte{
	0, 0, 0,
}

var _ifql_index_offsets []byte = []byte{
	0, 0, 2,
}

var _ifql_trans_targs []byte = []byte{
	2, 0, 0,
}

var _ifql_trans_actions []byte = []byte{
	0, 0, 0,
}

var _ifql_eof_actions []byte = []byte{
	0, 0, 1,
}

const ifql_start int = 1
const ifql_first_final int = 2
const ifql_error int = 0

const ifql_en_main int = 1

//line ifql.rl:40

type Machine struct {
	data       []byte
	cs         int
	ts, te     int
	act        int
	p, pe, eof int

	tokens []Token
}

func NewMachine(data []byte) Machine {
	m := Machine{
		data: data,
		pe:   len(data),
		eof:  len(data),
	}

//line ifql.rl:60
//line ifql.rl:61
//line ifql.rl:62
//line ifql.rl:63
//line ifql.rl:64
//line ifql.go:84
	{
		m.cs = ifql_start
	}

//line ifql.rl:65
	return m
}

func (m *Machine) Scan() (Token, error) {
	if m.p >= m.pe {
		return Token{Type: EOF}, fmt.Errorf("EOF")
	}

	//var ts, te int
	var tp TokenType

//line ifql.go:104
	{
		var _klen int
		var _trans int
		var _keys int
		if (m.p) == (m.pe) {
			goto _test_eof
		}
		if m.cs == 0 {
			goto _out
		}
	_resume:
		_keys = int(_ifql_key_offsets[m.cs])
		_trans = int(_ifql_index_offsets[m.cs])

		_klen = int(_ifql_single_lengths[m.cs])
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
				case (m.data)[(m.p)] < _ifql_trans_keys[_mid]:
					_upper = _mid - 1
				case (m.data)[(m.p)] > _ifql_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_ifql_range_lengths[m.cs])
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
				case (m.data)[(m.p)] < _ifql_trans_keys[_mid]:
					_upper = _mid - 2
				case (m.data)[(m.p)] > _ifql_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		m.cs = int(_ifql_trans_targs[_trans])

		if m.cs == 0 {
			goto _out
		}
		(m.p)++
		if (m.p) != (m.pe) {
			goto _resume
		}
	_test_eof:
		{
		}
		if (m.p) == (m.eof) {
			__acts := _ifql_eof_actions[m.cs]
			__nacts := uint(_ifql_actions[__acts])
			__acts++
			for ; __nacts > 0; __nacts-- {
				__acts++
				switch _ifql_actions[__acts-1] {
				case 0:
//line ifql.rl:15
					tp = Identifier

				case 1:
//line ifql.rl:22
					tp = EOF

//line ifql.go:195
				}
			}
		}

	_out:
		{
		}
	}

//line ifql.rl:79
	if tp == NoMatch {
		return Token{Type: EOF}, fmt.Errorf("not found")
	}

	return Token{
		Type: tp,
		Data: m.data[m.ts:m.te],
	}, nil
}
