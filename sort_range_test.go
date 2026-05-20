package semver

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAscendingMin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []Range
		expected []Range
	}{
		{
			name: "sort by min version ascending",
			input: []Range{
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.1.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.1.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
			},
		},
		{
			name: "already sorted",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
			},
		},
		{
			name: "reverse sorted",
			input: []Range{
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
			},
		},
		{
			name: "same min versions with different max",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
			},
		},
		{
			name: "minor version differences",
			input: []Range{
				{Min: MustNewVersion("1.5.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.2.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.2.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.5.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "patch version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.5"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.1"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.3"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.1"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.3"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.5"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name:     "empty slice",
			input:    []Range{},
			expected: []Range{},
		},
		{
			name: "single element",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]Range, len(tt.input))
			copy(input, tt.input)
			sort.Sort(AscendingMin(input))
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestAscendingMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []Range
		expected []Range
	}{
		{
			name: "sort by max version ascending",
			input: []Range{
				{Min: MustNewVersion("0.1.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("0.4.0"), Max: MustNewVersion("3.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("0.4.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("0.1.0"), Max: MustNewVersion("4.0.0")},
			},
		},
		{
			name: "already sorted",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
			},
		},
		{
			name: "reverse sorted",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
			},
		},
		{
			name: "same max versions with different min",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("5.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("5.0.0")},
			},
		},
		{
			name: "minor version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.5.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.1.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.1.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.5.0")},
			},
		},
		{
			name: "patch version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.7")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.3")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.5")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.3")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.5")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.7")},
			},
		},
		{
			name:     "empty slice",
			input:    []Range{},
			expected: []Range{},
		},
		{
			name: "single element",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]Range, len(tt.input))
			copy(input, tt.input)
			sort.Sort(AscendingMax(input))
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestDescendingMin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []Range
		expected []Range
	}{
		{
			name: "sort by min version descending",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.1.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.1.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
			},
		},
		{
			name: "already sorted",
			input: []Range{
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "reverse sorted",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "same min versions with different max",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
			},
		},
		{
			name: "minor version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.2.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.5.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.5.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.2.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "patch version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.1"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.5"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.3"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.5"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.3"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.1"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name:     "empty slice",
			input:    []Range{},
			expected: []Range{},
		},
		{
			name: "single element",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]Range, len(tt.input))
			copy(input, tt.input)
			sort.Sort(DescendingMin(input))
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestDescendingMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []Range
		expected []Range
	}{
		{
			name: "sort by max version descending",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("0.1.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("0.4.0"), Max: MustNewVersion("3.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("0.1.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("0.4.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "already sorted",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "reverse sorted",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("4.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("3.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "same max versions with different min",
			input: []Range{
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("5.0.0")},
				{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("5.0.0")},
			},
		},
		{
			name: "minor version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.1.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.5.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.5.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.1.0")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
		{
			name: "patch version differences",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.3")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.7")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.5")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.7")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.5")},
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.3")},
			},
		},
		{
			name:     "empty slice",
			input:    []Range{},
			expected: []Range{},
		},
		{
			name: "single element",
			input: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
			expected: []Range{
				{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]Range, len(tt.input))
			copy(input, tt.input)
			sort.Sort(DescendingMax(input))
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestSortInterface_Len(t *testing.T) {
	t.Parallel()

	ranges := []Range{
		{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
		{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
		{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
	}

	assert.Equal(t, 3, AscendingMin(ranges).Len())
	assert.Equal(t, 3, AscendingMax(ranges).Len())
	assert.Equal(t, 3, DescendingMin(ranges).Len())
	assert.Equal(t, 3, DescendingMax(ranges).Len())

	assert.Equal(t, 0, AscendingMin([]Range{}).Len())
}

func TestSortInterface_Swap(t *testing.T) {
	t.Parallel()

	ranges := []Range{
		{Min: MustNewVersion("1.0.0"), Max: MustNewVersion("2.0.0")},
		{Min: MustNewVersion("2.0.0"), Max: MustNewVersion("3.0.0")},
		{Min: MustNewVersion("3.0.0"), Max: MustNewVersion("4.0.0")},
	}

	AscendingMin(ranges).Swap(0, 2)
	assert.Equal(t, MustNewVersion("3.0.0"), ranges[0].Min)
	assert.Equal(t, MustNewVersion("1.0.0"), ranges[2].Min)
}
