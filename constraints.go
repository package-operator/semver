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

//nolint:revive // Receiver name differs from type to avoid shadowing
func (a and) Contains(other Constraint) bool {
	// Unwrap originalInputConstraint if needed
	otherUnwrapped := other
	if oic, ok := other.(*originalInputConstraint); ok {
		otherUnwrapped = oic.Constraint
	}

	// Special case: when checking if (A && B) contains (C && D)
	otherAnd, ok := otherUnwrapped.(and)
	if !ok {
		// For non-AND constraints, check if all of our constraints contain it
		for _, r := range a {
			if !r.Contains(other) {
				return false
			}
		}
		return true
	}

	// For each constraint in other AND, check if it's covered by our AND
	return a.containsAnd(otherAnd)
}

//nolint:revive // Receiver name differs from type to avoid shadowing
func (a and) containsAnd(other and) bool {
	for _, otherConstraint := range other {
		if !a.coversConstraint(otherConstraint) {
			return false
		}
	}
	return true
}

//nolint:revive // Receiver name differs from type to avoid shadowing
func (a and) coversConstraint(c Constraint) bool {
	// Check if at least one constraint in our AND contains it
	for _, ourConstraint := range a {
		if ourConstraint.Contains(c) {
			return true
		}
	}
	// Check if it's the same constraint (e.g., both have !=2)
	for _, ourConstraint := range a {
		if ourConstraint.String() == c.String() {
			return true
		}
	}
	return false
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

//nolint:revive // Receiver name differs from type to avoid shadowing
func (o or) Contains(other Constraint) bool {
	// Unwrap originalInputConstraint if needed
	otherUnwrapped := other
	if oic, ok := other.(*originalInputConstraint); ok {
		otherUnwrapped = oic.Constraint
	}

	// Special case: when checking if (A || B) contains (C || D),
	// we need to verify that each branch of other is contained by at least one branch of our OR.
	// This is because (A || B) represents the union of A and B, and it contains (C || D)
	// if and only if every version in (C || D) is in (A || B).
	if otherOr, ok := otherUnwrapped.(or); ok {
		// For each branch in other OR, check if at least one branch in our OR contains it
		for _, otherBranch := range otherOr {
			contained := false
			for _, branch := range o {
				if branch.Contains(otherBranch) {
					contained = true
					break
				}
			}
			if !contained {
				return false
			}
		}
		return true
	}

	// For non-OR constraints, check if any of our branches contains it
	for _, r := range o {
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
	return !other.Contains(&not.Range)
}

func (not not) String() string {
	return "!=" + not.Range.String()
}
