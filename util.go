package semver

import "unicode"

func compareSegment(v, o uint64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}
	return 0
}

func isBuildIdentifier(s string) bool {
	return isDigits(s) || isAlphaNumericIdentifier(s)
}

func isPreReleaseIdentifier(s string) bool {
	return isAlphaNumericIdentifier(s) || isNumericIdentifier(s)
}

func isAlphaNumericIdentifier(s string) bool {
	if len(s) == 1 && isNonDigit(rune(s[0])) {
		// must be a non-diget if len==1
		return false
	}
	// must contain one non-diget
	var foundNonDigit bool
	for _, char := range s {
		if isNonDigit(char) {
			foundNonDigit = true
		}
		if !isIdentifierChar(char) {
			return false
		}
	}
	return foundNonDigit
}

func isNumericIdentifier(s string) bool {
	if s == "0" {
		return true
	}
	for i, char := range s {
		if i == 0 && !isPositiveDigit(char) {
			// must start with positive digit
			return false
		}
		if !isDigit(char) {
			// must be only digits
			return false
		}
	}
	return true
}

func isDigits(s string) bool {
	for _, char := range s {
		if !isDigit(char) {
			return false
		}
	}
	return true
}

func isIdentifierChar(r rune) bool {
	return isDigit(r) || isNonDigit(r)
}

func isNonDigit(r rune) bool {
	return r == '-' || isLetter(r)
}

func isLetter(r rune) bool {
	unicode.IsLetter(r)
	switch r {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z':
		return true
	}
	return false
}

func isDigit(r rune) bool {
	return r == '0' || isPositiveDigit(r)
}

func isPositiveDigit(r rune) bool {
	switch r {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}
