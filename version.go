package semver

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// VersionList keeps a list of versions and provides helper functions on them.
type VersionList []Version

// StringList converts all Versions into strings and returns them.
func (l VersionList) StringList() []string {
	vs := make([]string, len(l))
	for i := range l {
		vs[i] = l[i].String()
	}
	return vs
}

// Prints the VersionList as comma and space ", " separated list.
func (l VersionList) String() string {
	return strings.Join(l.StringList(), ", ")
}

// Version represents a Semantic Versioning 2.0.0 Version.
type Version struct {
	Major, Minor, Patch uint64
	PreRelease          PreReleaseIdentifierList
	BuildMetadata       []string
}

// Same returns true if both Versions are the same.
func (v *Version) Same(o Version) bool {
	return v.Major == o.Major &&
		v.Minor == o.Minor &&
		v.Patch == o.Patch &&
		slices.Equal(v.PreRelease, o.PreRelease) &&
		slices.Equal(v.BuildMetadata, o.BuildMetadata)
}

// String returns a string representation of the Version.
func (v *Version) String() string {
	s := fmt.Sprintf("%s.%s.%s",
		printXonMaxInt(v.Major),
		printXonMaxInt(v.Minor),
		printXonMaxInt(v.Patch))
	if len(v.PreRelease) > 0 {
		s += "-" + v.PreRelease.String()
	}
	if len(v.BuildMetadata) > 0 {
		s += "+" + strings.Join(v.BuildMetadata, ".")
	}
	return s
}

func printXonMaxInt(d uint64) string {
	if d == maxUint64 {
		return "x"
	}
	return strconv.FormatUint(d, 10)
}

// Equal tests if both version are equal.
func (v *Version) Equal(o Version) bool {
	return v.Compare(o) == 0
}

// LessThan tests if one version is less than another one.
func (v *Version) LessThan(o Version) bool {
	return v.Compare(o) < 0
}

// GreaterThan tests if one version is greater than another one.
func (v *Version) GreaterThan(o Version) bool {
	return v.Compare(o) > 0
}

// Compare compares this version to another one. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version.
func (v *Version) Compare(o Version) int {
	if d := compareSegment(v.Major, o.Major); d != 0 {
		return d
	}
	if d := compareSegment(v.Minor, o.Minor); d != 0 {
		return d
	}
	if d := compareSegment(v.Patch, o.Patch); d != 0 {
		return d
	}
	return o.PreRelease.Compare(v.PreRelease)
}

// PreReleaseIdentifierList is a list of PreReleaseIdentifiers.
// A pre-release version MAY be denoted by appending a hyphen and
// a series of dot separated identifiers immediately following the patch version.
// Examples: 1.0.0-alpha, 1.0.0-alpha.1, 1.0.0-0.3.7, 1.0.0-x.7.z.92, 1.0.0-x-y-z.--.
type PreReleaseIdentifierList []PreReleaseIdentifier

// String returns the pre-release identifier list as a single string joined by '.'.
func (l PreReleaseIdentifierList) String() string {
	pre := make([]string, len(l))
	for i := range l {
		pre[i] = l[i].String()
	}
	return strings.Join(pre, ".")
}

// Compare compares this pre release identifier list to another one.
// It returns -1, 0, or 1 if the version smaller, equal, or larger than the other list.
func (l PreReleaseIdentifierList) Compare(o []PreReleaseIdentifier) int {
	preLen := len(l)
	otherLen := len(o)
	if preLen == 0 && otherLen == 0 {
		return 0
	}
	if preLen == 0 {
		return -1
	}
	if otherLen == 0 {
		return 1
	}

	prel := preLen
	if otherLen > preLen {
		prel = otherLen
	}

	for i := range prel {
		var (
			pre   PreReleaseIdentifier
			other PreReleaseIdentifier
		)
		if i < preLen {
			pre = l[i]
		}
		if i < otherLen {
			other = o[i]
		}
		if d := pre.Compare(other); d != 0 {
			return d
		}
	}
	return 0
}

// PreReleaseIdentifier can be alphanumeric or a number.
type PreReleaseIdentifier struct {
	str string
	num uint64
}

// Compare compares this pre release identifier to another one.
// It returns -1, 0, or 1 if the version smaller, equal, or larger than the other identifier.
func (s *PreReleaseIdentifier) Compare(o PreReleaseIdentifier) int {
	aNum, isANum := s.GetNumber()
	bNum, isBNum := o.GetNumber()

	switch {
	case !isANum && !isBNum:
		if s.str == o.str {
			return 0
		}
		if s.str > o.str {
			return -1
		}
		return 1

	case !isANum:
		// Numeric identifiers always have lower precedence
		// than non-numeric identifiers.
		return -1
	case !isBNum:
		return 1
	}

	if aNum == bNum {
		return 0
	}
	if aNum > bNum {
		return -1
	}
	return 1
}

// Interface returns either a string or uint64 depending on the underlying type.
func (s *PreReleaseIdentifier) Interface() any {
	if len(s.str) > 0 {
		return s.str
	}
	return s.num
}

// GetString returns a string and whether the underlying type is a string.
func (s *PreReleaseIdentifier) GetString() (string, bool) {
	return s.str, len(s.str) > 0
}

// GetNumber returns a uint64 and whether the underlying type is a uint64.
func (s *PreReleaseIdentifier) GetNumber() (uint64, bool) {
	return s.num, len(s.str) == 0
}

// String always returns a string representation.
func (s *PreReleaseIdentifier) String() string {
	if len(s.str) > 0 {
		return s.str
	}
	return strconv.FormatUint(s.num, 10)
}

// ToPreReleaseIdentifier converts the given string into a PreReleaseIdentifier.
func ToPreReleaseIdentifier(s string) PreReleaseIdentifier {
	num, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return PreReleaseIdentifier{str: s}
	}
	return PreReleaseIdentifier{num: num}
}
