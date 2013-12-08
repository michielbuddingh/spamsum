package spamsum

import (
	"encoding/binary"
	"fmt"
	"math/rand"
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
		if !test.shouldfail && sum.String() != test.input {
			t.Errorf("scanned sum %s is not equal to input string %s\n", sum.String(), test.input)
		}
	}
}

func TestHashBytes(t *testing.T) {
	tests := []struct {
		seed      int64
		length    int
		blocksize uint32
		expected  string
	}{
		{42, 16384, 384, "384:PnwCSZ6yE9r4UCZ1he34xas/E8AhHgdd2yM:PbSZ6yE9rGfExx"},
		{1000, 2048, 48, "48:Zo+v/bCSly4VhreHwHJdkHTzF7sjBU1YuD/QtFsByxoSJW+QiLlH:uSWSFteQHJd+Tp79mqSqyCt+5LlH"},
		{71268, 24, 3, "3:N0n6xmcFctn:7xmptn"},
	}

	for _, test := range tests {
		byteSlice := make([]byte, test.length)
		generator := rand.New(rand.NewSource(test.seed))

		for i := 0; i < test.length/4; i++ {
			binary.BigEndian.PutUint32(byteSlice[i*4:], generator.Uint32())
		}

		sum := HashBytes(byteSlice)

		if sum.String() != test.expected {
			t.Errorf("Expected %v, result was %v", test.expected, sum)
		}
	}
}

func TestBlocksizeAdjustment(t *testing.T) {
	byteSlice := make([]byte, 17921)
	generator := rand.New(rand.NewSource(191))

	i := 0
	for ; i < 24; i++ {
		binary.BigEndian.PutUint32(byteSlice[i*4:], generator.Uint32())
	}

	for ; i < 17921; i++ {
		byteSlice[i] = 0
	}

	sum := HashBytes(byteSlice)
	expected := "3:Bl5KOiWl/:ldZ/"

	if sum.String() != expected {
		t.Errorf("Expected %v, result was %v", expected, sum)
	}
}
