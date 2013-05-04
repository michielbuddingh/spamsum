// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0
package spamsum

import (
	"hash"
	"encoding/binary"
)

const (
	rollingWindow = 7
)

type rollState struct {
	window [rollingWindow]byte
	rollingSum, h2, shiftHash, position uint32
}

func newRollHash() hash.Hash32 {
	state := new(rollState)
	state.Reset()
	return state
}

func (rs * rollState) Write(p []byte) (n int, err error) {
	for _, b := range(p) {
		rs.h2 -= rs.rollingSum
		rs.h2 += rollingWindow * uint32(b)

		rs.rollingSum += uint32(b)
		rs.rollingSum -= uint32(rs.window[rs.position % rollingWindow])

		rs.window[rs.position % rollingWindow] = b
		rs.position += 1

		rs.shiftHash <<= 5
		rs.shiftHash ^= uint32(b)
	}
	return len(p), nil
}

func (rs * rollState) Reset() {
	for i := range(rs.window) { rs.window[i] = 0 }
	rs.rollingSum = 0
	rs.h2 = 0
	rs.shiftHash = 0
	rs.position = 0
}

func (rs * rollState) BlockSize() int {
	return rollingWindow
}

func (rs * rollState) Size() int {
	return 4
}

func (rs * rollState) Sum32() uint32 {
	return rs.rollingSum + rs.h2 + rs.shiftHash
}

func (rs * rollState) Sum(b []byte) (result []byte) {
	temp := rollState{
		rollingSum : rs.rollingSum,
		h2 : rs.h2,
		shiftHash : rs.shiftHash,
		position : rs.position,
	}

	temp.Write(b);
	binary.BigEndian.PutUint32(result, temp.Sum32())
	return
}
