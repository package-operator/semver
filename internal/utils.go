package internal

import "fmt"

func IsDigit(r rune) bool {
	return r == '0' || IsPositiveDigit(r)
}

func IsPositiveDigit(r rune) bool {
	switch r {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

type Position int

// String returns a string representation of the position.
// e.g. col 5:.
func (pos Position) String() string {
	return fmt.Sprintf("col %d", pos)
}
