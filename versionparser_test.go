package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const alotOfNines uint64 = 9999999999999999999

var toPR = ToPreReleaseIdentifier

func TestParser_success(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected Version
	}{
		{
			version: "0.0.4",
			expected: Version{
				Major: 0, Minor: 0, Patch: 4,
			},
		},
		{
			version: "0.0.4-rc.10",
			expected: Version{
				Major: 0, Minor: 0, Patch: 4,
				PreRelease: PreReleaseIdentifierList{
					toPR("rc"), toPR("10"),
				},
			},
		},
		{
			version: "1.2.3",
			expected: Version{
				Major: 1, Minor: 2, Patch: 3,
			},
		},
		{
			version: "10.20.30",
			expected: Version{
				Major: 10, Minor: 20, Patch: 30,
			},
		},
		{
			version: "1.1.2-prerelease+meta",
			expected: Version{
				Major: 1, Minor: 1, Patch: 2,
				PreRelease: []PreReleaseIdentifier{
					toPR("prerelease"),
				},
				BuildMetadata: []string{"meta"},
			},
		},
		{
			version: "1.1.2+meta-valid",
			expected: Version{
				Major: 1, Minor: 1, Patch: 2,
				BuildMetadata: []string{"meta-valid"},
			},
		},
		{
			version: "1.0.0-alpha.beta.1",
			expected: Version{
				Major: 1, Minor: 0, Patch: 0,
				PreRelease: []PreReleaseIdentifier{toPR("alpha"), toPR("beta"), toPR("1")},
			},
		},
		{
			version: "1.0.0-alpha.0valid",
			expected: Version{
				Major: 1, Minor: 0, Patch: 0,
				PreRelease: []PreReleaseIdentifier{toPR("alpha"), toPR("0valid")},
			},
		},
		{
			version: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
			expected: Version{
				Major: 1, Minor: 0, Patch: 0,
				PreRelease:    []PreReleaseIdentifier{toPR("alpha-a"), toPR("b-c-somethinglong")},
				BuildMetadata: []string{"build", "1-aef", "1-its-okay"},
			},
		},
		{
			version: "1.0.0-rc.1+build.1",
			expected: Version{
				Major: 1, Minor: 0, Patch: 0,
				PreRelease:    []PreReleaseIdentifier{toPR("rc"), toPR("1")},
				BuildMetadata: []string{"build", "1"},
			},
		},
		{
			version: "2.0.0-rc.1+build.123",
			expected: Version{
				Major: 2, Minor: 0, Patch: 0,
				PreRelease:    []PreReleaseIdentifier{toPR("rc"), toPR("1")},
				BuildMetadata: []string{"build", "123"},
			},
		},
		{
			version: "10.2.3-DEV-SNAPSHOT",
			expected: Version{
				Major: 10, Minor: 2, Patch: 3,
				PreRelease: []PreReleaseIdentifier{toPR("DEV-SNAPSHOT")},
			},
		},
		{
			version: "1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
			expected: Version{
				Major: 1, Minor: 2, Patch: 3,
				PreRelease:    []PreReleaseIdentifier{toPR("---RC-SNAPSHOT"), toPR("12"), toPR("9"), toPR("1--"), toPR("12")},
				BuildMetadata: []string{"788"},
			},
		},
		{
			version: "1.2.3----R-S.12.9.1--.12+meta",
			expected: Version{
				Major: 1, Minor: 2, Patch: 3,
				PreRelease:    []PreReleaseIdentifier{toPR("---R-S"), toPR("12"), toPR("9"), toPR("1--"), toPR("12")},
				BuildMetadata: []string{"meta"},
			},
		},
		{
			version: "1.0.0+0.build.1-rc.10000aaa-kk-0.1",
			expected: Version{
				Major: 1, Minor: 0, Patch: 0,
				BuildMetadata: []string{"0", "build", "1-rc", "10000aaa-kk-0", "1"},
			},
		},
		{
			version: "9999999999999999999.9999999999999999999.9999999999999999999",
			expected: Version{
				Major: alotOfNines, Minor: alotOfNines, Patch: alotOfNines,
			},
		},
		{
			version: "1.2.3-alpha.beta.gamma+banana",
			expected: Version{
				Major: 1, Minor: 2, Patch: 3,
				PreRelease:    []PreReleaseIdentifier{toPR("alpha"), toPR("beta"), toPR("gamma")},
				BuildMetadata: []string{"banana"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			t.Parallel()
			v, err := NewVersion(test.version)
			require.NoError(t, err)
			assert.Equal(t, test.expected, v)
		})
	}
}

func TestParser_error(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version     string
		expectedErr string
	}{
		{
			version:     "1",
			expectedErr: "col 2: missing minor",
		},
		{
			version:     "1.2",
			expectedErr: "col 4: missing patch",
		},
		{
			version:     "1.2.3-0123",
			expectedErr: `col 7: invalid pre release identifier "0123"`,
		},
		{
			version:     "1.2.3-0123.0123",
			expectedErr: `col 7: invalid pre release identifier "0123"`,
		},
		{
			version:     "1.1.2+.123",
			expectedErr: `col 7: build identifier empty`,
		},
		{
			version:     "+invalid",
			expectedErr: `col 1: starts with non-positive integer '+'`,
		},
		{
			version:     "-invalid",
			expectedErr: `col 1: starts with non-positive integer '-'`,
		},
		{
			version:     "-invalid+invalid",
			expectedErr: `col 1: starts with non-positive integer '-'`,
		},
		{
			version:     "-invalid+invalid.01",
			expectedErr: `col 1: starts with non-positive integer '-'`,
		},
		{
			version:     "alpha",
			expectedErr: `col 1: starts with non-positive integer 'a'`,
		},
		{
			version:     "alpha..",
			expectedErr: `col 1: starts with non-positive integer 'a'`,
		},
		{
			version:     "1.0.0-alpha_beta",
			expectedErr: `col 7: invalid pre release identifier "alpha_beta"`,
		},
		{
			version:     "1.0.0-alpha..",
			expectedErr: `col 13: pre release identifier empty`,
		},
		{
			version:     "1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
			expectedErr: `col 8: invalid character '2'`,
		},
		{
			version:     "-1.0.3-gamma+b7718",
			expectedErr: `col 2: starts with non-positive integer '-'`,
		},
		{
			version:     "1.0.\n3-gamma+b7718",
			expectedErr: `col 5: illegal character NEWLINE`,
		},
		{
			version:     "1..",
			expectedErr: `col 3: expected number, got nothing`,
		},
		{
			version:     "1.  .",
			expectedErr: `col 3: illegal character SPACE`,
		},
		{
			version:     "1.2 .",
			expectedErr: `col 4: illegal character SPACE`,
		},
		{
			version:     "1.2.3+   ",
			expectedErr: `col 7: illegal character SPACE`,
		},
		{
			version:     "",
			expectedErr: `col 1: missing major`,
		},
	}
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			t.Parallel()
			_, err := NewVersion(test.version)
			require.EqualError(t, err, test.expectedErr)
		})
	}
}

func TestParser_invalidBytes(t *testing.T) {
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
		_, err := parseVersion(test.version)
		require.EqualError(t, err, test.expectedErr)
	}
}
