// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0
package spamsum

import (
	"hash"
)

type SpamStreamSum struct {
	SpamSum
	spamsumState
}

func New(blockSize uint32) hash.Hash {
	sum := new(SpamStreamSum)

	sum.SpamSum.reset()
	sum.spamsumState.reset()

	sum.blocksize = blockSize
	return sum
}

func (sss *SpamStreamSum) Reset() {
	sss.spamsumState.reset()
	sss.SpamSum.reset()
}

func (sss *SpamStreamSum) Size() int {
	return SpamsumLength
}

func (sss *SpamStreamSum) Write(block []byte) (int, error) {
	processBlock(block, len(block), &sss.spamsumState, &sss.SpamSum)
	return len(block), nil
}

func (sss *SpamStreamSum) String() (result string) {
	writeTail(&sss.spamsumState, &sss.SpamSum)
	result = sss.SpamSum.String()
	sss.leftIndex -= 1
	sss.rightIndex -= 1
	return
}

func (sss *SpamStreamSum) Sum(block []byte) (result []byte) {
	var cloneState spamsumState = sss.spamsumState
	var cloneSum SpamSum = sss.SpamSum

	processBlock(block, len(block), &cloneState, &cloneSum)

	writeTail(&cloneState, &cloneSum)

	result = make([]byte, SpamsumLength)
	copy(result, cloneSum.leftPart[:cloneSum.leftIndex])
	return
}
