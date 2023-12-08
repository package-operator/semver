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
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := test.r.Check(test.v)
			assert.Equal(t, test.expected, r)
		})
	}
}

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
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := test.rA.Contains(test.rB)
			assert.Equal(t, test.expected, r)
		})
	}
}
