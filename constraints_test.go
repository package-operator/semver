package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}

// constraint that is always true.
type positiveConstraint struct{}

func (c *positiveConstraint) Check(_ Version) bool {
	return true
}

func (c *positiveConstraint) Contains(_ Constraint) bool {
	return true
}

// constraint that is always false.
type negativeConstraint struct{}

func (c *negativeConstraint) Check(_ Version) bool {
	return false
}

func (c *negativeConstraint) Contains(_ Constraint) bool {
	return false
}
