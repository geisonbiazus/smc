package smc

type TokenCollector interface {
	OpenBrace()
	CloseBrace()
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
		l.collector.CloseBrace()
	}
}
