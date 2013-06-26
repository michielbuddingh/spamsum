// Copyright 2013, Michiel Buddingh, All rights reserved.  Use of this
// code is governed by version 2.0 or later of the Apache License,
// available at http://www.apache.org/licenses/LICENSE-2.0

// This is a really ugly script that turns any git repository into a
// test suite for the Go spamsum implementation.  It inspects all
// revisions of all files in the current head, takes their Spamsum,
// and compares this against the result of the original spamsum tool.
//
// Note that this script requires the C sources of the original
// spamsum tool to work, and suffers from bad error handling in many
// places.  The time to run scales quadratically with the number of
// revisions, so its best suited for mid-sized git repositories.
package main

import (
	"bytes"
	"fmt"
	"github.com/michielbuddingh/spamsum"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const spamsumpath = "./spamsum"
const spamsum_comparepath = "./spamsum_compare"

var count = 0
var comparisoncount = 0

func main() {
	files, err := exec.Command("git", "ls-files").Output()
	if err != nil {
		println(err)
	}
	filelist := strings.Split(string(files), "\n")
	for _, file := range filelist {
		// Iterate over all revisions of a file.
		allRevisions(file)
	}
	log.Printf("%d files processed\n", count)
	log.Printf("%d comparisons\n", comparisoncount)
}

func allRevisions(filename string) {
	grabCommit := regexp.MustCompile("commit ([0-9a-f]*)")
	commitlog, _ := exec.Command("git", "log", "--pretty=short", filename).Output()
	commits := grabCommit.FindAllSubmatch(commitlog, -1)

	sums := make([]string, 0)

	if len(commits) > 1 {
		for _, commit := range commits {
			contents, err := exec.Command("git", "show", "--format=raw", string(commit[1])+":"+filename).Output()

			if err == nil && len(contents) > 0 && !strings.HasPrefix(string(contents), "fatal") {
				sum1 := createSpamSum(contents)
				sum2 := createOriginalSpamSum(contents)
				if sum1 != sum2 {
					log.Printf("revision %s of file %s has differing spamsums", string(commit[1]), filename)
				}
				count++
				sums = append(sums, sum1)
			}
		}
	}

	for idx, left := range sums {
		for i := idx + 1; i < len(sums); i++ {
			first := compareSpamSum(left, sums[i])
			second := compareOriginalSpamSum(left, sums[i])
			if first != second {
				log.Printf("Difference in comparison between %s and %s, %d, %d\n", left, sums[i], first, second)
			}
			comparisoncount++
		}
	}

}

func compareSpamSum(left, right string) int {
	var leftSum, rightSum spamsum.SpamSum
	fmt.Sscan(left, &leftSum)
	fmt.Sscan(right, &rightSum)
	score := leftSum.Compare(rightSum)
	return int(score)
}

func compareOriginalSpamSum(left, right string) int {
	scoretext, _ := exec.Command(spamsum_comparepath, left, right).Output()
	var score int
	score, _ = strconv.Atoi(string(scoretext))
	return score
}

func createSpamSum(contents []byte) string {
	reader := bytes.NewReader(contents)
	sum, _ := spamsum.HashReadSeeker(reader, int64(len(contents)))
	return sum.String()
}

func createOriginalSpamSum(contents []byte) string {
	reader := bytes.NewReader(contents)
	cmd := exec.Command(spamsumpath, "-")
	cmd.Stdin = reader
	if sumbytes, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(sumbytes))
	}
	return "nil"
}
