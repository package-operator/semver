package semver

import "sort"

// AscendingMin sorts ranges Ascending by minimal version via the sort standard lib package.
// resulting order: 1.0.0 - 4.0.0, 1.1.0 - 2.0.0, 2.0.0 - 5.0.0.
type AscendingMin []Range

var _ sort.Interface = AscendingMin{}

// Returns the number of items of the slice.
// Implements sort.Interface.
func (l AscendingMin) Len() int {
	return len(l)
}

// Returns true if item[i] should sort before item[j] (descending order).
// Implements sort.Interface.
func (l AscendingMin) Less(i, j int) bool {
	return l[i].Min.LessThan(l[j].Min)
}

// Swaps the position of two items in the list.
// Implements sort.Interface.
func (l AscendingMin) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// AscendingMax sorts ranges Ascending by max version via the sort standard lib package.
// resulting order: 1.0.0 - 2.0.0, 0.4.0 - 3.0.0, 0.1.0 - 4.0.0.
type AscendingMax []Range

var _ sort.Interface = AscendingMax{}

// Returns the number of items of the slice.
// Implements sort.Interface.
func (l AscendingMax) Len() int {
	return len(l)
}

// Returns true if item[j] is less than item[i].
// Implements sort.Interface.
func (l AscendingMax) Less(i, j int) bool {
	return l[i].Max.LessThan(l[j].Max)
}

// Swaps the position of two items in the list.
// Implements sort.Interface.
func (l AscendingMax) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// DescendingMin sorts ranges Descending by min version via the sort standard lib package.
// resulting order: 2.0.0, 1.1.0, 1.0.0.
type DescendingMin []Range

var _ sort.Interface = DescendingMin{}

// Returns the number of items of the slice.
// Implements sort.Interface.
func (l DescendingMin) Len() int {
	return len(l)
}

// Returns true if item[j] is less than item[i].
// Implements sort.Interface.
func (l DescendingMin) Less(i, j int) bool {
	return l[i].Min.GreaterThan(l[j].Min)
}

// Swaps the position of two items in the list.
// Implements sort.Interface.
func (l DescendingMin) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// DescendingMax sorts ranges Descending by max version via the sort standard lib package.
// resulting order: 2.0.0, 1.1.0, 1.0.0.
type DescendingMax []Range

var _ sort.Interface = DescendingMax{}

// Returns the number of items of the slice.
// Implements sort.Interface.
func (l DescendingMax) Len() int {
	return len(l)
}

// Returns true if item[i] should sort before item[j] (descending order).
// Implements sort.Interface.
func (l DescendingMax) Less(i, j int) bool {
	return l[i].Max.GreaterThan(l[j].Max)
}

// Swaps the position of two items in the list.
// Implements sort.Interface.
func (l DescendingMax) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
