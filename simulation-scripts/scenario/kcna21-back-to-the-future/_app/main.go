package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/namsral/flag"
	"github.com/schollz/progressbar/v3"
)

func main() {
	var timeBackwards bool
	var ctfFlag string
	flag.BoolVar(&timeBackwards, "tenet", true, "Turn time around..dnuora emit nruT")
	flag.StringVar(&ctfFlag, "something", "", "Nothing to see here")
	flag.Parse()

	realFlag := decryptFlag(ctfFlag)

	progress(timeBackwards, realFlag)
}

func progress(tenet bool, ctfFlag string) {
	var scriptArray [3]string
	scriptArray[0] = "Enabling time management system"
	scriptArray[1] = "Starting logging systems"
	scriptArray[2] = "State being recorded"

	fmt.Println(scriptArray[0])
	fmt.Println(scriptArray[1])
	fmt.Println(scriptArray[2])

	total := 99
	bar := progressbar.Default(100)
	for i := 0; i < total; i++ {
		bar.Add(1)
		time.Sleep(121 * time.Millisecond)
	}
	if tenet {
		for j := total; j >= 1; j-- {
			bar.Set(j - 1)
			bar.Add(1)
			time.Sleep(121 * time.Millisecond)
		}
		fmt.Println()
		fmt.Println(Reverse(scriptArray[2]))
		fmt.Println(Reverse(scriptArray[1]))
		fmt.Println(Reverse(scriptArray[0]))
	} else {
		bar.Add(1)
		fmt.Println(ctfFlag)
	}
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
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

func decryptFlag(cryptedFlag string) string {
	key := "d85bfa4bcacf4007987b7ecbe32eea6f"

	bytes, _ := hex.DecodeString(cryptedFlag)

	// flag must be shorted than the key
	keySized := key[:len(bytes)]

	// I know this isn't real encryption, but it's good enough
	flag := xor(bytes, []byte(keySized))
	return string(flag)
}
