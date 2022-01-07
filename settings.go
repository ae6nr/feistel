package main

import (
	"crypto"
	_ "crypto/sha256"
	"hash"

	_ "golang.org/x/crypto/ripemd160"
	_ "golang.org/x/crypto/sha3"
)

const maxkeys = 16   // sets the maximum number of keys allowed
const nentbytes = 16 // number of bytes entropy to prepend to messages

var hfunc func() hash.Hash = crypto.SHA256.New // import _ "crypto/sha256"
const hsize int = 32                           // hash size in bytes
// var hfunc func() hash.Hash = crypto.RIPEMD160.New // import _ "golang.org/x/crypto/ripemd160"
// const hsize int = 20                              // hash size in bytes
// var hfunc func() hash.Hash = crypto.SHA3_512.New // import _ "golang.org/x/crypto/sha3"
// const hsize int = 64                             // hash size in bytes

const blksize int = 2 * hsize
