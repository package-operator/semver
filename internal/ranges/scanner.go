package ranges

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"pkg.package-operator.run/semver/internal"
)

// An ErrorHandler may be provided to Scanner.Init.
type ErrorHandler func(pos internal.Position, msg string)

// Scanner implements a scanner for systemd unit files.
// It takes a []byte as source which can then be tokenized
// through repeated calls to the Scan method.
type Scanner struct {
	pos internal.Position
	src []byte
	err ErrorHandler

	// scanning state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)

	ErrorCount int // number of errors encountered
}

func (s *Scanner) Init(src []byte, errHandler ErrorHandler) {
	s.src = src
	s.err = errHandler
	s.pos = 0

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.ErrorCount = 0

	s.next()
}

func (s *Scanner) error(msg string) {
	if s.err != nil {
		s.err(s.pos, msg)
	}
	s.ErrorCount++
}

func (s *Scanner) next() {
	if s.rdOffset >= len(s.src) {
		// we are at the end of our buffer
		// -> EOF
		s.offset = len(s.src)
		s.ch = -1
		return
	}

	s.offset = s.rdOffset
	switch s.ch {
	case '\n':
		s.error("illegal character NEWLINE")
	default:
		s.pos++
	}

	r, w := rune(s.src[s.rdOffset]), 1
	switch {
	case r == 0:
		s.error("illegal character NUL")

	case r >= utf8.RuneSelf:
		r, w = utf8.DecodeRune(s.src[s.rdOffset:])
		if r == utf8.RuneError && w == 1 {
			s.error("illegal UTF-8 encoding")
		}
	}
	s.ch = r
	s.rdOffset += w
}

func (s *Scanner) scanNumber() string {
	offs := s.offset - 1
	for unicode.IsDigit(s.ch) && s.ch != -1 {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanSpace() string {
	offs := s.offset - 1
	for s.ch == ' ' && s.ch != -1 {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) followedByEqual(tok0, tok1 Token) Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Scanner) Scan() (pos internal.Position, tok Token, lit uint64) {
	pos = s.pos
	ch := s.ch
	s.next()

	switch ch {
	case -1:
		tok = EOF
	case ' ':
		tok = SPACE
		_ = s.scanSpace()

	case '=':
		tok = EQUAL
	case '!':
		// MUST be followed by =
		if s.ch != '=' {
			s.error(fmt.Sprintf("unexpected character %#U", s.ch))
			tok = ILLEGAL
		} else {
			tok = NOT_EQUAL
			s.next()
		}
	case '>':
		tok = s.followedByEqual(GREATER, GREATER_EQUAL)
	case '<':
		tok = s.followedByEqual(LESS, LESS_EQUAL)

	case '~':
		tok = TILDE
	case '^':
		tok = CARET

	case '.':
		tok = DOT
	case '-':
		tok = HYPHEN
	case 'x', 'X', '*':
		tok = WILDCARD
	case '|':
		if s.ch != '|' {
			s.error(fmt.Sprintf("unexpected character %#U", s.ch))
			tok = ILLEGAL
		} else {
			tok = OR
			s.next()
		}
	case ',':
		tok = AND
	case '&':
		if s.ch != '&' {
			s.error(fmt.Sprintf("unexpected character %#U", s.ch))
			tok = ILLEGAL
		} else {
			tok = AND
			s.next()
		}

	default:
		if ch == '0' && !internal.IsDigit(s.ch) {
			tok = NUMBER
			lit = 0
			return
		}
		if internal.IsDigit(ch) && ch != '0' {
			tok = NUMBER
			num := s.scanNumber()

			var err error
			lit, err = strconv.ParseUint(num, 10, 0)
			if err != nil {
				panic(err)
			}
			return
		}

		s.error(fmt.Sprintf("unexpected character %#U", s.ch))
		tok = ILLEGAL

		if tok == ILLEGAL {
			return
		}
	}
	return
}
