package main

import "crypto/hmac"

// Implements a single round of the Feistel cipher
func round(in [2 * hsize]byte, key []byte) (out [2 * hsize]byte) {

	// split input in half
	L := in[0:hsize]         // left
	R := in[hsize : 2*hsize] // right
	copy(out[0:hsize], R)    // copy input right to output left

	// nonlinear function
	h := hmac.New(hfunc, key) // hash function
	h.Write(R)
	x := h.Sum([]byte{}) // digest
	copy(out[hsize:2*hsize], L)
	for i, b := range x { // XOR bytes
		out[hsize+i] ^= b
	}

	return out
}

// Feistel cipher
// The input is twice the HMAC digest size
// keys contains an array of keys used as keys for the HMAC
// enc=true will encode, enc=false will decode
// number of rounds is determined by number of keys supplied
func feistel(in [2 * hsize]byte, keys [][]byte, enc bool) (out [2 * hsize]byte) {

	// Feistel rounds
	var idx int
	for i := 0; i < len(keys); i++ {
		if enc {
			idx = i
		} else {
			idx = len(keys) - 1 - i
		}
		in = round(in, keys[idx])
	}

	// final reversal
	copy(out[0:hsize], in[hsize:2*hsize])
	copy(out[hsize:2*hsize], in[0:hsize])
	return out
}
