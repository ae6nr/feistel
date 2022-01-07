package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// Reads a file of base64 encoded keys, each on a different line
// limited to maxkeys number of keys
func readKeyFile(fname string) ([][]byte, error) {
	// open file
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	// parse keys
	str := string(data)
	strs := strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")
	var keys [maxkeys][]byte
	if len(strs) > maxkeys {
		return nil, fmt.Errorf("too many keys (%d/%d)", len(strs), maxkeys)
	}

	// decode base64
	for i, s := range strs {
		keys[i], err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		}
	}

	// return keys
	return keys[:len(strs)], nil
}
