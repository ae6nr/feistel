package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Feistel cipher demo.")
	fmt.Println("By: Bryan Redd")

	var m string // mode (encrypt or decrypt)
	var enc bool // encode=true, decode=false
	var encstr []byte
	var err error
	var keyfname string // name of the key file
	var outfname string // name of the output file
	var infname string  // name of the input file
	var result []byte   // the result of encryption/decryption

	flag.StringVar(&m, "m", "", "Indicate whether to (enc)rypt or (dec)rypt.")
	flag.StringVar(&keyfname, "k", "", "Indicates filename of file containing base64 keys.")
	flag.StringVar(&outfname, "o", "", "Indicates filename of output file.")
	flag.StringVar(&infname, "i", "", "Indicates filename of input file.")
	flag.Parse()

	// read input
	if infname == "" {
		fmt.Println("Specify an input filename using -i flag.")
		return
	}
	fmt.Printf("Reading %s...\n", infname)
	encstr, err = os.ReadFile(infname)
	check(err)

	if outfname == "" {
		fmt.Println("Specify an output filename using -o flag.")
		return
	}

	// read keys
	if keyfname == "" {
		fmt.Println("Specify a key filename using -k flag.")
		return
	}
	keys, err := readKeyFile(keyfname) // any base64 encoded file, each key is on new line
	check(err)
	fmt.Printf("Using keys from %s...\n", keyfname)

	// determine mode
	if m == "enc" {
		enc = true
		entropy := make([]byte, nentbytes) // add entropy so message isn't the same every time
		_, err = rand.Read(entropy)
		check(err)
		encstr = append(entropy[:], encstr...)
		fmt.Println("Encrypting...")
	} else if m == "dec" {
		enc = false
		fmt.Println("Decrypting...")
	} else {
		fmt.Println("use -m flag to indicate (enc)ryption or (dec)ryption.")
		return
	}

	// determine number of blocks needed
	nblks := len(encstr)/blksize + 1 // number of blocks needed
	if len(encstr)%blksize == 0 {
		nblks--
	}

	var state [2 * hsize]byte

	// encrypt/decrypt
	for i := 0; i < nblks; i++ {

		var in [2 * hsize]byte

		// create input
		n := (i + 1) * blksize
		m := blksize
		if n > len(encstr) {
			n = len(encstr)
			m = len(encstr) % blksize
		}
		copy(in[:m], encstr[i*blksize:n])

		if enc { // encrypt

			// XOR with previous ciphertext so no two blocks are same
			for j := 0; j < blksize; j++ {
				in[j] ^= state[j]
			}

			// encrypt using Feistel cipher
			out := feistel(in, keys[:], enc)
			result = append(result, out[:]...)

			// update state to be previous ciphertext
			state = out

		} else { // decrypt

			// decrypt using Feistel cipher
			out := feistel(in, keys[:], enc)

			// get original plaintext
			for j := 0; j < blksize; j++ {
				out[j] ^= state[j]
			}
			result = append(result, out[:]...)

			// update state to be ciphertext
			state = in
		}

	}

	// save or display results
	if enc {
		os.WriteFile(outfname, result[:], 0644)
	} else {
		os.WriteFile(outfname, bytes.Trim(result[nentbytes:], "\x00"), 0644)
	}

	fmt.Printf("Result saved to %s.\n", outfname)
}
