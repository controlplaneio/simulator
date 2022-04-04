package main

import (
	"fmt"
	"time"

	"github.com/namsral/flag"
	"github.com/schollz/progressbar/v3"
)

func main() {
	var scriptArray [3]string
	scriptArray[0] = "Enabling time management system"
	scriptArray[1] = "Starting logging systems"
	scriptArray[2] = "State being recorded"
	var tenet bool
	flag.BoolVar(&tenet, "tenet", false, "Turn time around..dnuora emit nruT")
	flag.Parse()

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
		fmt.Println(Reverse(scriptArray[2]))
		fmt.Println(Reverse(scriptArray[1]))
		fmt.Println(Reverse(scriptArray[0]))
	} else {
		bar.Add(1)
		fmt.Println("flag")
	}
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
