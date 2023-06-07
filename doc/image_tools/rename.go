package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {

	h := sha256.Sum256([]byte("qiyihuo1"))

	fmt.Printf("%v", hex.EncodeToString(h[5:15]))

}
