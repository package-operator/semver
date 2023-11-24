package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion_String(t *testing.T) {
	t.Parallel()
	// input equals output test
	tests := []string{
		"0.0.4", "10.20.30", "1.1.2-prerelease+meta",
		"1.0.0-alpha.beta.1", "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
		"9999999999999999999.9999999999999999999.9999999999999999999",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
	}
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			t.Parallel()
			v, err := NewVersion(test)
			require.NoError(t, err)
			assert.Equal(t, test, v.String())
		})
	}
}

func TestPreReleaseIdentifier_Compare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pre      PreReleaseIdentifier
		other    PreReleaseIdentifier
		expected int
	}{
		{
			name:     "alpha equals alpha",
			pre:      toPR("alpha"),
			other:    toPR("alpha"),
			expected: 0,
		},
		{
			name:     "alpha before rc",
			pre:      toPR("rc"),
			other:    toPR("alpha"),
			expected: -1,
		},
		{
			name:     "rc after alpha",
			pre:      toPR("alpha"),
			other:    toPR("rc"),
			expected: 1,
		},
		{
			name:     "5 after 1",
			pre:      toPR("1"),
			other:    toPR("5"),
			expected: 1,
		},
		{
			name:     "1 before 5",
			pre:      toPR("5"),
			other:    toPR("1"),
			expected: -1,
		},
		{
			name:     "5 equals 5",
			pre:      toPR("5"),
			other:    toPR("5"),
			expected: 0,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			d := test.pre.Compare(test.other)
			assert.Equal(t, test.expected, d)
		})
	}
}

func TestPreReleaseIdentifierList_Compare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pre      []PreReleaseIdentifier
		other    []PreReleaseIdentifier
		expected int
	}{
		{
			name: "alpha before alpha.1",
			pre: []PreReleaseIdentifier{
				toPR("alpha"),
			},
			other: []PreReleaseIdentifier{
				toPR("alpha"), toPR("1"),
			},
			expected: 1,
		},
		{
			name: "alpha.1 before alpha.beta",
			pre: []PreReleaseIdentifier{
				toPR("alpha"), toPR("1"),
			},
			other: []PreReleaseIdentifier{
				toPR("alpha"), toPR("beta"),
			},
			expected: 1,
		},
		{
			name: "beta before beta.11",
			pre: []PreReleaseIdentifier{
				toPR("beta"),
			},
			other: []PreReleaseIdentifier{
				toPR("beta"), toPR("11"),
			},
			expected: 1,
		},
		{
			name: "none after rc.1",
			pre:  []PreReleaseIdentifier{},
			other: []PreReleaseIdentifier{
				toPR("rc"), toPR("1"),
			},
			expected: -1,
		},
		{
			name: "rc.1 before none",
			pre: []PreReleaseIdentifier{
				toPR("rc"), toPR("1"),
			},
			other:    []PreReleaseIdentifier{},
			expected: 1,
		},
		{
			name:     "none equals none",
			pre:      []PreReleaseIdentifier{},
			other:    []PreReleaseIdentifier{},
			expected: 0,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			pre := PreReleaseIdentifierList(test.pre)
			d := pre.Compare(test.other)
			assert.Equal(t, test.expected, d)
		})
	}
}
