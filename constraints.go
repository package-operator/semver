package semver

import "strings"

// Constraint interface is common to all version constraints.
type Constraint interface {
	// Check if the version is allowed by the constraint or not.
	Check(v Version) bool
	// Check if a range is contained within a constraint.
	Contains(or Constraint) bool
	// Returns the string representation of this constraint.
	String() string
}

// safes original parser input to return via the String method.
type originalInputConstraint struct {
	Constraint
	original string
}

var _ Constraint = (*originalInputConstraint)(nil)

func (oi *originalInputConstraint) String() string {
	return oi.original
}

// and is a list of Ranges that all have to pass.
type and []Constraint

var _ Constraint = and{}

func (and and) Check(v Version) bool {
	for _, r := range and {
		if !r.Check(v) {
			return false
		}
	}
	return true
}

func (and and) Contains(or Constraint) bool {
	for _, r := range and {
		if !r.Contains(or) {
			return false
		}
	}
	return true
}

func (and and) String() string {
	parts := make([]string, len(and))
	for i := range and {
		parts[i] = and[i].String()
	}
	return strings.Join(parts, " && ")
}

// or is a list of Constraints that need at least one match.
type or []Constraint

var _ Constraint = or{}

func (or or) Check(v Version) bool {
	for _, r := range or {
		if r.Check(v) {
			return true
		}
	}
	return false
}

func (or or) Contains(other Constraint) bool {
	for _, r := range or {
		if r.Contains(other) {
			return true
		}
	}
	return false
}

func (or or) String() string {
	parts := make([]string, len(or))
	for i := range or {
		parts[i] = or[i].String()
	}
	return strings.Join(parts, " || ")
}

// not negates the given Range.
type not struct{ Range }

var _ Constraint = not{}

func (not not) Check(v Version) bool {
	return !not.Range.Check(v)
}

func (not not) Contains(other Constraint) bool {
	return !not.Range.Contains(other)
}

func (not not) String() string {
	return "!=" + not.Range.String()
}
