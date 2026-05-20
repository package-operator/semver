package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnd(t *testing.T) {
	t.Parallel()
	t.Run("all true", func(t *testing.T) {
		t.Parallel()
		and := and{
			&positiveConstraint{}, &positiveConstraint{}, &positiveConstraint{},
		}

		assert.True(t, and.Check(Version{}))
		assert.True(t, and.Contains(&Range{}))
	})

	t.Run("one false", func(t *testing.T) {
		t.Parallel()
		and := and{
			&positiveConstraint{}, &negativeConstraint{}, &positiveConstraint{},
		}

		assert.False(t, and.Check(Version{}))
		assert.False(t, and.Contains(&Range{}))
	})

	t.Run("compaction structure", func(t *testing.T) {
		t.Parallel()
		c, err := NewConstraint("1 - 2 && 2 - 3 && 1.1 - 10")
		require.NoError(t, err)
		oic := c.(*originalInputConstraint)
		r, isSingleRange := oic.Constraint.(*Range)
		if isSingleRange {
			assert.True(t, r.Min.Same(Version{Major: 2}) && r.Max.Same(Version{Major: 2}),
				"intersection should be {2.0.0}, got %s", r.String())
		}
	})

	t.Run("compaction Check", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			constraint string
			version    string
			expected   bool
		}{
			{"adjacent AND matches at boundary", "1 - 2 && 2 - 3", "2.0.0", true},
			{"adjacent AND rejects left-only version", "1 - 2 && 2 - 3", "1.5.0", false},
			{"adjacent AND rejects right-only version", "1 - 2 && 2 - 3", "2.5.0", false},
			{"adjacent AND rejects left boundary", "1 - 2 && 2 - 3", "1.0.0", false},
			{"adjacent AND rejects right boundary", "1 - 2 && 2 - 3", "3.0.0", false},
			{"adjacent AND rejects below both", "1 - 2 && 2 - 3", "0.5.0", false},
			{"adjacent AND rejects above both", "1 - 2 && 2 - 3", "3.5.0", false},
			{"overlapping AND matches intersection start", "1.0.0 - 3.0.0 && 2.0.0 - 4.0.0", "2.0.0", true},
			{"overlapping AND matches intersection mid", "1.0.0 - 3.0.0 && 2.0.0 - 4.0.0", "2.5.0", true},
			{"overlapping AND matches intersection end", "1.0.0 - 3.0.0 && 2.0.0 - 4.0.0", "3.0.0", true},
			{"overlapping AND rejects left-only", "1.0.0 - 3.0.0 && 2.0.0 - 4.0.0", "1.5.0", false},
			{"overlapping AND rejects right-only", "1.0.0 - 3.0.0 && 2.0.0 - 4.0.0", "3.5.0", false},
			{"identical AND matches inside", "1.0.0 - 2.0.0 && 1.0.0 - 2.0.0", "1.5.0", true},
			{"identical AND rejects outside", "1.0.0 - 2.0.0 && 1.0.0 - 2.0.0", "2.5.0", false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				c, err := NewConstraint(tt.constraint)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, c.Check(MustNewVersion(tt.version)))
			})
		}
	})

	t.Run("compaction Contains", func(t *testing.T) {
		t.Parallel()
		assert.False(t,
			MustNewConstraint("2 - 3 && 1 - 2").Contains(MustNewConstraint("1 - 3")),
			"adjacent AND does not contain full union range")
	})
}

func TestAnd_String(t *testing.T) {
	t.Parallel()
	and := and{
		&Range{
			Min: MustNewVersion("1.1.0"),
			Max: MustNewVersion("1.2.5"),
		},
		&Range{
			Min: MustNewVersion("2.1.0"),
			Max: MustNewVersion("2.2.5"),
		},
	}
	assert.Equal(t, "1.1.0 - 1.2.5 && 2.1.0 - 2.2.5", and.String())
}

func TestOr(t *testing.T) {
	t.Parallel()
	t.Run("one true", func(t *testing.T) {
		t.Parallel()
		or := or{
			&negativeConstraint{}, &negativeConstraint{}, &positiveConstraint{},
		}

		assert.True(t, or.Check(Version{}))
		assert.True(t, or.Contains(&Range{}))
	})

	t.Run("all false", func(t *testing.T) {
		t.Parallel()
		or := or{
			&negativeConstraint{}, &negativeConstraint{}, &negativeConstraint{},
		}

		assert.False(t, or.Check(Version{}))
		assert.False(t, or.Contains(&Range{}))
	})

	t.Run("overlapping branches Contains", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			rA       string
			rB       string
			expected bool
		}{
			{
				name:     "overlapping OR branches cover full range",
				rA:       "1 - 3 || 2 - 5",
				rB:       "1 - 5",
				expected: true,
			},
			{
				name:     "adjacent OR branches cover full range",
				rA:       "1 - 2 || 2 - 3",
				rB:       "1 - 3",
				expected: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, tt.expected,
					MustNewConstraint(tt.rA).Contains(MustNewConstraint(tt.rB)))
			})
		}
	})
}

func TestOr_String(t *testing.T) {
	t.Parallel()
	or := or{
		&Range{
			Min: MustNewVersion("1.1.0"),
			Max: MustNewVersion("1.2.5"),
		},
		&Range{
			Min: MustNewVersion("2.1.0"),
			Max: MustNewVersion("2.2.5"),
		},
	}
	assert.Equal(t, "1.1.0 - 1.2.5 || 2.1.0 - 2.2.5", or.String())
}

func TestNot(t *testing.T) {
	t.Parallel()
	t.Run("true becomes false", func(t *testing.T) {
		t.Parallel()
		not := not{
			Range{},
		}

		assert.False(t, not.Check(Version{}))
		assert.False(t, not.Contains(&Range{}))
	})

	t.Run("false becomes true", func(t *testing.T) {
		t.Parallel()
		not := not{
			Range{
				Min: MustNewVersion("1.0.0"),
				Max: MustNewVersion("2.0.0"),
			},
		}

		assert.True(t, not.Check(MustNewVersion("0.2.0")))
		assert.True(t, not.Contains(&Range{
			Min: MustNewVersion("3.0.0"),
			Max: MustNewVersion("4.0.0"),
		}))
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			rA       string
			rB       string
			expected bool
		}{
			{
				name:     "wider exclusion does not contain narrower exclusion",
				rA:       "1 - 5 && !=2",
				rB:       "1 - 5 && !=2.0.0",
				expected: false,
			},
			{
				name:     "narrower exclusion contains wider exclusion",
				rA:       "1 - 5 && !=2.0.0",
				rB:       "1 - 5 && !=2",
				expected: true,
			},
			{
				name:     "narrow not contains wide not",
				rA:       "!=2.0.0",
				rB:       "!=2",
				expected: true,
			},
			{
				name:     "wide not does not contain narrow not",
				rA:       "!=2",
				rB:       "!=2.0.0",
				expected: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, tt.expected,
					MustNewConstraint(tt.rA).Contains(MustNewConstraint(tt.rB)))
			})
		}
	})
}

func TestNot_String(t *testing.T) {
	t.Parallel()
	not := not{
		Range{
			Min: MustNewVersion("1.1.2"),
			Max: MustNewVersion("1.1.2"),
		},
	}
	assert.Equal(t, "!=1.1.2", not.String())
}

// constraint that is always true.
type positiveConstraint struct{}

func (c *positiveConstraint) Check(_ Version) bool {
	return true
}

func (c *positiveConstraint) Contains(_ Constraint) bool {
	return true
}

func (c *positiveConstraint) String() string {
	return "positiveStub"
}

// constraint that is always false.
type negativeConstraint struct{}

func (c *negativeConstraint) Check(_ Version) bool {
	return false
}

func (c *negativeConstraint) Contains(_ Constraint) bool {
	return false
}

func (c *negativeConstraint) String() string {
	return "negativeStub"
}
