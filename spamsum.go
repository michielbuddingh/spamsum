// Copyright 2013, Michiel Buddingh, All rights reserved.  Use of this
// code is governed by version 2.0 or later of the Apache License,
// available at http://www.apache.org/licenses/LICENSE-2.0

// Package spamsum implements Andrew Tridgell's fuzzy hash for spam
// detection.  The output should be identical to that produced by
// existing tools like ssdeep.
package spamsum

import (
	"bytes"
	"fmt"
	"io"
)

const (
	rollingWindow = 7
	minBlockSize  = 3
	SpamsumLength = 64
	ReadSize      = 8192
	offset32      = uint32(0x28021967)
	prime32       = uint32(16777619)
)

type SpamSum struct {
	blocksize             uint32
	leftPart              [SpamsumLength]byte
	rightPart             [SpamsumLength / 2]byte
	leftIndex, rightIndex int
}

// String produces the canonical representation of a spamsum. a
// positive number indicating the block size, up to 64 base64
// characters each encoding a 6-bit hash of an (approximately)
// BlockSize()-sized block, up to 32 base64 characters encoding a
// 6-bit hash for BlockSize() * 2.  The three fields are separated by
// a colon.
func (ss *SpamSum) String() string {
	return fmt.Sprintf("%d:%s:%s",
		ss.blocksize,
		string(ss.leftPart[:ss.leftIndex]),
		string(ss.rightPart[:ss.rightIndex]))
}

// BlockSize returns the approximate block size used in this sum.
// Note that this size is only the expected value.  The individual
// 6-bit block hashes may encode far smaller or far larger blocks.
func (ss *SpamSum) BlockSize() int {
	return int(ss.blocksize)
}

// HashBytes takes a byte slice, and takes its SpamSum, calculating
// the optimal block size in several passes.  Since adding more data
// to such a sum would invalidate the block size calculation, this
// SpamSum can not be added to.
func HashBytes(b []byte) (*SpamSum, error) {
	wrapper := io.NewSectionReader(bytes.NewReader(b), 0, int64(len(b)))
	return HashReadSeeker(wrapper, wrapper.Size())
}

// HashReadSeeker requires an implementation of io.ReadSeeker, and a length
// value indicating its size, and takes its SpamSum, calculating
// the optimal block size in several passes.  It is assumed that Seeks upto
// the specified length are allowed. Since adding more data
// to such a sum would invalidate the block size calculation, this
// SpamSum can not be added to.  Any errors returned will originate
// from the implementation of ReadSeeker.
func HashReadSeeker(source io.ReadSeeker, length int64) (*SpamSum, error) {
	sum := new(SpamSum)
	sum.blocksize = minBlockSize

	for int64(sum.blocksize*SpamsumLength) < length {
		sum.blocksize *= 2
	}

	sss := spamsumState{}
source_iteration:
	for {
		sss.reset()
		sum.reset()

		if _, err := source.Seek(0, 0); err != nil {
			return nil, err
		}
		block := make([]byte, ReadSize)

	block_read_loop:
		for {
			var num int
			var err error
			if num, err = source.Read(block); num == 0 {
				break block_read_loop
			} else {
				processBlock(block, num, &sss, sum)
			}

			if err != nil {
				return nil, err
			}
		}

		writeTail(&sss, sum)

		if sum.blocksize > minBlockSize && (sum.leftIndex-1) < (SpamsumLength/2) {
			sum.blocksize /= 2
		} else {
			break source_iteration
		}
	}

	return sum, nil
}

type spamsumState struct {
	// fields for the rolling hash
	window                              [rollingWindow]byte
	rollingSum, h2, shiftHash, position uint32

	// FNV-1 style hash fields
	left, right uint32
}

const b64 string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func processBlock(block []byte, length int, sss *spamsumState, sum *SpamSum) {
	for i := 0; i < length; i++ {
		sss.h2 -= sss.rollingSum
		sss.h2 += rollingWindow * uint32(block[i])

		sss.rollingSum += uint32(block[i])
		sss.rollingSum -= uint32(sss.window[sss.position%rollingWindow])

		sss.window[sss.position%rollingWindow] = block[i]
		sss.position += 1

		sss.shiftHash <<= 5
		sss.shiftHash ^= uint32(block[i])

		roll := sss.rollingSum + sss.h2 + sss.shiftHash

		// left and right are Fowler/No/Voll-1 hashes with a
		// slightly different starting value.
		sss.left *= prime32
		sss.left ^= uint32(block[i])

		sss.right *= prime32
		sss.right ^= uint32(block[i])

		// Assuming the output of the rolling sum is uniformly
		// distributed, this condition will occur once every
		// blocksize bytes.  This means that the expected value
		// for the length of the blocks hashed is blocksize.
		if roll%sum.blocksize == (sum.blocksize - 1) {
			sum.leftPart[sum.leftIndex] = b64[sss.left%64]
			// Note that this means that the first 63 bytes of the
			// hash will encode the first 63*blocksize blocks,
			// and the last byte will encode the remainder, be it
			// one block, or 4GB.
			if sum.leftIndex < SpamsumLength-1 {
				sum.leftIndex += 1
				sss.left = offset32
			}
		}

		// As for the previous condition, but for blocksize * 2
		if roll%(sum.blocksize*2) == ((sum.blocksize * 2) - 1) {
			sum.rightPart[sum.rightIndex] = b64[sss.right%64]
			if sum.rightIndex < (SpamsumLength/2)-1 {
				sum.rightIndex += 1
				sss.right = offset32
			}
		}
	}
}

func writeTail(sss *spamsumState, sum *SpamSum) {
	roll := sss.rollingSum + sss.h2 + sss.shiftHash
	if roll != 0 {
		sum.leftPart[sum.leftIndex] = b64[sss.left%64]
		sum.rightPart[sum.rightIndex] = b64[sss.right%64]
		sum.leftIndex += 1
		sum.rightIndex += 1
	}
}

func (sss *spamsumState) reset() {
	for i := range sss.window {
		sss.window[i] = 0
	}

	sss.rollingSum = 0
	sss.h2 = 0
	sss.shiftHash = 0
	sss.position = 0

	sss.left = offset32
	sss.right = offset32
}

func (sum *SpamSum) reset() {
	for i := range sum.leftPart {
		sum.leftPart[i] = 0
	}

	for i := range sum.rightPart {
		sum.rightPart[i] = 0
	}

	sum.leftIndex, sum.rightIndex = 0, 0
}
