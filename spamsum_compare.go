// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0

package spamsum

import (
	"math"
)

const (
	insCost    = 1
	delCost    = 1
	changeCost = 3
)

// Compare two SpamSums, returning a value between 0 and 100.
// This method is currently not bug-for-bug compatible with the
// original spamsum.
func (from SpamSum) Compare(to SpamSum) (similarity uint32) {
	q := float32(from.blocksize) / float32(to.blocksize)
	if q == 1 {
		similarity = uint32(max(
			score(from.leftPart[:from.leftIndex],
				to.leftPart[:to.leftIndex],
				int(from.blocksize)),
			score(from.rightPart[:from.rightIndex],
				to.rightPart[:to.rightIndex],
				int(to.blocksize))))

	} else if q == 2 {
		similarity = uint32(score(
			from.leftPart[:from.leftIndex],
			to.rightPart[:to.rightIndex],
			int(from.blocksize)))
	} else if q == 0.5 {
		similarity = uint32(score(
			from.rightPart[:from.rightIndex],
			to.leftPart[:to.leftIndex],
			int(to.blocksize)))
	} else {
		similarity = 0
	}
	return
}

func score(from, to []byte, blocksize int) (score int) {
	if !hasCommonSubstring(from, to) {
		return 0
	}

	from = eliminateRepetition(from)
	to = eliminateRepetition(to)

	score = editDistance(from, to)

	score *= SpamsumLength
	score /= len(from) + len(to)

	score = (score * 100) / 64

	score = 100 - score

	maxscore := blocksize / minBlockSize * min(len(from), len(to))
	score = min(score, maxscore)

	return score
}

func editDistance(from, to []byte) int {
	// memoize turns a recursive levenshtein function into one that uses an
	// array to cache results.  Uses |from| * |to| ints of memory.
	memoize := func(calculate func(a, b []byte) int) func(a, b []byte) int {
		var memo []int
		ffl, ttl := len(from), len(to)
		memo = make([]int, ffl*ttl)

		return func(from, to []byte) int {
			fl, tl := len(from), len(to)

			if fl == 0 {
				return tl
			}
			if tl == 0 {
				return fl
			}

			index := ((tl - 1) * ffl) + fl - 1
			if memo[index] == 0 {
				memo[index] = calculate(from, to)

			}
			return memo[index]
		}
	}

	var levenshteinRecursive func(from, to []byte) int

	// to see uncached results, just remove the memoize()
	levenshteinRecursive = memoize(func(from, to []byte) (distance int) {
		// This algorithm is not tuned for anything but legibility, complexity
		// is O(|from| * |to|).  The original code has the option of swapping
		// adjacent characters; as far as I can deduce, this is never used due
		// to the cost penalty, so it is omitted here.
		fl, tl := len(from), len(to)

		if fl == 0 {
			return tl
		}
		if tl == 0 {
			return fl
		}

		var cost = changeCost

		if from[fl-1] == to[tl-1] {
			cost = 0
		}

		return min(
			levenshteinRecursive(from[:fl-1], to)+delCost,
			levenshteinRecursive(from, to[:tl-1])+delCost,
			levenshteinRecursive(from[:fl-1], to[:tl-1])+cost)
	})

	return levenshteinRecursive(from, to)
}

// eliminateRepetition reduces sequences of repeating bytes
// longer than 3 bytes to length 3.
func eliminateRepetition(from []byte) (to []byte) {
	to = make([]byte, len(from))
	copy(to, from[:3])

	i, j := 3, 3
	for ; i < len(from); i++ {
		if from[i-3] != from[i] ||
			from[i-2] != from[i] ||
			from[i-1] != from[i] {
			to[j] = from[i]
			j++
		}
	}

	return to[:j]
}

// hasCommonSubstring returns true if the two byte slices
// passed have a common substring of at least seven bytes.
func hasCommonSubstring(seq1, seq2 []byte) (found bool) {
shift_offset:
	for shift := len(seq1) - 7; shift >= 7-len(seq2); shift-- {
		firstbound, secondbound := max(0, shift), max(0, -shift)
		common := 0
		for i, j := firstbound, secondbound; j < len(seq2) && i < len(seq1); i++ {
			if seq1[i] != seq2[j] {
				common = 0
			} else if common == 6 {
				found = true
				break shift_offset
			} else {
				common++
			}
			j++
		}
	}
	return
}

// min returns the minimum of its arguments
func min(args ...int) int {
	min := int(math.MaxInt32)
	for _, m := range args {
		if m < min {
			min = m
		}
	}
	return min
}

// max returns the maximum of its arguments
func max(args ...int) int {
	max := int(-math.MaxInt32)
	for _, m := range args {
		if m > max {
			max = m
		}
	}
	return max
}
