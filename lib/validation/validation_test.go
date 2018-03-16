package validation

import (
	"testing"
)

// testString holds data for testing CheckString
var testString = []struct {
	regex  string
	data   string
	expect bool
}{
	{"^$", "", true},
	{"a(b", "asd", false},
	{"^[a-zA-Z]{1,5}$", "Hello", true},
	{"^[a-z]{1,5}$", "Hello", false},
}

// testCount holds data for testing CheckRuneCount
var testCount = []struct {
	data   string
	min    uint
	max    uint
	expect bool
}{
	{"", 0, 0, true},
	{"", 1, 0, false},
	{"Hello", 5, 5, true},
	{"世界", 0, 2, true},
	{"Hello", 0, 4, false},
	{"世界", 0, 1, false},
}

// TestCheckString
func TestCheckString(t *testing.T) {
	for _, tt := range testString {
		if ret := CheckString(tt.data, tt.regex); ret != tt.expect {
			t.Errorf("TestCheckString (%s, %s) expected %t, got %t\n", tt.regex, tt.data, tt.expect, ret)
		}
	}
}

// TestCheckRuneCount
func TestCheckRuneCount(t *testing.T) {
	for _, tt := range testCount {
		if ret := CheckRuneCount(tt.data, tt.min, tt.max); ret != tt.expect {
			t.Errorf("TestCheckRuneCount expected %t, got %t\n", tt.expect, ret)
		}
	}
}

// EOF
