package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
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
	words := extractWords(string(contents))

	rand.Seed(time.Now().UnixNano())

	fmt.Println("Generating chain...")
	chain := NewChain()
	var lastNode *Node
	for _, word := range words {
		node := chain.AddWord(word)
		if lastNode != nil {
			lastNode.AddTarget(word)
		}
		lastNode = node
	}

	fmt.Println("Generating sentences")
	for i := 0; i < 15; i++ {
		sequence := chain.RandomSequence(7)
		fmt.Println(wordsToSentence(sequence))
	}
}

func extractWords(contents string) []string {
	split := strings.Fields(contents)
	res := make([]string, len(split))
	for _, s := range split {
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		if strings.HasPrefix(s, "(") {
			s = s[1:]
		}
		if strings.HasSuffix(s, ")") {
			s = s[:len(s)-1]
		}
		if strings.HasPrefix(s, "\"") || strings.HasPrefix(s, "\u201C") ||
			strings.HasPrefix(s, "\xD3") {
			s = string([]rune(s)[1:])
		}
		if strings.HasSuffix(s, "\"") || strings.HasSuffix(s, "\u201D") ||
			strings.HasPrefix(s, "\xD3") {
			runes := []rune(s)
			s = string(runes[:len(runes)-1])
		}
		var stopper string
		if s != "mr." && s != "ms." && s != "m." && s != "mrs." && s != "dr." {
			if strings.HasSuffix(s, ".") || strings.HasSuffix(s, "!") ||
				strings.HasSuffix(s, "?") || strings.HasSuffix(s, ",") ||
				strings.HasSuffix(s, ";") {
				stopper = string(s[len(s)-1])
				s = s[:len(s)-1]
			}
		} else {
			s = strings.ToUpper(s[:1]) + s[1:]
		}
		if len(s) > 0 {
			res = append(res, s)
			if stopper != "" {
				res = append(res, stopper)
			}
		}
	}
	return res
}

func wordsToSentence(words []string) string {
	res := strings.Join(words, " ")
	if len(res) == 0 {
		return ""
	}
	for _, x := range []string{";", ",", "!", ".", "?"} {
		res = strings.Replace(res, " "+x, x, -1)
	}
	return strings.ToUpper(res[:1]) + res[1:]
}
