// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0
package spamsum

import (
	"io"
	"bytes"
	"fmt"
)

const (
	minBlockSize = 3
	SpamsumLength = 64
)

type Summable interface {
	io.ReadSeeker
	Size() int64;
}

type SpamSum struct {
	blocksize uint32
	leftPart [SpamsumLength]byte
	rightPart [SpamsumLength / 2]byte
}

func (ss * SpamSum) String() string {
	return fmt.Sprintf("%d:%s:%s", ss.blocksize, string(ss.leftPart[:]), string(ss.rightPart[:]))
}

func (ss * SpamSum) BlockSize() int {
	return int(ss.blocksize)
}

func HashBytes (b []byte) (* SpamSum, error) {
	wrapper := io.NewSectionReader(bytes.NewReader(b),0,int64(len(b)))
	return HashReadSeeker(wrapper, wrapper.Size());
}


const b64 string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

func HashReadSeeker (source io.ReadSeeker, length int64) (* SpamSum, error) {
	sum := new(SpamSum)
	sum.blocksize = minBlockSize

	for int64(sum.blocksize * SpamsumLength) < length {
		sum.blocksize *= 2
	}

source_iteration:for {
		rolling := newRollHash()
		left := newModifiedFNV()
		right := newModifiedFNV()
		leftIndex, rightIndex := 0, 0

		if _, err := source.Seek(0, 0); err != nil {
			return nil, err
		}
		block := make([]byte, sum.blocksize)

	block_read_loop: for {
			var num int; var err error
			if num, err = source.Read(block); num == 0 {
				break block_read_loop
			} else {
				for i := 0; i < num; {
					l, r, j := false, false, i

				scan_trigger_condition: for ; j < num; {
						rolling.Write([]byte{ block[j] });

						roll := rolling.Sum32()
						// Trigger condition 1
						if roll % sum.blocksize == (sum.blocksize - 1) {
							l = true
						}

						// Trigger condition 2
						if roll % (sum.blocksize * 2) == ((sum.blocksize * 2) - 1) {
							r = true
						}

						j++

						if (l || r) {
							break scan_trigger_condition
						}

					}

					left.Write(block[i:j])
					right.Write(block[i:j])

					if l {
						sum.leftPart[leftIndex] = b64[left.Sum32() % 64]
						if leftIndex < SpamsumLength - 1 {
							leftIndex += 1
							left = newModifiedFNV()
						}
					}

					if r {
						sum.rightPart[rightIndex] = b64[right.Sum32() % 64]
						if rightIndex < (SpamsumLength/2) - 1 {
							rightIndex += 1
							right = newModifiedFNV()
						}
					}


					i = j
				}

			}

			if err != nil {
				return nil, err
			}
		}

		if rolling.Sum32() != 0 {
			sum.leftPart[leftIndex] = b64[left.Sum32() % 64]
			sum.rightPart[rightIndex] = b64[right.Sum32() % 64]
		}

		if sum.blocksize > minBlockSize && leftIndex < (SpamsumLength / 2) {
			sum.blocksize /= 2
		} else {
			break source_iteration
		}
	}


	return sum, nil
}