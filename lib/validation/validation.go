package validation

import (
	"regexp"
	"unicode/utf8"
)

// CheckString returns true if the string matches the regex
func CheckString(value, regex string) bool {
	matched, err := regexp.MatchString(regex, value)
	if err != nil {
		return false
	}

	return matched
}

// CheckRuneCount accepts a min and max value, returns true if the Rune count is equal or between these values
func CheckRuneCount(value string, min, max uint) bool {
	if min > max {
		return false
	}

	count := uint(utf8.RuneCountInString(value))

	if count < min {
		return false
	}

	if count > max {
		return false
	}

	return true
}

// EOF
