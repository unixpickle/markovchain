// Command textchain generates text by building a Markov
// chain on a text document.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/unixpickle/markovchain"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s text_doc.txt history_size\n", os.Args[0])
		os.Exit(1)
	}

	historySize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid history size:", os.Args[2])
		os.Exit(1)
	}

	body, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read text:", err)
		os.Exit(1)
	}

	fields := strings.Fields(string(body))
	fieldChan := make(chan string, 10)
	go func() {
		for _, field := range fields {
			fieldChan <- field
		}
		close(fieldChan)
	}()

	log.Println("Building chain...")
	chain := markovchain.NewChainText(fieldChan, historySize)
	log.Println("Generating text...")
	state := randomStart(chain)
	for i := 0; i < 100; i++ {
		ts := state.(markovchain.TextState)
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(ts[len(ts)-1])
		state = randomTransition(chain, state)
	}
	fmt.Println("")
}

func randomStart(ch *markovchain.Chain) markovchain.State {
	var allStates []markovchain.State
	ch.Iterate(func(s *markovchain.StateTransitions) bool {
		allStates = append(allStates, s.State)
		return true
	})
	state := allStates[rand.Intn(len(allStates))]

	// Run through the markov chain to land at a more
	// likely state.
	for i := 0; i < 10; i++ {
		newState := randomTransition(ch, state)
		if newState == nil {
			break
		}
		state = newState
	}

	return state
}

func randomTransition(ch *markovchain.Chain, state markovchain.State) markovchain.State {
	entry := ch.Lookup(state)
	if entry == nil || len(entry.Targets) == 0 {
		return nil
	}

	prob := rand.Float64()
	var curProb float64
	for i, x := range entry.Probabilities {
		curProb += x
		if curProb > prob {
			return entry.Targets[i]
		}
	}

	return entry.Targets[len(entry.Targets)-1]
}
