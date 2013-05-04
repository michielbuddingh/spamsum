// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// spamsum uses a modified version of 32-bit FNV-1.  This file
// is a copy-and-paste hackjob of the standard library FNV
// implementation
package spamsum

import (
	"hash"
)

type (
	sum32  uint32
)

const (
	offset32 = 0x28021967
	prime32  = 16777619
)

// New32 returns a new 32-bit FNV-1 hash.Hash.
func newModifiedFNV() hash.Hash32 {
	var s sum32 = offset32
	return &s
}

func (s *sum32) Reset()  { *s = offset32 }

func (s *sum32) Sum32() uint32  { return uint32(*s) }

func (s *sum32) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash *= prime32
		hash ^= sum32(c)
	}
	*s = hash
	return len(data), nil
}


func (s *sum32) Size() int  { return 4 }

func (s *sum32) BlockSize() int  { return 1 }

func (s *sum32) Sum(in []byte) []byte {
	v := uint32(*s)
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
