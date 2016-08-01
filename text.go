package markovchain

import "strings"

// A TextState is a comparable list of textual tokens.
type TextState []string

// Compare alphabetically compares t to s1 token by token,
// returning the comparison result of the first token to
// disagree.
// Both states must be TextTokens of the same length.
func (t TextState) Compare(s1 State) Comparison {
	t1 := s1.(TextState)
	if len(t1) != len(t) {
		panic("can only compare equal-length text tokens")
	}
	for i, x := range t {
		res := strings.Compare(x, t1[i])
		if res == -1 {
			return Less
		} else if res == 1 {
			return Greater
		}
	}
	return Equal
}

// NewChainText creates a Markov chain for transitions
// between TextStates of length n.
// It takes a channel of words from which transition
// probabilities can be derived.
//
// In other words, it creates a Markov chain which
// predicts the following word based on the previous
// n words in the string.
func NewChainText(words <-chan string, n int) *Chain {
	stateChan := make(chan State)
	go func() {
		var curState TextState
		for word := range words {
			if len(curState) < n {
				curState = append(curState, word)
				if len(curState) == n {
					stateChan <- curState
				}
			} else {
				newState := make(TextState, n)
				copy(newState, curState[1:])
				newState[n-1] = word
				curState = newState
				stateChan <- curState
			}
		}
		close(stateChan)
	}()
	return NewChainChan(stateChan)
}
