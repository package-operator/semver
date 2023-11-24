package ranges

import (
	"strconv"
)

// Token is the type for lexical tokens of a systemd configuration file.
type Token int

// List of Tokens.
//
//nolint:revive,stylecheck // naming is fine here
const (
	// Special.
	ILLEGAL Token = iota
	EOF
	SPACE

	// Values - Essentially everything that does not fit elsewhere.
	NUMBER

	// Operators and delimiters.
	EQUAL         // =
	NOT_EQUAL     // !=
	GREATER       // >
	GREATER_EQUAL // >=
	LESS          // <
	LESS_EQUAL    // <=

	TILDE // ~
	CARET // ^

	DOT      // .
	HYPHEN   // -
	WILDCARD // x, X or *
	OR       // ||
	AND      // , or &&
)

var tokens = [...]string{
	// Special
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	SPACE:   "SPACE",

	// Values
	NUMBER: "NUMBER",

	// Operators and delimiters
	EQUAL:         "EQUAL",
	NOT_EQUAL:     "NOT_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",

	TILDE: "TILDE",
	CARET: "CARET",

	DOT:      "DOT",
	HYPHEN:   "HYPHEN",
	WILDCARD: "WILDCARD",
	OR:       "OR",
	AND:      "AND",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}
