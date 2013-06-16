// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0

package spamsum

import (
	"hash"
)

type SpamSumWriter struct {
	SpamSum
	spamsumState
}

// StartFixedBlocksize creates a SpamSumWriter with a fixed block size,
// that implements the hash.Hash interface, and accepts an arbitrary
// number of bytes through Write().  Note that the SpamSum algorithm
// does not handle arbitrary length inputs well.  If the input stream
// is significantly longer than SpamLength * blocksize, the tail end
// of the stream will, for most intents and purposes, not generate
// hash blocks.  Please consider the HashBytes or HashReadSeeker
// functions instead.
func StartFixedBlocksize(blockSize uint32) hash.Hash {
	sum := new(SpamSumWriter)

	sum.SpamSum.reset()
	sum.spamsumState.reset()

	sum.blocksize = blockSize
	return sum
}

// Reset sets the state of the SpamSumWriter to its initial value,
// while keeping the blocksize parameter as is.
func (sss *SpamSumWriter) Reset() {
	sss.spamsumState.reset()
	sss.SpamSum.reset()
}

func (sss *SpamSumWriter) Size() int {
	return SpamsumLength
}

// Write a byte slice to the SpamSumWriter.  Returns the length of the
// byte slice, and nil.
func (sss *SpamSumWriter) Write(block []byte) (int, error) {
	processBlock(block, len(block), &sss.spamsumState, &sss.SpamSum)
	return len(block), nil
}

func (sss *SpamSumWriter) String() (result string) {
	writeTail(&sss.spamsumState, &sss.SpamSum)
	result = sss.SpamSum.String()
	// writeTail increments leftIndex and rightIndex, 'finishing'
	// the sum.  Since Write() is still allowed, we decrement
	// again to return to an 'unfinished' state.
	sss.leftIndex -= 1
	sss.rightIndex -= 1
	return
}

// Sum is implemented mostly for the sake of compatibility with
// hash.Hash.  While the SpamSum algorithm creates variable-length
// hashes, Sum is supposed to return a fixed-length slice of Size()
// bytes.  The implementation returns a slice where the non-zero bytes
// contain a base64-encoded 6-bit hash for a `BlockSize()`-sized
// block.  The block hashes continue up to the end of the slice, or up
// to the first zero byte.
func (sss *SpamSumWriter) Sum(block []byte) (result []byte) {
	var cloneState spamsumState = sss.spamsumState
	var cloneSum SpamSum = sss.SpamSum

	processBlock(block, len(block), &cloneState, &cloneSum)

	writeTail(&cloneState, &cloneSum)

	result = make([]byte, SpamsumLength)
	copy(result, cloneSum.leftPart[:cloneSum.leftIndex])
	return
}
