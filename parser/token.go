package parser

type TokenType int

const (
	NoMatch TokenType = iota
	Identifier
	Operator
	Expression
	EOF
)

type Token struct {
	Type TokenType
	Data []byte
}
