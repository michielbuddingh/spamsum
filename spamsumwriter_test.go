// Copyright 2013, Michiel Buddingh, All rights reserved.
// Use of this code is governed by version 2.0 or later of the Apache
// License, available at http://www.apache.org/licenses/LICENSE-2.0

package spamsum

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		seed      int64
		length    int
		blocksize uint32
		expected  string
	}{
		{42, 16384, 384, "384:PnwCSZ6yE9r4UCZ1he34xas/E8AhHgdd2yM:PbSZ6yE9rGfExx"},
		{1000, 2048, 48, "48:Zo+v/bCSly4VhreHwHJdkHTzF7sjBU1YuD/QtFsByxoSJW+QiLlH:uSWSFteQHJd+Tp79mqSqyCt+5LlH"},
		{1000, 1048576, 24576, "24576:xL2L/P40/cnWGr7tsP+mgdQGvnb1UV+gQ8ZwU:ErPP/2WItsPTgdD/bqQ4yU"},
		{71268, 24, 3, "3:N0n6xmcFctn:7xmptn"},
	}

	for _, test := range tests {
		generator := rand.New(rand.NewSource(test.seed))
		writer := StartFixedBlocksize(test.blocksize)
		for i := 0; i < test.length/4; i++ {
			binary.Write(writer, binary.BigEndian, generator.Uint32())
		}
		if writer.String() != test.expected {
			t.Errorf("Expected %v, result was %v", test.expected, writer)
		}
	}
}

func TestWriterIntermediate(t *testing.T) {
	generator := rand.New(rand.NewSource(3181))
	writer := StartFixedBlocksize(768)

	for i := 0; i < 4096; i++ {
		binary.Write(writer, binary.BigEndian, generator.Uint32())
	}

	expectedIntermediate := "768:Mz4Rjllf5YbJAdWsgdL7rPjcYK1TkTN:MzsjjfqbJkzcLPjcYEW"
	actualIntermediate := writer.String()

	if actualIntermediate != expectedIntermediate {
		t.Errorf("Expected %v, actual %v", expectedIntermediate, actualIntermediate)
	}

	for i := 0; i < 4096; i++ {
		binary.Write(writer, binary.BigEndian, generator.Uint32())
	}

	expectedFinal := "768:Mz4Rjllf5YbJAdWsgdL7rPjcYK1TkTQSPNZHrHmwGS8VUSYy20b9n6TZd:MzsjjfqbJkzcLPjcYEcNpXG7VUSR2Q6n"
	actualFinal := writer.String()

	if actualFinal != expectedFinal {
		t.Errorf("Expected %v, actual %v", expectedFinal, actualFinal)
	}
}

func TestWriterReset(t *testing.T) {
	generator := rand.New(rand.NewSource(3181))
	writer := StartFixedBlocksize(768)
	emtpySlice := make([]byte, 0)

	for i := 0; i < 4096; i++ {
		binary.Write(writer, binary.BigEndian, generator.Uint32())
	}

	beforeReset := writer.String()
	beforeResetBinary := writer.Sum(emtpySlice)

	writer.Reset()
	generator = rand.New(rand.NewSource(3181))

	for i := 0; i < 4096; i++ {
		binary.Write(writer, binary.BigEndian, generator.Uint32())
	}

	afterReset := writer.String()
	afterResetBinary := writer.Sum(emtpySlice)

	if beforeReset != afterReset {
		t.Errorf("Same data written to the same writer, but different results!")
	}

	if len(afterResetBinary) != len(beforeResetBinary) {
		t.Errorf("Binary spamsums are not even the same size")
	}

	for i, _ := range beforeResetBinary {
		if beforeResetBinary[i] != afterResetBinary[i] {
			t.Errorf("Binary spamsums before and after reset differ at byte %d", i)
			break
		}
	}
}

func TestSize(t *testing.T) {
	writer := StartFixedBlocksize(16)
	if writer.Size() != SpamsumLength {
		t.Errorf("Max result size should always be equal to SpamsumLength, which is %d\n", SpamsumLength)
	}
}
