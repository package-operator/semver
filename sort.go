package semver

import "sort"

// Ascending sorts versions Ascending via the sorts standard lib package.
// resulting order: 1.0.0, 1.1.0, 2.0.0.
type Ascending []Version

var _ sort.Interface = Ascending{}

// Returns the number of items of the slice.
// Implements sort.Interface.
func (l Ascending) Len() int {
	return len(l)
}

// Returns true if item[j] is less than item[i].
// Implements sort.Interface.
func (l Ascending) Less(i, j int) bool {
	return l[i].LessThan(l[j])
}

// Swaps the position of two items in the list.
// Implements sort.Interface.
func (l Ascending) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Descending sorts versions Descending via the sorts standard lib package.
// resulting order: 2.0.0, 1.1.0, 1.0.0.
type Descending []Version

var _ sort.Interface = Descending{}

// Returns the number of items of the slice.
// Implements sort.Interface.
func (l Descending) Len() int {
	return len(l)
}

// Returns true if item[j] is less than item[i].
// Implements sort.Interface.
func (l Descending) Less(i, j int) bool {
	return l[i].GreaterThan(l[j])
}

// Swaps the position of two items in the list.
// Implements sort.Interface.
func (l Descending) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
