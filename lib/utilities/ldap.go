package utilities

import (
	"regexp"
	"strconv"
)

// ExtractLDAPerrorCode extracts the ldap int code from an error string
func ExtractLDAPerrorCode(err error) int {
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindStringSubmatch(err.Error())

	if len(matches) != 1 {
		return -1
	}

	ret, err := strconv.Atoi(matches[0])
	if err != nil {
		return -1
	}

	return ret
}

// EOF
