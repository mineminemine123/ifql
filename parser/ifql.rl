package parser

import "fmt"

%%{
	machine ifql;

	action start {
		fmt.Printf("D! begin p: %d pe: %d eof: %d\n", m.p, m.pe, m.eof)
	}

	action operator {
        tp = Operator
	}
	action identifier {
        tp = Identifier
	}
	action expression {
        tp = Expression
	}

	action end {
		tp = EOF;
	}



	ws = [\t\v\f ];

    operator = '=' %operator;
    identifier = [\w] %identifier;
    expression = 'expression' %expression;
    

    program = identifier ws operator ws expression;

    main := identifier %eof(end);
    
	write data;
}%%

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

func (m *Machine) Scan() (Token, error) {
	if m.p >= m.pe {
		return Token{Type:EOF}, fmt.Errorf("EOF")
	}

	//var ts, te int
	var tp TokenType


	%% write exec;

	if tp == NoMatch {
		return Token{Type:EOF}, fmt.Errorf("not found")
	}

	return Token{
        Type: tp,
        Data: m.data[m.ts:m.te],
    } , nil
}
