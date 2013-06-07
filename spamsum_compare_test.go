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

func TestCompare(t *testing.T) {
	tests := []struct {
		left, right         string
		similarity_expected int
	}{
		// these are not values produced by the original spamsum
		// score algorithm
		{

			"12582912:UVxeXup8VuH8rD//pcrHBrlG5FWgYJ70A:O4XuptH8D//pcrHmgfL",
			"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A:e4XuptH8D//4crHMmUfL",
			85},

		{"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A:e4XuptH8D//4crHMmUfL",
			"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WGYJ70A:e4XuptH8D//4crHMMUfL",
			96},
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
		println(similarity)
	}
}
