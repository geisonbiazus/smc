package smc

type TokenCollector interface {
	OpenBrace(line, pos int)
	ClosedBrace(line, pos int)
	Colon(line, pos int)
	OpenParen(line, pos int)
	ClosedParen(line, pos int)
	OpenAngle(line, pos int)
	ClosedAngle(line, pos int)
	Dash(line, pos int)
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
		l.collector.OpenBrace(1, 1)
	case "}":
		l.collector.ClosedBrace(1, 1)
	case ":":
		l.collector.Colon(1, 1)
	case "(":
		l.collector.OpenParen(1, 1)
	case ")":
		l.collector.ClosedParen(1, 1)
	case "<":
		l.collector.OpenAngle(1, 1)
	case ">":
		l.collector.ClosedAngle(1, 1)
	case "-":
		l.collector.Dash(1, 1)
	}
}
