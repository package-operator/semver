package ranges

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.package-operator.run/semver/internal"
)

const example0 = `       4`

var example0tokens = []tokenEntry{
	{pos: 1, tok: SPACE},
	{pos: 8, tok: NUMBER, lit: 4},
	{pos: 8, tok: EOF},
}

const example1 = `>=0.0.3 <0.0.4`

var example1tokens = []tokenEntry{
	{pos: 1, tok: GREATER_EQUAL},
	{pos: 3, tok: NUMBER, lit: 0},
	{pos: 4, tok: DOT},
	{pos: 5, tok: NUMBER, lit: 0},
	{pos: 6, tok: DOT},
	{pos: 7, tok: NUMBER, lit: 3},
	{pos: 8, tok: SPACE},
	{pos: 9, tok: LESS},
	{pos: 10, tok: NUMBER, lit: 0},
	{pos: 11, tok: DOT},
	{pos: 12, tok: NUMBER, lit: 0},
	{pos: 13, tok: DOT},
	{pos: 14, tok: NUMBER, lit: 4},
	{pos: 14, tok: EOF},
}

const example2 = `^0.2.3`

var example2tokens = []tokenEntry{
	{pos: 1, tok: CARET},
	{pos: 2, tok: NUMBER, lit: 0},
	{pos: 3, tok: DOT},
	{pos: 4, tok: NUMBER, lit: 2},
	{pos: 5, tok: DOT},
	{pos: 6, tok: NUMBER, lit: 3},
	{pos: 6, tok: EOF},
}

const example3 = `~1.x`

var example3tokens = []tokenEntry{
	{pos: 1, tok: TILDE},
	{pos: 2, tok: NUMBER, lit: 1},
	{pos: 3, tok: DOT},
	{pos: 4, tok: WILDCARD},
	{pos: 4, tok: EOF},
}

const example4 = `>= 1.2, != 1.4.5`

var example4tokens = []tokenEntry{
	{pos: 1, tok: GREATER_EQUAL},
	{pos: 3, tok: SPACE},
	{pos: 4, tok: NUMBER, lit: 1},
	{pos: 5, tok: DOT},
	{pos: 6, tok: NUMBER, lit: 2},
	{pos: 7, tok: AND},
	{pos: 8, tok: SPACE},
	{pos: 9, tok: NOT_EQUAL},
	{pos: 11, tok: SPACE},
	{pos: 12, tok: NUMBER, lit: 1},
	{pos: 13, tok: DOT},
	{pos: 14, tok: NUMBER, lit: 4},
	{pos: 15, tok: DOT},
	{pos: 16, tok: NUMBER, lit: 5},
	{pos: 16, tok: EOF},
}

const example5 = `=1.0.2||=1.4.5`

var example5tokens = []tokenEntry{
	{tok: EQUAL},
	{tok: NUMBER, lit: 1},
	{tok: DOT},
	{tok: NUMBER, lit: 0},
	{tok: DOT},
	{tok: NUMBER, lit: 2},
	{tok: OR},
	{tok: EQUAL},
	{tok: NUMBER, lit: 1},
	{tok: DOT},
	{tok: NUMBER, lit: 4},
	{tok: DOT},
	{tok: NUMBER, lit: 5},
	{tok: EOF},
}

const example6 = `=11&&<12`

var example6tokens = []tokenEntry{
	{tok: EQUAL},
	{tok: NUMBER, lit: 11},
	{tok: AND},
	{tok: LESS},
	{tok: NUMBER, lit: 12},
	{tok: EOF},
}

type tokenEntry struct {
	pos internal.Position
	tok Token
	lit uint64
}

func TestScanner(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Input        string
		Tokens       []tokenEntry
		TestPosition bool
	}{
		{
			Input:        example0,
			Tokens:       example0tokens,
			TestPosition: true,
		},
		{
			Input:        example1,
			Tokens:       example1tokens,
			TestPosition: true,
		},
		{
			Input:        example2,
			Tokens:       example2tokens,
			TestPosition: true,
		},
		{
			Input:        example3,
			Tokens:       example3tokens,
			TestPosition: true,
		},
		{
			Input:        example4,
			Tokens:       example4tokens,
			TestPosition: true,
		},
		{
			Input:  example5,
			Tokens: example5tokens,
		},
		{
			Input:  example6,
			Tokens: example6tokens,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Input, func(t *testing.T) {
			t.Parallel()
			// Init
			var s Scanner
			s.Init([]byte(test.Input), nil)

			// Scan
			tokens := []tokenEntry{}
			for {
				pos, tok, lit := s.Scan()
				t.Logf("%s\t%s\t%d\n", pos, tok, lit)
				te := tokenEntry{
					tok: tok,
					lit: lit,
				}
				if test.TestPosition {
					te.pos = pos
				}
				tokens = append(tokens, te)
				if tok == EOF {
					break
				}
			}
			assert.Equal(t, test.Tokens, tokens)
		})
	}
}
