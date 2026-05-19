package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Test(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		r        Range
		v        Version
		expected bool
	}{
		{
			name:     "1.0.0 - 2.0.0 contains 1.0.0",
			r:        Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			v:        MustNewVersion("1.0.0"),
			expected: true,
		},
		{
			name:     "1.0.0 - 2.0.0 contains 2.0.0",
			r:        Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			v:        MustNewVersion("2.0.0"),
			expected: true,
		},
		{
			name:     "1.0.0 - 2.0.0 does not contain 0.1.0",
			r:        Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			v:        MustNewVersion("0.1.0"),
			expected: false,
		},
		{
			name:     "1.0.0 - 2.0.0 does not contain 2.1.0",
			r:        Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			v:        MustNewVersion("2.1.0"),
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := test.r.Check(test.v)
			assert.Equal(t, test.expected, r)
		})
	}
}

//nolint:maintidx // Table-driven test with many edge cases
func TestRange_Contains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		rA       Constraint
		rB       Constraint
		expected bool
	}{
		{
			name:     "1.0.0 - 2.0.0 contains 1.0.0 - 2.0.0",
			rA:       &Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			rB:       &Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			expected: true,
		},
		{
			name:     ">=3 contains >=3",
			rA:       MustNewConstraint(">=3"),
			rB:       MustNewConstraint(">=3"),
			expected: true,
		},
		{
			name:     "1.0.0 - 2.0.0 contains 1.3.0 - 1.4.0",
			rA:       &Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			rB:       &Range{Min: MustNewVersion("1.3.0"), Max: MustNewVersion("1.4.0")},
			expected: true,
		},
		{
			name:     "1.0.0 - 2.0.0 does not contain 0.1.0 - 1.0.0",
			rA:       &Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			rB:       &Range{Min: MustNewVersion("0.1.0"), Max: MustNewVersion("1.0.0")},
			expected: false,
		},
		{
			name:     "1.0.0 - 2.0.0 does not contain 2.1.0 - 3.0.0",
			rA:       &Range{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			rB:       &Range{Min: MustNewVersion("2.1.0"), Max: MustNewVersion("3.0.0")},
			expected: false,
		},
		{
			name:     "1 - 2 || 0 - 1 does not contain 2.1.0 - 3.0.0",
			rA:       MustNewConstraint("1-2 || 0-1"),
			rB:       &Range{Min: MustNewVersion("2.1.0"), Max: MustNewVersion("3.0.0")},
			expected: false,
		},
		{
			name:     "2.1.0 - 3.0.0 does not contain 1 - 2 || 0 - 1",
			rA:       &Range{Min: MustNewVersion("2.1.0"), Max: MustNewVersion("3.0.0")},
			rB:       MustNewConstraint("1-2 || 0-1"),
			expected: false,
		},
		{
			name:     "2.0.0 - 3.0.0 does contain 2-2.4 && 2.1-2.2",
			rA:       &Range{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
			rB:       MustNewConstraint("2-2.4 && 2.1-2.2"),
			expected: true,
		},
		{
			name:     "2.0.0 - 3.0.0 does contain 2-2.4 || 2.1-2.2 && !=5.0.0",
			rA:       &Range{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
			rB:       MustNewConstraint("2-2.4 || 2.1-2.2 && !=5.0.0"),
			expected: true,
		},
		{
			name:     "2.0.0 - 3.0.0 does not contain 2-2.4 && !=2.0.0 || 2.1-2.2 && !=2.0.0",
			rA:       &Range{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
			rB:       MustNewConstraint("2-2.4 && !=2.0.0 || 2.1-2.2 && !=2.0.0"),
			expected: false,
		},
		{
			name:     "2.0.0 - 3.0.0 does contain 2-2.4 || 2.5-2.6 && !=4.0.0",
			rA:       &Range{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
			rB:       MustNewConstraint("2-2.4 || 2.5-2.6&& !=4.0.0"),
			expected: true,
		},
		{
			name:     "2.0.0 - 3.0.0 || 5.0.0 - 6.0.0 does contain !=4.0.0 && 5.2.0-6.0.0",
			rA:       MustNewConstraint("2 - 3 || 5 - 6"),
			rB:       MustNewConstraint("!=4 && 5.2-6"),
			expected: true,
		},
		{
			name:     "4.12-4.14 does contain ~4.12",
			rA:       MustNewConstraint("4.12-4.14"),
			rB:       MustNewConstraint("~4.12"),
			expected: true,
		},
		{
			name:     "4.12-4.14 does contain >=4.12,<4.13",
			rA:       MustNewConstraint("4.12-4.14"),
			rB:       MustNewConstraint(">=4.12,<4.13"),
			expected: true,
		},
		{
			name:     "~4.11 does not contain ~4.11 && !=4.11.3",
			rA:       MustNewConstraint("~4.11"),
			rB:       MustNewConstraint("~4.11 && !=4.11.3"),
			expected: false,
		},
		{
			name:     "4.11-4.13 does not contain 4.11-4.13 && !=4.11.3",
			rA:       MustNewConstraint("4.11-4.13"),
			rB:       MustNewConstraint("4.11-4.13 && !=4.11.3"),
			expected: false,
		},
		{
			name:     "4.11-4.13 && !=4.11.3 does not contain 4.11-4.13",
			rA:       MustNewConstraint("4.11-4.13 && !=4.11.3"),
			rB:       MustNewConstraint("4.11-4.13"),
			expected: false,
		},
		{
			name:     "!=4.11.3 does not contain 4.11-4.13",
			rA:       MustNewConstraint("!=4.11.3"),
			rB:       MustNewConstraint("4.11-4.13"),
			expected: false,
		},
		{
			name:     "=4.11.3 does not contain 4.11-4.13",
			rA:       MustNewConstraint("=4.11.3"),
			rB:       MustNewConstraint("4.11-4.13"),
			expected: false,
		},
		{
			name:     "4.11-4.13 does contain =4.11.3",
			rA:       MustNewConstraint("4.11-4.13"),
			rB:       MustNewConstraint("=4.11.3"),
			expected: true,
		},
		{
			name:     "4.11-4.13 does contain !=4.14.3",
			rA:       MustNewConstraint("4.11-4.13"),
			rB:       MustNewConstraint("!=4.14.3"),
			expected: true,
		},
		{
			name:     "2 - 3 && 1 - 2 does contain 1 - 3",
			rA:       MustNewConstraint("2 - 3 && 1 - 2"),
			rB:       MustNewConstraint("1 - 3"),
			expected: true,
		},
		{
			name:     "4.12.x - 4.14.x && != 4.13.5 does contain 4.13.0 - 4.13.4 && 4.13.6 - 4.13.8",
			rA:       MustNewConstraint("4.12.x - 4.14.x && != 4.13.5"),
			rB:       MustNewConstraint("4.13.0 - 4.13.4 && 4.13.6 - 4.13.8"),
			expected: true,
		},

		// Edge cases: Range.Contains(OR) - must contain ALL branches
		{
			name:     "1-3 does NOT contain 1-1.5 || 2.5-3.5 (second branch extends beyond)",
			rA:       MustNewConstraint("1-3"),
			rB:       MustNewConstraint("1-1.5 || 2.5-3.5"),
			expected: false, // 3.5 > 3
		},
		{
			name:     "1-4 does contain 1-1.5 || 2.5-3.5 (both branches within bounds)",
			rA:       MustNewConstraint("1-4"),
			rB:       MustNewConstraint("1-1.5 || 2.5-3.5"),
			expected: true,
		},
		{
			name:     "1-3 does NOT contain 0.5-1.5 || 2-3 (first branch extends below)",
			rA:       MustNewConstraint("1-3"),
			rB:       MustNewConstraint("0.5-1.5 || 2-3"),
			expected: false, // 0.5 < 1
		},
		{
			name:     ">=2 does NOT contain >=1.5 || >=3 (first branch too wide)",
			rA:       MustNewConstraint(">=2"),
			rB:       MustNewConstraint(">=1.5 || >=3"),
			expected: false, // >=1.5 includes versions < 2
		},
		{
			name:     "<=3 does NOT contain <=2 || <=4 (second branch too wide)",
			rA:       MustNewConstraint("<=3"),
			rB:       MustNewConstraint("<=2 || <=4"),
			expected: false, // <=4 extends beyond <=3
		},
		{
			name:     "=2.0.0 does NOT contain =2.0.0 || =3.0.0 (includes extra version)",
			rA:       MustNewConstraint("=2.0.0"),
			rB:       MustNewConstraint("=2.0.0 || =3.0.0"),
			expected: false,
		},

		// Edge cases: OR.Contains(OR)
		{
			name:     "1-4 || 5-6 does contain 1.5-3.5 || 5.2-5.8 (each branch covered)",
			rA:       MustNewConstraint("1-4 || 5-6"),
			rB:       MustNewConstraint("1.5-3.5 || 5.2-5.8"),
			expected: true,
		},
		{
			name:     "1-2 || 3-4 does NOT contain 1.5-2 || 2.5-3.5 (gap not covered)",
			rA:       MustNewConstraint("1-2 || 3-4"),
			rB:       MustNewConstraint("1.5-2 || 2.5-3.5"),
			expected: false, // 2.5 falls in gap between 2 and 3
		},
		{
			name:     "1-3 || 5-7 || 9-11 does contain 1-2 || 6-7 || 10-11",
			rA:       MustNewConstraint("1-3 || 5-7 || 9-11"),
			rB:       MustNewConstraint("1-2 || 6-7 || 10-11"),
			expected: true,
		},

		// Edge cases: AND with multiple !=
		{
			name:     "1-5 && !=2 && !=3 does contain 1.5-4.5 && !=2 && !=3",
			rA:       MustNewConstraint("1-5 && !=2 && !=3"),
			rB:       MustNewConstraint("1.5-4.5 && !=2 && !=3"),
			expected: true,
		},
		{
			name:     "1-5 && !=2 does NOT contain 1-5 && !=3 (different exclusions)",
			rA:       MustNewConstraint("1-5 && !=2"),
			rB:       MustNewConstraint("1-5 && !=3"),
			expected: false,
		},
		{
			name:     "1-10 && !=3 && !=5 does contain 2-4 && !=3 && !=5",
			rA:       MustNewConstraint("1-10 && !=3 && !=5"),
			rB:       MustNewConstraint("2-4 && !=3 && !=5"),
			expected: true,
		},

		// Edge cases: Mixed AND/OR with !=
		{
			name:     "1-5 && !=3 does NOT contain 1-2 || 3-4 (OR includes excluded version)",
			rA:       MustNewConstraint("1-5 && !=3"),
			rB:       MustNewConstraint("1-2 || 3-4"),
			expected: false, // rB includes 3.0.0
		},
		{
			name:     "1-5 && !=3 does contain 1-2 || 4-5 (OR avoids excluded version)",
			rA:       MustNewConstraint("1-5 && !=3"),
			rB:       MustNewConstraint("1-2 || 4-5"),
			expected: true,
		},
		{
			name:     "1-2 || 3-4 does NOT contain 1.5-3.5 && !=2 (AND crosses gap)",
			rA:       MustNewConstraint("1-2 || 3-4"),
			rB:       MustNewConstraint("1.5-3.5 && !=2"),
			expected: false,
		},

		// Edge cases: Tilde and Caret with OR
		{
			name:     "~1 does contain ~1.2 || ~1.5",
			rA:       MustNewConstraint("~1"),
			rB:       MustNewConstraint("~1.2 || ~1.5"),
			expected: true,
		},
		{
			name:     "~1.2 || ~1.5 does NOT contain ~1 (gaps between ranges)",
			rA:       MustNewConstraint("~1.2 || ~1.5"),
			rB:       MustNewConstraint("~1"),
			expected: false,
		},
		{
			name:     "^1.0.0 does contain ^1.2.0 || ^1.5.0",
			rA:       MustNewConstraint("^1.0.0"),
			rB:       MustNewConstraint("^1.2.0 || ^1.5.0"),
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := test.rA.Contains(test.rB)
			assert.Equal(t, test.expected, r)
		})
	}
}
