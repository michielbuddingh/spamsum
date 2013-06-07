package spamsum

import (
	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	tests := []struct {
		input      string
		shouldfail bool
	}{
		{"49152:dihMNzhZt62oh9+onrqMPr/KwJsvD/mMplt:Hxxpj", false},
		{"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A:e4XuptH8D//4crHMmUfL", false},
		{"18446744073709551616:dihMNzhZt62oh9+onrqMPr/KwJsvD/mMplt:H.soa", true},
		{"49152:dihMNzhZt62oh9+onrqMPr/KwJsvD/mMplt.Hxxpj", true},
		{"22:i3wkMEgPthpID7YoQDjrdAjGBwBIg8Qow0iLSAhIi3AQSItCCEiLUhBIOch1MEiJBCRIiVQkCEiJ:UxUp", true},
	}

	for _, test := range tests {
		var sum SpamSum
		_, err := fmt.Sscan(test.input, &sum)
		if test.shouldfail && err == nil {
			t.Errorf("Should not be able to parse %s\n", test.input)
		}
		if !test.shouldfail && err != nil {
			t.Errorf("Parse failed with error: %v", err)
		}
	}
}
