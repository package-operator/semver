package semver

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"

	"pkg.package-operator.run/semver/internal"
)

// MustNewVersion parses the given string into a Version object and panics on error.
func MustNewVersion(src string) Version {
	v, err := NewVersion(src)
	if err != nil {
		panic(err)
	}
	return v
}

// NewVersion parses the given string into a Version object.
func NewVersion(src string) (Version, error) {
	return parseVersion([]byte(src))
}

// parse bytes into a Version.
func parseVersion(src []byte) (Version, error) {
	var p parser
	p.init(src)
	return p.parse()
}

type parser struct {
	pos internal.Position // column number
	src []byte

	// parser state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)

	// parser position
	// 0 => Major, 1 => Minor, 2 => Patch
	// 3 => Pre Release
	// 4 => Build
	logicalPosition int
}

func (p *parser) init(src []byte) {
	p.src = src
	p.pos = 0

	p.ch = -1
	p.offset = 0
	p.rdOffset = 0

	_ = p.next()
}

func (p *parser) next() error {
	if p.rdOffset >= len(p.src) {
		// we are at the end of our buffer
		// -> EOF
		p.offset = len(p.src)
		p.ch = -1
		return nil
	}

	p.offset = p.rdOffset
	switch p.ch {
	case '\n':
		return fmt.Errorf("%s: illegal character NEWLINE", p.pos)
	case ' ':
		return fmt.Errorf("%s: illegal character SPACE", p.pos)
	default:
		p.pos++
	}

	r, w := rune(p.src[p.rdOffset]), 1
	switch {
	case r == 0:
		return fmt.Errorf("%s: illegal character NUL", p.pos-1)

	case r >= utf8.RuneSelf:
		r, w = utf8.DecodeRune(p.src[p.rdOffset:])
		if r == utf8.RuneError && w == 1 {
			return fmt.Errorf("%s: illegal UTF-8 encoding", p.pos-1)
		}
	}
	p.ch = r
	p.rdOffset += w
	return nil
}

func (p *parser) parse() (Version, error) {
	// TODO add column number to errors
	v := Version{}

parser:
	for {
		// col = p.col
		ch := p.ch
		if err := p.next(); err != nil {
			return v, err
		}

		if ch == -1 {
			break parser
		}

		if p.logicalPosition > 2 {
			// version x.y.z parsed
			switch ch {
			case '-':
				// pre release
				var err error
				v.PreRelease, err = p.scanPreRelease()
				if err != nil {
					return Version{}, err
				}
				p.logicalPosition++
				goto parser

			case '+':
				// build metadata
				var err error
				v.BuildMetadata, err = p.scanBuildMeta()
				if err != nil {
					return Version{}, err
				}
				p.logicalPosition++
				goto parser

			default:
				return Version{}, fmt.Errorf("%s: invalid character %q", p.pos, p.ch)
			}
		}

		num, err := p.scanNumber(ch)
		if err != nil {
			return Version{}, err
		}
		switch p.logicalPosition {
		case 0:
			v.Major = num
			_ = p.scanDot()
			p.logicalPosition++

		case 1:
			v.Minor = num
			_ = p.scanDot()
			p.logicalPosition++

		case 2:
			v.Patch = num
			p.logicalPosition++
		}
	}

	switch p.logicalPosition {
	case 0:
		return v, fmt.Errorf("%s: missing major", p.pos+1)
	case 1:
		return v, fmt.Errorf("%s: missing minor", p.pos+1)
	case 2:
		return v, fmt.Errorf("%s: missing patch", p.pos+1)
	}

	return v, nil
}

func (p *parser) scanDot() error {
	if p.ch != '.' {
		return fmt.Errorf("%s: invalid character %q", p.pos, p.ch)
	}
	if err := p.next(); err != nil {
		return err
	}
	return nil
}

func (p *parser) scanNumber(ch rune) (uint64, error) {
	pos := p.pos
	offs := p.offset - 1
	for isDigit(p.ch) && p.ch != -1 {
		if err := p.next(); err != nil {
			return 0, err
		}
	}
	out := string(p.src[offs:p.offset])
	if len(out) == 0 || out == "." {
		return 0, fmt.Errorf("%s: expected number, got nothing", pos)
	}
	if out != "0" && !isPositiveDigit(rune(out[0])) {
		return 0, fmt.Errorf("%s: starts with non-positive integer %q", p.pos-1, ch)
	}
	return strconv.ParseUint(out, 10, 0)
}

func (p *parser) scanBuildMeta() ([]string, error) {
	var prParts []string
	for {
		pos := p.pos
		s, end, err := p.scanString()
		if err != nil {
			return nil, err
		}
		if len(s) == 0 {
			// must be non-empty
			return nil, fmt.Errorf("%s: build identifier empty", pos)
		}
		if !isBuildIdentifier(s) {
			return nil, fmt.Errorf("%s: invalid build identifier %q", pos, s)
		}
		prParts = append(prParts, s)
		if end {
			break
		}
	}
	return prParts, nil
}

func (p *parser) scanPreRelease() ([]PreReleaseIdentifier, error) {
	if p.logicalPosition != 3 {
		return nil, errors.New("pre release not after patch")
	}

	var prParts []PreReleaseIdentifier
	for {
		pos := p.pos
		s, end, err := p.scanString()
		if err != nil {
			return nil, err
		}
		if len(s) == 0 {
			// must be non-empty
			return nil, fmt.Errorf("%s: pre release identifier empty", pos)
		}
		if !isPreReleaseIdentifier(s) {
			return nil, fmt.Errorf("%s: invalid pre release identifier %q", pos, s)
		}
		prParts = append(prParts, ToPreReleaseIdentifier(s))
		if end {
			break
		}
	}
	return prParts, nil
}

func (p *parser) scanString() (s string, end bool, err error) {
	offs := p.offset
	// scan until build meta, dot or end
	for p.ch != '+' && p.ch != '.' && p.ch != -1 {
		if err = p.next(); err != nil {
			return
		}
	}
	out := string(p.src[offs:p.offset])
	if p.ch == '.' {
		// skip dot
		if err = p.next(); err != nil {
			return
		}
	}
	return out, p.ch == '+' || p.ch == -1, nil
}
