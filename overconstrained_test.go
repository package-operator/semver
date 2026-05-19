package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverConstrainedDetection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		constraint  string
		shouldError bool
		errorMsg    string
	}{
		// Valid constraints (should NOT error)
		{
			name:        "valid range",
			constraint:  "1.0.0 - 2.0.0",
			shouldError: false,
		},
		{
			name:        "valid AND with overlap",
			constraint:  ">=1.0.0 && <=2.0.0",
			shouldError: false,
		},
		{
			name:        "valid adjacent ranges",
			constraint:  "1.0.0 - 2.0.0 && 2.0.0 - 3.0.0",
			shouldError: false,
		},
		{
			name:        "valid overlapping ranges",
			constraint:  "1.0.0 - 2.5.0 && 2.0.0 - 3.0.0",
			shouldError: false,
		},
		{
			name:        "valid NOT with range",
			constraint:  "1.0.0 - 2.0.0 && !=1.5.0",
			shouldError: false,
		},

		// Over-constrained: min > max (detected as non-overlapping ranges)
		{
			name:        "impossible min > max",
			constraint:  ">=2.0.0 && <1.0.0",
			shouldError: true,
			errorMsg:    "over-constrained, ranges do not overlap",
		},
		{
			name:        "impossible min > max with exact versions",
			constraint:  ">=5.0.0 && <=3.0.0",
			shouldError: true,
			errorMsg:    "over-constrained, ranges do not overlap",
		},

		// Over-constrained: non-overlapping ranges
		{
			name:        "non-overlapping ranges gap between",
			constraint:  "1.0.0 - 2.0.0 && 3.0.0 - 4.0.0",
			shouldError: true,
			errorMsg:    "over-constrained, ranges do not overlap",
		},
		{
			name:        "non-overlapping ranges adjacent but not touching",
			constraint:  "1.0.0 - 1.9.9 && 2.0.0 - 3.0.0",
			shouldError: true,
			errorMsg:    "over-constrained, ranges do not overlap",
		},
		{
			name:        "non-overlapping with multiple ranges",
			constraint:  "1.0.0 - 1.5.0 && 2.0.0 - 2.5.0 && 3.0.0 - 3.5.0",
			shouldError: true,
			errorMsg:    "over-constrained, ranges do not overlap",
		},

		// Over-constrained: NOT excludes exact version
		{
			name:        "equal and not equal same version",
			constraint:  "=1.0.0 && !=1.0.0",
			shouldError: true,
			errorMsg:    "over-constrained, =1.0.0 AND !=1.0.0 excludes all versions",
		},
		{
			name:        "equal and not equal different order",
			constraint:  "!=2.5.0 && =2.5.0",
			shouldError: true,
			errorMsg:    "over-constrained, =2.5.0 AND !=2.5.0 excludes all versions",
		},

		// Edge cases
		{
			name:        "touching ranges at boundary (valid)",
			constraint:  "1.0.0 - 2.0.0 && 2.0.0 - 3.0.0",
			shouldError: false, // 2.0.0 is in both ranges
		},
		{
			name:        "same range twice (valid but redundant)",
			constraint:  "1.0.0 - 2.0.0 && 1.0.0 - 2.0.0",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewConstraint(tt.constraint)

			if tt.shouldError {
				require.Error(t, err, "expected error for constraint: %s", tt.constraint)
				assert.Contains(t, err.Error(), tt.errorMsg,
					"error message should contain '%s' for constraint: %s", tt.errorMsg, tt.constraint)
			} else {
				require.NoError(t, err, "constraint should be valid: %s", tt.constraint)
			}
		})
	}
}

func TestOverConstrainedWithOperators(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		constraint  string
		shouldError bool
		errorMsg    string
	}{
		// Valid combinations - these may trigger redundancy warnings but are not over-constrained
		{
			name:        "tilde with wider range gets redundancy check",
			constraint:  "~1.2.3 && >=1.2.0",
			shouldError: true, // Redundancy check, not over-constrained
			errorMsg:    "redundant",
		},
		{
			name:        "caret with wider range gets redundancy check",
			constraint:  "^1.2.0 && >=1.0.0",
			shouldError: true, // Redundancy check, not over-constrained
			errorMsg:    "redundant",
		},

		// Invalid combinations - truly over-constrained
		{
			name:        "tilde incompatible with lower range",
			constraint:  "~2.0.0 && <1.0.0",
			shouldError: true,
			errorMsg:    "over-constrained",
		},
		{
			name:        "caret incompatible with higher range",
			constraint:  "^1.0.0 && >=5.0.0",
			shouldError: true,
			errorMsg:    "over-constrained",
		},
		{
			name:        "greater than less than impossible",
			constraint:  ">5.0.0 && <3.0.0",
			shouldError: true,
			errorMsg:    "over-constrained",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewConstraint(tt.constraint)

			if tt.shouldError {
				require.Error(t, err, "expected error for constraint: %s", tt.constraint)
				assert.Contains(t, err.Error(), tt.errorMsg,
					"error should mention '%s' for: %s", tt.errorMsg, tt.constraint)
			} else {
				require.NoError(t, err, "constraint should be valid: %s", tt.constraint)
			}
		})
	}
}

// TestOverConstrainedErrorPosition verifies that error messages include position info.
func TestOverConstrainedErrorPosition(t *testing.T) {
	t.Parallel()

	_, err := NewConstraint(">=2.0.0 && <1.0.0")
	require.Error(t, err)

	// Error should include column position
	assert.Contains(t, err.Error(), "col", "error should include position information")
	assert.Contains(t, err.Error(), "over-constrained", "error should mention over-constrained")
}
