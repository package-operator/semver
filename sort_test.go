package semver

import (
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAscendingDescendingSort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "simple",
			input: []string{
				"1.2.4",
				"1.2.3",
				"1.0.0",
				"1.3.0",
				"2.0.0",
				"0.4.2",
			},
			expected: []string{
				"0.4.2",
				"1.0.0",
				"1.2.3",
				"1.2.4",
				"1.3.0",
				"2.0.0",
			},
		},
		{
			name: "pre-releases-1",
			input: []string{
				"1.0.0-beta.2",
				"1.0.0-alpha",
				"1.0.0-alpha.beta",
				"1.0.0",
				"1.0.0-rc.1",
				"1.0.0-alpha.1",
				"1.0.0-beta",
				"1.0.0-beta.11",
			},
			expected: []string{
				"1.0.0-alpha",
				"1.0.0-alpha.1",
				"1.0.0-alpha.beta",
				"1.0.0-beta",
				"1.0.0-beta.2",
				"1.0.0-beta.11",
				"1.0.0-rc.1",
				"1.0.0",
			},
		},
		{
			name: "pre-releases-2",
			input: []string{
				"1.3.0-rc.5",
				"1.3.0-rc.0",
				"1.3.0-alpha+banana",
				"1.3.0-alpha.2+banana",
				"1.2.3",
				"1.0.0",
				"1.3.0",
				"2.0.0",
				"0.4.2",
			},
			expected: []string{
				"0.4.2",
				"1.0.0",
				"1.2.3",
				"1.3.0-alpha+banana",
				"1.3.0-alpha.2+banana",
				"1.3.0-rc.0",
				"1.3.0-rc.5",
				"1.3.0",
				"2.0.0",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var list []Version
			for _, vs := range test.input {
				v, err := NewVersion(vs)
				require.NoError(t, err)
				list = append(list, v)
			}

			// test Ascending
			sort.Sort(Ascending(list))
			ascOut := make([]string, len(list))
			for i := range list {
				ascOut[i] = list[i].String()
			}
			assert.Equal(t, test.expected, ascOut)

			// test Descending
			slices.Reverse(test.expected)
			desc := Descending(list)
			sort.Sort(desc)

			descOut := make([]string, len(list))
			for i := range list {
				descOut[i] = list[i].String()
			}
			assert.Equal(t, test.expected, descOut)
		})
	}
}
