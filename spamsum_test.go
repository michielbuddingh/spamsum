package spamsum

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	tests := []struct {
		input      string
		shouldfail bool
	}{
		{"49152:dihMNzhZt62oh9+onrqMPr/KwJsvD/mMplt:Hxxpj", false},
		{"12582912:kVxeXup8VuH8rD//4crHBrlGXm5WgYJ70A:e4XuptH8D//4crHMmUfL", false},
		{"18446744073709551616:dihMNzhZt62oh9+onrqMPr/KwJsvD/mMplt:H.soa", true},
		{"49152:dihMNzhZt62oh9+onrqMPr/KwJsvD/mMplt.Hxxpj", true},
		{"22:i3wkMEgPthpID7YoQDjrdAjGBwBIg8Qow0iLSAhIi3AQSItCCEiLUhBIOch1MEiJBCRIiVQkCEiJ:UxUp", true},
	}

	for _, test := range tests {
		var sum SpamSum
		_, err := fmt.Sscan(test.input, &sum)
		if test.shouldfail && err == nil {
			t.Errorf("Should not be able to parse %s\n", test.input)
		}
		if !test.shouldfail && err != nil {
			t.Errorf("Parse failed with error: %v", err)
		}
		if !test.shouldfail && sum.String() != test.input {
			t.Errorf("scanned sum %s is not equal to input string %s\n", sum.String(), test.input)
		}
	}
}

func TestHashReadSeeker(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"LAND.MAP", "768:tlBecdq6/+dgZUTp+gAdA3T9Y02xEFshHOl3O98FzbXfBfhPcGxGB3whvm9HvMB1:O"},
		{"embedded_video_quicktime.doc", "192:o50PBwxGc+ZrnCe9pz1aZ8GHiLUd0935:G8cOz9pzJ3"},
	}

	for _, test := range tests {
		path := filepath.Join("testdata", test.filename)
		file, openerr := os.Open(path)
		if openerr != nil {
			t.Fatal(openerr)
		}
		defer file.Close()
		stat, staterr := file.Stat()
		if staterr != nil {
			t.Fatal(staterr)
		}

		sum, sumerr := HashReadSeeker(file, stat.Size())
		if sumerr != nil {
			t.Fatal(sumerr)
		}

		if sum.String() != test.expected {
			t.Errorf("Expected %s hashing %s, result was %v", test.expected, test.filename, sum)
		}
	}
}

func TestHashBytes(t *testing.T) {
	tests := []struct {
		seed      int64
		length    int
		blocksize uint32
		expected  string
	}{
		{6065, 1024, 24, "24:D4JsKhbN85qJzgs+JLY4DffT9hhD6Wa333cRhDEPVreO:LKFN85qJMHJj1hkuDEPVeO"},
		{1029936, 1025, 12, "12:RePpJA8PW0JP1uTMCa9qpRCwtnacOzIayUpkmp6v12qVVIFSpNKDrASjqPoOaY1L:UPPWE4TMY37nYzslTN2gTDZpaCji8ZmM"},
		{1252877, 22624, 192, "192:V5cZcnyOVaMvLF4f8mkfu4u95tgALGPVxn8QhXSd1CsvQ+D3QMfFiz/uxuVge/7P:nIyvGkWN/iHImSc6vAzWgyeIiBzPbgzk"},
		{1497046, 22624, 192, "192:BTBLFZFxOyNbTjMRjkLOBKiDKe2cRfzKQACMTsZGJRaYjx44gkX2iJ4nURozFp9S:B5OyN/QjkL6KiH2JgoNIDreMuxqRkxJZ"},
	}

	for _, test := range tests {
		byteSlice := make([]byte, test.length)
		generator := rand.New(rand.NewSource(test.seed))

		generator.Read(byteSlice)

		sum := HashBytes(byteSlice)

		if sum.String() != test.expected {
			t.Errorf("Expected %v, result was %v", test.expected, sum)
		}
	}
}

func TestBlocksizeAdjustment(t *testing.T) {
	byteSlice := make([]byte, 17921)
	generator := rand.New(rand.NewSource(191))

	i := 0
	for ; i < 24; i++ {
		binary.BigEndian.PutUint32(byteSlice[i*4:], generator.Uint32())
	}

	for ; i < 17921; i++ {
		byteSlice[i] = 0
	}

	sum := HashBytes(byteSlice)
	expected := "3:Bl5KOiWl/:ldZ/"

	if sum.String() != expected {
		t.Errorf("Expected %v, result was %v", expected, sum)
	}
}
