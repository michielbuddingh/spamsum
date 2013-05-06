spamsum
=======

#### A fuzzy checksum for matching spam ####

This is a native go implementation of spamsum.

spamsum was developer by Andrew Tridgell to hash email messages for computationally inexpensive SPAM detection. See <http://junkcode.samba.org/#spamsum>.

The state of this package
-------------------------

* Not ready for general consumption
* It seems to generate results identical to that of the [spamsum tool](https://junkcode.samba.org/ftp/unpacked/junkcode/spamsum/) and [ssdeep](http://ssdeep.sf.net).  This has only been tested on a small number of files.
* It is about twice as slow as the spamsum tool; about 40MB/s on a 3Ghz Core i3
* It does not yet support fuzzy comparison, which is the only reason why you'd want to use a fuzzy hash anyway.

How to use
----------

Unfortunately, the default operation for spamsum is to iterate over the data several times to determine an optimal block size, so it's not sensible to implement the `hash.Hash` interface.

Instead, the package exports the functions `HashBytes(b [] byte)` and `HashReadSeeker(source io.ReadSeeker, length int64)`.

	if file, err := os.Open("filename"); err != nil {
		log.Fatal(err)
	} else if stat, err := file.Stat(); err != nil {
		log.Fatal(err)
	} else {
		sum, err := spamsum.HashReadSeeker(file, stat.Size())
		// etc.
	}

Any errors returned by `HashReadSeeker` will originate from the `io.ReadSeeker` functions.

### Alternatively ###

If it is acceptable to set a fixed blocksize beforehand, the `SpamStreamSum` type can be used, which _does_ implement the `hash.Hash` interface.  The `Sum(b []byte) []byte` method is not terribly useful; it will return a slice where the non-zero bytes contain a base64-encoded 6-bit hash for a `BlockSize()`-sized block.  Use the `String()` method to obtain a more useful representation.

### License ###

Use of this code is governed by version 2.0 or later of the Apache
License, available at <http://www.apache.org/licenses/LICENSE-2.0>
