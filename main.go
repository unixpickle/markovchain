package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: markovchain <input.txt>")
		os.Exit(1)
	}
	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	sentences := TokenizeText(string(contents))
	chain := NewChain()
	for _, sentence := range sentences {
		chain.AddSentence(sentence)
	}

	for i := 0; i < 15; i++ {
		fmt.Println(chain.RandomSentence(5))
	}
}
