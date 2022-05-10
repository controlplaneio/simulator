package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

func mainz() {
	plain := os.Args[1]

	key := "d85bfa4bcacf4007987b7ecbe32eea6f"

	keySized := key[:len(plain)]

        // I know this isn't real encryption, but it's good enough
        flag := xor([]byte(plain), []byte(keySized))
	fmt.Println(hex.EncodeToString(flag))
}

func xor(b1 []byte, b2 []byte) []byte {
	if len(b1) != len(b2) {
		fmt.Println("Incorrect key length")
		os.Exit(-1)
	}
	outbytes := make([]byte, len(b1))
	for i := 0; i < len(b1); i++ {
		outbytes[i] = b1[i] ^ b2[i]
	}
	return outbytes
}

func main() {
	mainz()
}
