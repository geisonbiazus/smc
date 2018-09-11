package smc

import (
	"bufio"
	"io"
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
	line      int
}

func NewLexer(collector TokenCollector) *Lexer {
	return &Lexer{collector: collector}
}

func (l *Lexer) Lex(input io.Reader) {
	l.line = 0
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		l.line++
		l.lexLine(scanner.Text())
	}

}

func (l *Lexer) lexLine(input string) {
	l.pos = 1
	for l.pos <= len(input) {
		if !l.findToken(input) {
			l.addError(input)
		}
	}
}

func (l *Lexer) findToken(input string) bool {
	return l.ignorePossibleWhitespace(input) ||
		l.findSingleCharToken(input) ||
		l.findName(input)
}

func (l *Lexer) findSingleCharToken(input string) bool {
	token := input[l.pos-1 : l.pos]
	switch token {
	case "{":
		l.collector.OpenBrace(l.line, l.pos)
	case "}":
		l.collector.ClosedBrace(l.line, l.pos)
	case ":":
		l.collector.Colon(l.line, l.pos)
	case "(":
		l.collector.OpenParen(l.line, l.pos)
	case ")":
		l.collector.ClosedParen(l.line, l.pos)
	case "<":
		l.collector.OpenAngle(l.line, l.pos)
	case ">":
		l.collector.ClosedAngle(l.line, l.pos)
	case "-":
		l.collector.Dash(l.line, l.pos)
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
		l.collector.Name(name, l.line, l.pos)
		l.pos += len(name)
		return true
	}
	return false
}

func (l *Lexer) matchRegexp(input string, r *regexp.Regexp) (token string, ok bool) {
	sub := input[l.pos-1:]
	match := r.FindStringSubmatch(sub)
	if len(match) == 0 {
		return "", false
	}
	return match[0], true
}

func (l *Lexer) addError(input string) {
	l.collector.Error(l.line, l.pos)
	l.pos++
}
