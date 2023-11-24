package semver

import (
	"fmt"
)

// Range represents a min to max version range.
// Min and Max are inclusive >=/<=
// e.g. 1.0.0 - 1.1.0 would contain both 1.0.0 and 1.1.0.
type Range struct {
	Min Version
	Max Version
}

func (r *Range) String() string {
	return fmt.Sprintf("%s - %s", r.Min.String(), r.Max.String())
}

// Check if the given version is contained in the range.
func (r *Range) Check(v Version) bool {
	if v.LessThan(r.Min) {
		return false
	}
	if v.GreaterThan(r.Max) {
		return false
	}
	return true
}

// Checks if the given constraint fits into this range.
func (r *Range) Contains(other Constraint) bool {
	return rangeContains(*r, other)
}

// reverse contains comparison to compare a range against another constraint.
func rangeContains(r Range, other Constraint) bool {
	switch v := other.(type) {
	case *Range:
		return rangeContainsRange(r, *v)

	case not:
		return !rangeContainsRange(r, v.Range)

	case and:
		for _, ac := range v {
			if !rangeContains(r, ac) {
				return false
			}
		}
		return true

	case or:
		for _, ac := range v {
			if rangeContains(r, ac) {
				return true
			}
		}
		return false
	}
	return false
}

func rangeContainsRange(rA, rB Range) bool {
	if rA.Min.Compare(rB.Min) <= 0 &&
		rA.Max.Compare(rB.Max) >= 0 {
		return true
	}
	return false
}
