// Package semver implements logic to work with Sementic Versioning 2.0.0 in Go.
// It provides:
// - Parser for semantic versions
// - Validation of semantic versions
// - Sorting of semantic versions
// - Parser for semantic version range constraints
// - Range constraint matching
//   - whether version contained in range
//   - whether a range is contained in another range
package semver

import (
	"fmt"

	"pkg.package-operator.run/semver/internal"
	"pkg.package-operator.run/semver/internal/ranges"
)

const maxUint64 = ^uint64(0)

// MustNewConstraint parses the given string into a Version Constraint or panics.
func MustNewConstraint(data string) Constraint {
	c, err := NewConstraint(data)
	if err != nil {
		panic(err)
	}
	return c
}

// NewConstraint parses the given string into a Version Constraint.
func NewConstraint(data string) (Constraint, error) {
	c, err := parseConstraint([]byte(data))
	if err != nil {
		return nil, err
	}
	return &originalInputConstraint{
		Constraint: c,
		original:   data,
	}, nil
}

// parseConstraint bytes into a Version Constraint.
func parseConstraint(data []byte) (Constraint, error) {
	var p parserState
	p.init(data)
	c, err := p.parse()
	if err != nil {
		return nil, err
	}
	return c, nil
}

type parserState struct {
	scanner ranges.Scanner

	c Constraint

	or       or           // active || combined ranges or And constraints
	and      and          // active && combined ranges
	operator ranges.Token // EQUAL,NOT_EQUAL, GREATER, LESS, GREATER_EQUAL, LESS_EQUAL

	expectingNumber bool     // if we expect a number next
	lastSemverPos   int      // previous position after versionClose()
	semverPos       int      // 0=Major, 1=Minor, 2=Patch
	version         *Version // active version being parsed
	max             bool     // false = min part of the range, true = max part of the range
	activeRange     *Range   // active range being parsed
	errors          []string // scanner errors
}

func (p *parserState) init(src []byte) *parserState {
	p.scanner.Init(src, func(pos internal.Position, msg string) {
		p.errors = append(p.errors, fmt.Sprintf("%s: %s", pos, msg))
	})
	p.resetRange()
	return p
}

func (p *parserState) addNumberToVersion(num uint64) {
	if p.activeRange == nil {
		p.activeRange = &Range{}
	}
	if p.max {
		p.version = &p.activeRange.Max
	} else {
		p.version = &p.activeRange.Min
	}

	switch p.semverPos {
	case 0:
		p.version.Major = num
	case 1:
		p.version.Minor = num
	case 2:
		p.version.Patch = num
	}
}

func (p *parserState) closeVersion(pos internal.Position) error {
	if p.version == nil {
		return nil
	}
	if p.expectingNumber {
		// semver clause incomplete!
		return fmt.Errorf("%s: semver clause incomplete", pos)
	}

	p.lastSemverPos = p.semverPos
	p.semverPos = 0
	if !p.max {
		// move to max part of range next.
		p.max = true
	}
	p.version = nil // no active version
	return nil
}

func (p *parserState) closeRange(pos internal.Position) error {
	if p.activeRange == nil {
		return nil
	}
	if err := p.closeVersion(pos); err != nil {
		return err
	}
	r := p.activeRange

	switch p.operator {
	case ranges.EQUAL, ranges.NOT_EQUAL:
		r.Max = r.Min
		switch p.lastSemverPos {
		case 0:
			r.Max.Minor = maxUint64
			r.Max.Patch = maxUint64

		case 1:
			r.Max.Patch = maxUint64
		}

	case ranges.GREATER:
		switch p.lastSemverPos {
		// 1.x.x -> 2.x.x
		case 0:
			r.Min.Major++

		// 1.2.x -> 1.3.x
		case 1:
			r.Min.Minor++

		// 1.2.0 -> 1.2.1
		default:
			r.Min.Patch++
		}
		// x.x.x
		r.Max = Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64}

	case ranges.HYPHEN:
		// x.0 => x.x
		if r.Max.Major == maxUint64 {
			r.Max.Minor = maxUint64
		}
		// 1.x.0 => 1.x.x
		if r.Max.Minor == maxUint64 {
			r.Max.Patch = maxUint64
		}

	case ranges.LESS:
		r.Max = r.Min
		switch {
		// 1.2.0 => 1.1.x
		case r.Max.Patch == 0 && r.Max.Minor > 0:
			r.Max.Patch = maxUint64
			r.Max.Minor--
		// 1.0.0 => 0.x.x
		case r.Max.Minor == 0:
			r.Max.Patch = maxUint64
			r.Max.Minor = maxUint64
			r.Max.Major--
		// 1.2.3 => 1.2.2
		default:
			r.Max.Patch--
		}
		r.Min = Version{} // 0.0.0

	case ranges.LESS_EQUAL:
		r.Max = r.Min
		r.Min = Version{} // 0.0.0

	case ranges.GREATER_EQUAL:
		r.Max = Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64}

	case ranges.TILDE:
		r.Max = r.Min
		r.Max.Patch = maxUint64
		if r.Max.Minor == 0 {
			r.Max.Minor = maxUint64
		}

	case ranges.CARET:
		r.Max = r.Min
		if r.Min.Major != 0 {
			r.Max.Minor = maxUint64
		}
		r.Max.Patch = maxUint64

	default:
		return fmt.Errorf("%s: range closed without operator", pos)
	}

	var c Constraint
	c = r
	// negate result
	if p.operator == ranges.NOT_EQUAL {
		c = not{Range: *r}
	}

	p.and = append(p.and, c)
	var err error
	p.and, err = compactAndValidateLogicalAND(pos, p.and)
	if err != nil {
		return err
	}

	// reset
	p.resetRange()
	return nil
}

func (p *parserState) resetRange() {
	p.max = false
	p.activeRange = nil
	p.operator = 0
}

func (p *parserState) close(pos internal.Position) error {
	if err := p.closeRange(pos); err != nil {
		return err
	}

	switch len(p.and) {
	case 0:
	case 1:
		p.or = append(p.or, p.and[0])
	default:
		p.or = append(p.or, p.and)
	}

	switch len(p.or) {
	case 0:
	case 1:
		p.c = p.or[0]
	default:
		p.c = p.or
	}
	return nil
}

func (p *parserState) parse() (Constraint, error) {
parse:
	for {
		pos, tok, lit := p.scanner.Scan()
		if len(p.errors) > 0 {
			return nil, fmt.Errorf("%s", p.errors[0])
		}

		switch tok {
		case ranges.ILLEGAL:
			goto parse

		case ranges.SPACE:
			if err := p.closeVersion(pos); err != nil {
				return nil, err
			}

		case ranges.EQUAL, ranges.NOT_EQUAL,
			ranges.GREATER, ranges.GREATER_EQUAL,
			ranges.LESS, ranges.LESS_EQUAL,
			ranges.TILDE, ranges.CARET:
			if err := p.closeRange(pos); err != nil {
				return nil, err
			}
			p.operator = tok

		case ranges.AND:
			if p.activeRange == nil {
				return nil, fmt.Errorf("%s: AND empty range constraint", pos)
			}
			if err := p.closeRange(pos); err != nil {
				return nil, fmt.Errorf("%s: AND %w", pos, err)
			}

		case ranges.OR:
			if p.activeRange == nil {
				return nil, fmt.Errorf("%s: OR empty range constraint", pos)
			}
			if err := p.closeRange(pos); err != nil {
				return nil, fmt.Errorf("%s: OR %w", pos, err)
			}
			// Shift current AND constraint into OR
			if len(p.and) == 1 {
				p.or = append(p.or, p.and[0])
			} else {
				p.or = append(p.or, p.and)
			}
			p.and = nil

		case ranges.HYPHEN:
			if p.operator == ranges.HYPHEN {
				// we are already within a HYPON range.
				// seeing a HYPON again is an error.
				return nil, fmt.Errorf(`%s: double hyphen in range constraint`, pos)
			}
			if err := p.closeVersion(pos); err != nil {
				return nil, err
			}
			p.operator = tok
			p.max = true
			p.expectingNumber = true

		case ranges.EOF:
			if err := p.close(pos); err != nil {
				return nil, err
			}
			break parse

		case ranges.NUMBER:
			p.addNumberToVersion(lit)
			p.expectingNumber = false

		case ranges.WILDCARD:
			if p.max {
				p.addNumberToVersion(maxUint64)
			} else {
				p.addNumberToVersion(0)
			}
			if p.semverPos != 0 {
				p.semverPos--
			}
			p.expectingNumber = false

		case ranges.DOT:
			p.expectingNumber = true
			p.semverPos++
			if p.semverPos > 2 {
				return nil, fmt.Errorf("%s: found 3rd dot when parsing semver", pos)
			}
		}
	}
	if p.c == nil {
		return nil, fmt.Errorf("%s: empty", internal.Position(1))
	}

	return p.c, nil
}

func compactAndValidateLogicalAND(pos internal.Position, and and) (and, error) {
	if len(and) < 2 {
		return and, nil
	}

	var fullyDefinedConstraints []Constraint

	// find min version and max version
	var (
		maxVersion *Version
		minVersion *Version
	)
	for _, c := range and {
		r, ok := c.(*Range)
		switch {
		case ok && isMinUnconstraint(*r):
			if maxVersion != nil {
				return nil, fmt.Errorf(
					"%s: <=%s overlaps with <=%s in logical AND",
					pos, r.Max.String(), maxVersion,
				)
			}
			maxVersion = &r.Max

		case ok && isMaxUnconstraint(*r):
			if minVersion != nil {
				return nil, fmt.Errorf(
					"%s: >=%s overlaps with >=%s in logical AND",
					pos, r.Min.String(), minVersion,
				)
			}
			minVersion = &r.Min

		case ok:
			if minVersion != nil && maxVersion != nil {
				existingRange := Range{Min: *minVersion, Max: *maxVersion}
				if !existingRange.Contains(r) {
					return nil, fmt.Errorf(
						"%s: non overlapping ranges %q and %q in logical AND",
						pos, r.String(), existingRange.String(),
					)
				}
				fullyDefinedConstraints = append(fullyDefinedConstraints, c)
			} else {
				minVersion = &r.Min
				maxVersion = &r.Max
			}

		default:
			fullyDefinedConstraints = append(fullyDefinedConstraints, c)
		}
	}
	if minVersion != nil && maxVersion != nil {
		fullyDefinedConstraints = append(fullyDefinedConstraints, &Range{
			Min: *minVersion,
			Max: *maxVersion,
		})
	}
	return fullyDefinedConstraints, nil
}

func isMinUnconstraint(r Range) bool {
	return r.Min.Same(Version{})
}

func isMaxUnconstraint(r Range) bool {
	return r.Max.Same(Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64})
}
