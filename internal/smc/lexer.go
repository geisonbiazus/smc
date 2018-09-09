package smc

type TokenCollector interface {
	OpenBrace()
	ClosedBrace()
	Colon()
	OpenParen()
	ClosedParen()
	OpenAngle()
	ClosedAngle()
	Dash()
}

type Lexer struct {
	collector TokenCollector
}

func NewLexer(collector TokenCollector) *Lexer {
	return &Lexer{collector: collector}
}

func (l *Lexer) Lex(input string) {
	switch input {
	case "{":
		l.collector.OpenBrace()
	case "}":
		l.collector.ClosedBrace()
	case ":":
		l.collector.Colon()
	case "(":
		l.collector.OpenParen()
	case ")":
		l.collector.ClosedParen()
	case "<":
		l.collector.OpenAngle()
	case ">":
		l.collector.ClosedAngle()
	case "-":
		l.collector.Dash()
	}
}
