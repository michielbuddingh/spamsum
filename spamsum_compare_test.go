// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0

package spamsum

import (
	"fmt"
	"testing"
)

func TestEliminateRepetition(t *testing.T) {
	teststrings := []struct {
		input, expected string
	}{
		{"AAAABC", "AAABC"},
		{"Qddddddddd", "Qddd"},
		{"AtrU||||v*****pn", "AtrU|||v***pn"},
	}

	for _, pair := range teststrings {
		shortened := string(eliminateRepetition([]byte(pair.input)))
		if shortened != pair.expected {
			t.Errorf("%v shortened should be %v, is %v", pair.input, pair.expected, shortened)
		}
	}
}

func TestHasCommonSubstring(t *testing.T) {
	tests := []struct {
		left, right string
		expected    bool
	}{
		{"Hello, world", "Hello there", false},
		{"abcdefg", "abcdefg", true},
		{"", "", false},
		{"0123456789ABCDEF", "ABCDEF0123456789", true},
		{"321abcdefg321", "abcdefg", true},
		{"123b4567", "123c4567", false},
	}

	for _, test := range tests {
		result := hasCommonSubstring([]byte(test.left), []byte(test.right))
		if result != test.expected {
			condition := "not "
			if test.expected {
				condition = ""
			}
			t.Errorf("\"%v\" and \"%v\" should %shave a common substring of length 7", test.left, test.right, condition)
		}
		mirroredResult := hasCommonSubstring([]byte(test.right), []byte(test.left))
		if mirroredResult != result {
			t.Errorf("Symmetry error for %v and %v", test.left, test.right)
		}
	}
}

func TestEditDistance(t *testing.T) {
	tests := []struct {
		left, right   string
		dist_expected int
	}{
		{"abcdefg", "abcdefg", 0},
		{"abcdefg", "abcqefg", 2},
		{"ABCDEFG", "ABCEDFG", 2},
		{"ooooAAA", "AAAoooo", 6},
		{"oAoooAA", "AAoooAo", 4},
		{"", "1234567", 7},
		{"", "", 0},
		{"HIJKLMN", "JKLMNOPQRST", 8},
		{"UVxeXup8VuH8rD//pcrHBrlG5FWgYJ70A",
			"kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A", 7},
		{"O4XuptH8D//pcrHmgfL", "e4XuptH8D//4crHMmUfL", 7},
		{"kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A",
			"kVxeXup8VuH8rD//4crHBrlGXm5WGYJ70A", 2},
		{"2Ewd+NvN88y3GdkvBC+9lKMHhDh",
			"2Ewd+NvNrgdkvBC+9lKMHhDh", 7},
		{"vEnWHH6d/4H/4Z2fvNoF8Sy2yt/YUC",
			"xLnWHH6d/4H/4HHHHHHHH4CnrJuN0QhsSyjTU9/j4hbp96khuYhwX", 51},
	}

	for _, test := range tests {
		result := editDistance([]byte(test.left), []byte(test.right))
		if result != test.dist_expected {
			t.Errorf("\"%v\" and \"%v\" should have a distance of %d, was %d", test.left, test.right, test.dist_expected, result)
		}
		mirroredResult := editDistance([]byte(test.left), []byte(test.right))
		if mirroredResult != result {
			t.Errorf("Symmetry error, editDistance(%s, %s) should be editDistance(%s, %s)", test.left, test.right, test.right, test.left)
		}
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		left, right    string
		blocksize      int
		score_expected int
	}{
		{"2Ewd+NvN88y3GdkvBC+9lKMHhDh",
			"2Ewd+NvNrgdkvBC+9lKMHhDh", 6, 48},
		{"7iExTmgeXCcGYX1CRRX1PRRX88p0RRpdV/ISGcEvNOk+l/oX9QUopsAoX9QUopIo",
			"7iExTmgeXCcGYX1CRRX1PRRXrZGcEvNOk+l/oX9QUopsAoX9QUopIHKl057DRMHD",
			12, 80},
		{"vEnWHH6d/4H/4Z2fvNoF8Sy2yt/YUC",
			"xLnWHH6d/4H/4HHHHHHHH4CnrJuN0QhsSyjTU9/j4hbp96khuYhwX", 24, 43},
	}
	for _, test := range tests {
		result := score([]byte(test.left), []byte(test.right), test.blocksize)
		if result != test.score_expected {
			t.Errorf("\"%v\" and \"%v\" should have a score of %d, was %d", test.left, test.right, test.score_expected, result)
		}
	}

}

func TestCompare(t *testing.T) {
	tests := []struct {
		left, right         string
		similarity_expected uint32
	}{
		// these are not values produced by the original spamsum
		// score algorithm
		{

			"12582912:UVxeXup8VuH8rD//pcrHBrlG5FWgYJ70A:O4XuptH8D//pcrHmgfL",
			"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A:e4XuptH8D//4crHMmUfL",
			91},

		{"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A:e4XuptH8D//4crHMmUfL",
			"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WGYJ70A:e4XuptH8D//4crHMMUfL",
			99},
		{"96:aaUi0DTEnLMZMVd2jnEMyFrsdy9LdeGatg3Uogbqs0uBUZoXLn1IvwwDaK:aaf0PU8YMnElrcULdSWgbqs0uBb1IIK",
			"192:aaf6PU8YMnElrcULdSWgbqs0uBb1IIAfsR6OZWjZDx:aaf6PUcYrfLdSWgms0uBb1TA0lZ8ZDx", 80},
		{"48:wX0GLBZET14EHWFIUXs0hPbaL3RdNhI6h0:wPLBS4EecWT6hdNhs",
			"48:w+wNj5GLBX/8jrT14EHWFIUXs0hPbaL3qd9hI6h0:w+zLBX/w14EecWT6ad9hs", 77},
		{"12:7iExTmgeXCcGYX1CRRX1PRRX88p0RRpdV/ISGcEvNOk+l/oX9QUopsAoX9QUopIo:2Ewd+NvN88y3GdkvBC+9lKMHhDh", "12:7iExTmgeXCcGYX1CRRX1PRRXrZGcEvNOk+l/oX9QUopsAoX9QUopIHKl057DRMHD:2Ewd+NvNrgdkvBC+9lKMHhDh", 88},
		{"24:R9mMhMDnWm8m86dmW4zm8mW4zm/mhkcnZ/uLkcHrBCaDrvNQxhwQmq8SywwboX+6:vEnWHH6d/4H/4Z2fvNoF8Sy2yt/YUC",
			"48:xLnWHH6d/4H/4HHHHHHHH4CnrJuN0QhsSyjTU9/j4hbp96khuYhwX:NWHH6dQHQHHHHHHHH4CnV1QeSyj8j4hG",
			43},
	}

	for _, test := range tests {
		var left, right SpamSum
		if _, err := fmt.Sscan(test.left, &left); err != nil {
			t.Errorf("Could not scan string %s, %v", test.left, err)
		}
		if _, err := fmt.Sscan(test.right, &right); err != nil {
			t.Errorf("Could not scan string %s, %v", test.right, err)
		}
		similarity := left.Compare(right)
		if similarity != test.similarity_expected {
			t.Errorf("%s, %s\nSimilariy score should be %d, was %d", left, right, test.similarity_expected, similarity)
		}
	}
}
