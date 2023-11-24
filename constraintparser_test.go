package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx
func TestConstraintParser_success(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected Constraint
	}{
		{
			name:  "simple range",
			input: `1.2.3 - 1.3.4`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 3},
				Max: Version{Major: 1, Minor: 3, Patch: 4},
			},
		},
		{
			name:  "wildcard patch range",
			input: `1.2.x - 1.3.x`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 0},
				Max: Version{Major: 1, Minor: 3, Patch: maxUint64},
			},
		},
		{
			name:  "wildcard minor range",
			input: `1.x - 2.x`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 0, Patch: 0},
				Max: Version{Major: 2, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "wildcard major range start",
			input: `x - 3`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 3, Minor: 0, Patch: 0},
			},
		},
		{
			name:  "wildcard major range end",
			input: `1 - x`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 0, Patch: 0},
				Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "or range",
			input: `1.2.3 - 1.3.4 || 2.3.4 - 4.5.3`,
			expected: or{
				&Range{
					Min: Version{Major: 1, Minor: 2, Patch: 3},
					Max: Version{Major: 1, Minor: 3, Patch: 4},
				},
				&Range{
					Min: Version{Major: 2, Minor: 3, Patch: 4},
					Max: Version{Major: 4, Minor: 5, Patch: 3},
				},
			},
		},
		{
			name:  "complex OR AND",
			input: `=1.2.3 || >1.2.3 <5.4.0 && 1.2.3 - 2.3.4 || !=3`,
			expected: or{
				&Range{
					Min: Version{Major: 1, Minor: 2, Patch: 3},
					Max: Version{Major: 1, Minor: 2, Patch: 3},
				},
				and{
					&Range{
						Min: Version{Major: 1, Minor: 2, Patch: 4},
						Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
					},
					&Range{
						Min: Version{Major: 0, Minor: 0, Patch: 0},
						Max: Version{Major: 5, Minor: 3, Patch: maxUint64},
					},
					&Range{
						Min: Version{Major: 1, Minor: 2, Patch: 3},
						Max: Version{Major: 2, Minor: 3, Patch: 4},
					},
				},
				not{
					Range{
						Min: Version{Major: 3, Minor: 0, Patch: 0},
						Max: Version{Major: 3, Minor: 0, Patch: 0},
					},
				},
			},
		},
		{
			name:  "equal",
			input: `=1.2.3`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 3},
				Max: Version{Major: 1, Minor: 2, Patch: 3},
			},
		},
		{
			name:  "not equal",
			input: `!=1.2.3`,
			expected: not{
				Range{
					Min: Version{Major: 1, Minor: 2, Patch: 3},
					Max: Version{Major: 1, Minor: 2, Patch: 3},
				},
			},
		},
		{
			name:  "greater",
			input: `>1.2.3`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 4},
				Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "greater minor",
			input: `>1.2`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 3, Patch: 0},
				Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "greater minor wildcard",
			input: `>1.2.x`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 3, Patch: 0},
				Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "greater major",
			input: `>1`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 0, Patch: 0},
				Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "greater major wildcard",
			input: `>1.x`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 0, Patch: 0},
				Max: Version{Major: maxUint64, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "greater equal",
			input: `>=1.2.3`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 3},
				Max: Version{Major: maxUint64, Minor: 0, Patch: 0},
			},
		},
		{
			name:  "less",
			input: `<1.2.3`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 1, Minor: 2, Patch: 2},
			},
		},
		{
			name:  "less minor",
			input: `<1.2`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 1, Minor: 1, Patch: maxUint64},
			},
		},
		{
			name:  "less minor wildcard",
			input: `<1.2.x`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 1, Minor: 1, Patch: maxUint64},
			},
		},
		{
			name:  "less major",
			input: `<1`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 0, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "less major wildcard",
			input: `<1.x`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 0, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "less equal",
			input: `<=1.2.3`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 1, Minor: 2, Patch: 3},
			},
		},
		{
			name:  "tilde major",
			input: `~1`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 0, Patch: 0},
				Max: Version{Major: 1, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "tilde minor wildcard",
			input: `~1.x`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 0, Patch: 0},
				Max: Version{Major: 1, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "tilde minor",
			input: `~2.3`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 3, Patch: 0},
				Max: Version{Major: 2, Minor: 3, Patch: maxUint64},
			},
		},
		{
			name:  "tilde patch wildcard",
			input: `~2.3.x`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 3, Patch: 0},
				Max: Version{Major: 2, Minor: 3, Patch: maxUint64},
			},
		},
		{
			name:  "tilde patch",
			input: `~1.2.3`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 3},
				Max: Version{Major: 1, Minor: 2, Patch: maxUint64},
			},
		},
		{
			name:  "caret stable patch",
			input: `^1.2.3`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 3},
				Max: Version{Major: 1, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "caret stable patch wildcard",
			input: `^1.2.x`,
			expected: &Range{
				Min: Version{Major: 1, Minor: 2, Patch: 0},
				Max: Version{Major: 1, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "caret stable minor",
			input: `^2.3`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 3, Patch: 0},
				Max: Version{Major: 2, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "caret stable minor wildcard",
			input: `^2.x`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 0, Patch: 0},
				Max: Version{Major: 2, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "caret stable major",
			input: `^2`,
			expected: &Range{
				Min: Version{Major: 2, Minor: 0, Patch: 0},
				Max: Version{Major: 2, Minor: maxUint64, Patch: maxUint64},
			},
		},
		{
			name:  "caret unstable patch",
			input: `^0.2.3`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 2, Patch: 3},
				Max: Version{Major: 0, Minor: 2, Patch: maxUint64},
			},
		},
		{
			name:  "caret unstable minor",
			input: `^0.2`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 2, Patch: 0},
				Max: Version{Major: 0, Minor: 2, Patch: maxUint64},
			},
		},
		{
			name:  "caret unstable minor",
			input: `^0`,
			expected: &Range{
				Min: Version{Major: 0, Minor: 0, Patch: 0},
				Max: Version{Major: 0, Minor: 0, Patch: maxUint64},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c, err := NewConstraint(test.input)
			require.NoError(t, err)

			assert.Equal(t, test.expected, c)
		})
	}
}

func TestConstraintParser_error(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input       string
		expectedErr string
	}{
		{
			input:       `1.2.3.3.4`,
			expectedErr: "col 6: found 3rd dot when parsing semver",
		},
		{
			input:       `||`,
			expectedErr: "col 1: OR empty range constraint",
		},
		{
			input:       `|b`,
			expectedErr: "col 2: unexpected character U+0062 'b'",
		},
		{
			input:       `&&`,
			expectedErr: "col 1: AND empty range constraint",
		},
		{
			input:       `&x`,
			expectedErr: "col 2: unexpected character U+0078 'x'",
		},
		{
			input:       `1.2.  3`,
			expectedErr: "col 5: semver clause incomplete",
		},
		{
			input:       `1.2 -- 3`,
			expectedErr: "col 6: double hyphen in range constraint",
		},
		{
			input:       `= 1.2.  3`,
			expectedErr: "col 7: semver clause incomplete",
		},
		{
			input:       `= \n`,
			expectedErr: "col 4: unexpected character U+006E 'n'",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			_, err := NewConstraint(test.input)
			require.EqualError(t, err, test.expectedErr)
		})
	}
}

func TestConstraintParser_invalidBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version     []byte
		expectedErr string
	}{
		{
			version:     []byte{0},
			expectedErr: "col 1: illegal character NUL",
		},
		{
			version:     []byte("\xc3\x28"),
			expectedErr: "col 1: illegal UTF-8 encoding",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(string(test.version), func(t *testing.T) {
			t.Parallel()
			_, err := parseConstraint(test.version)
			require.EqualError(t, err, test.expectedErr)
		})
	}
}
