package smc

import (
	"regexp"
)

type TokenCollector interface {
	OpenBrace(line, pos int)
	ClosedBrace(line, pos int)
	Colon(line, pos int)
	OpenParen(line, pos int)
	ClosedParen(line, pos int)
	OpenAngle(line, pos int)
	ClosedAngle(line, pos int)
	Dash(line, pos int)
	Name(name string, line, pos int)
	Error(line, pos int)
}

type Lexer struct {
	collector TokenCollector
	pos       int
}

func NewLexer(collector TokenCollector) *Lexer {
	return &Lexer{collector: collector}
}

func (l *Lexer) Lex(input string) {
	l.pos = 0
	for l.pos < len(input) {
		l.lexLine(input)
	}
}

func (l *Lexer) lexLine(input string) {
	if !l.findToken(input) {
		l.addError(input)
	}
}

func (l *Lexer) findToken(input string) bool {
	return l.ignorePossibleWhitespace(input) ||
		l.findSingleCharToken(input) ||
		l.findName(input)
}

func (l *Lexer) findSingleCharToken(input string) bool {
	token := input[l.pos : l.pos+1]
	switch token {
	case "{":
		l.collector.OpenBrace(1, l.pos+1)
	case "}":
		l.collector.ClosedBrace(1, l.pos+1)
	case ":":
		l.collector.Colon(1, l.pos+1)
	case "(":
		l.collector.OpenParen(1, l.pos+1)
	case ")":
		l.collector.ClosedParen(1, l.pos+1)
	case "<":
		l.collector.OpenAngle(1, l.pos+1)
	case ">":
		l.collector.ClosedAngle(1, l.pos+1)
	case "-":
		l.collector.Dash(1, l.pos+1)
	default:
		return false
	}
	l.pos++
	return true
}

var whitespaceRegex = regexp.MustCompile("^\\s+")

func (l *Lexer) ignorePossibleWhitespace(input string) bool {
	if whitespaces, ok := l.matchRegexp(input, whitespaceRegex); ok {
		l.pos += len(whitespaces)
		return true
	}
	return false
}

var nameRegex = regexp.MustCompile("^\\w+")

func (l *Lexer) findName(input string) bool {
	if name, ok := l.matchRegexp(input, nameRegex); ok {
		l.collector.Name(name, 1, l.pos+1)
		l.pos += len(name)
		return true
	}
	return false
}

func (l *Lexer) matchRegexp(input string, r *regexp.Regexp) (token string, ok bool) {
	sub := input[l.pos:]
	match := r.FindStringSubmatch(sub)
	if len(match) == 0 {
		return "", false
	}
	return match[0], true
}

func (l *Lexer) addError(input string) {
	l.collector.Error(1, l.pos+1)
	l.pos++
}
