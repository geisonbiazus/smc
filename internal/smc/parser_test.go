package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestParser(t *testing.T) {
	assertParserResult(t,
		"a:b c:d {}",
		FSMSyntax{
			Headers: []Header{
				{Name: "a", Value: "b"},
				{Name: "c", Value: "d"},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b:c:d {}",
		FSMSyntax{
			Headers: []Header{
				{Name: "a", Value: "b"},
				{Name: "c", Value: "d"},
			},
			Errors: []SyntaxError{
				{Type: ErrorHeader, LineNumber: 1, Position: 4, Msg: "HEADER|COLON"},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b{c d e f}",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{StateSpec{Name: "c"}, []SubTransition{{"d", "e", []string{"f"}}}},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b{c d e {f g} \n h i j k}",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{StateSpec{Name: "c"}, []SubTransition{{"d", "e", []string{"f", "g"}}}},
				{StateSpec{Name: "h"}, []SubTransition{{"i", "j", []string{"k"}}}},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b { c d e - \n f g h i }",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{StateSpec{Name: "c"}, []SubTransition{{"d", "e", []string{}}}},
				{StateSpec{Name: "f"}, []SubTransition{{"g", "h", []string{"i"}}}},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b { c d - e }",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{StateSpec{Name: "c"}, []SubTransition{{"d", "", []string{"e"}}}},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b { c - d e }",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{StateSpec{Name: "c"}, []SubTransition{{"", "d", []string{"e"}}}},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b { c { d e f \n g h i }}",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{
					StateSpec{Name: "c"}, []SubTransition{
						{"d", "e", []string{"f"}},
						{"g", "h", []string{"i"}},
					},
				},
			},
			Done: true,
		})

	assertParserResult(t,
		"a:b { c { - - - } g { h i j } }",
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{StateSpec{Name: "c"}, []SubTransition{{"", "", []string{}}}},
				{StateSpec{Name: "g"}, []SubTransition{{"h", "i", []string{"j"}}}},
			},
			Done: true,
		})

	assertParserResult(t,
		`a:b {
			c {
				d e {f g}
				h i j
			}
		}`,
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{
					StateSpec{Name: "c"}, []SubTransition{
						{"d", "e", []string{"f", "g"}},
						{"h", "i", []string{"j"}},
					},
				},
			},
			Done: true,
		})

	assertParserResult(t,
		`a:b {
			(c) d e f
			(g) h i -
			j : c : g - - -
		}`,
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{
					StateSpec{Name: "c", AbstractState: true}, []SubTransition{
						{"d", "e", []string{"f"}},
					},
				},
				{
					StateSpec{Name: "g", AbstractState: true}, []SubTransition{
						{"h", "i", []string{}},
					},
				},
				{
					StateSpec{Name: "j", SuperStates: []string{"c", "g"}}, []SubTransition{
						{"", "", []string{}},
					},
				},
			},
			Done: true,
		})

	assertParserResult(t,
		`a:b {
			c <d <e >f >g h i j
		}`,
		FSMSyntax{
			Headers: []Header{{Name: "a", Value: "b"}},
			Logic: []Transition{
				{
					StateSpec{
						Name:         "c",
						EntryActions: []string{"d", "e"},
						ExitActions:  []string{"f", "g"},
					}, []SubTransition{
						{"h", "i", []string{"j"}},
					},
				},
			},
			Done: true,
		})
}

func assertParserResult(t *testing.T, input string, expected FSMSyntax) {
	t.Helper()
	builder := NewSyntaxBuilder()
	parser := NewParser(builder)
	lexer := NewLexer(parser)

	lexer.Lex(bytes.NewBufferString(input))
	assert.DeepEqual(t, expected, builder.FSM())
}
