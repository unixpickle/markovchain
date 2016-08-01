package markovchain

import "sort"

// A Comparison represents the ordinal relationship
// between two comparable things.
type Comparison int

const (
	Equal Comparison = iota
	Less
	Greater
)

// A State represents a state in a Markov chain.
// Since states are stored in transition Chains, they
// are comparable to facilitate efficient lookups.
type State interface {
	// Compare the receiver to s1.
	// If the receiver is greater than s1, this returns
	// Greater, etc.
	Compare(s1 State) Comparison
}

// StateTransitions stores the transitions going out
// of a given state in a Markov chain.
//
// A StateTransitions instance is equivalent to a
// row/column in a markov matrix.
type StateTransitions struct {
	// State is the state whose outgoing transitions are
	// stored in this object.
	State State

	// Targets is the list of possible states which could
	// result from this state.
	Targets []State

	// Probabilities stores the transition probability for
	// each target in Targets.
	Probabilities []float64
}

func (s *StateTransitions) registerTarget(state State) int {
	idx := searchStates(s.Targets, state)
	if idx == len(s.Targets) {
		s.Targets = append(s.Targets, state)
		s.Probabilities = append(s.Probabilities, 0)
	} else if s.Targets[idx].Compare(state) != Equal {
		s.Targets = append(s.Targets, nil)
		s.Probabilities = append(s.Probabilities, 0)
		copy(s.Targets[idx+1:], s.Targets[idx:])
		copy(s.Probabilities[idx+1:], s.Probabilities[idx:])
		s.Probabilities[idx] = 0
		s.Targets[idx] = state
	}
	return idx
}

// A Chain stores the transition probabilities between
// states in a Markov process.
//
// A Chain is equivalent to a Markov matrix.
type Chain struct {
	Entries []*StateTransitions
}

// NewChainChan creates a Chain by reading a channel of
// states and calculating transition probabilities based
// on that channel.
func NewChainChan(ch <-chan State) *Chain {
	res := &Chain{}

	var lastState State
	for state := range ch {
		if lastState != nil {
			res.addEntry(lastState, state)
		}
		lastState = state
	}
	res.registerState(lastState)
	res.normalizeProbabilities()

	return res
}

// Lookup finds the transition information going out
// of the given state.
// If no entry exists, nil is returned.
func (c *Chain) Lookup(state State) *StateTransitions {
	idx := c.searchEntries(state)
	if idx == len(c.Entries) {
		return nil
	} else if c.Entries[idx].State.Compare(state) == Equal {
		return c.Entries[idx]
	} else {
		return nil
	}
}

func (c *Chain) addEntry(oldState, newState State) {
	entry := c.registerState(oldState)
	targetIdx := entry.registerTarget(newState)
	entry.Probabilities[targetIdx]++
}

func (c *Chain) normalizeProbabilities() {
	for _, entry := range c.Entries {
		var sum float64
		for _, x := range entry.Probabilities {
			sum += x
		}
		for i := range entry.Probabilities {
			entry.Probabilities[i] /= sum
		}
	}
}

func (c *Chain) registerState(state State) *StateTransitions {
	idx := c.searchEntries(state)
	if idx == len(c.Entries) {
		c.Entries = append(c.Entries, &StateTransitions{State: state})
	} else if c.Entries[idx].State.Compare(state) != Equal {
		c.Entries = append(c.Entries, nil)
		copy(c.Entries[idx+1:], c.Entries[idx:])
		c.Entries[idx] = &StateTransitions{State: state}
	}
	return c.Entries[idx]
}

func (c *Chain) searchEntries(state State) int {
	return sort.Search(len(c.Entries), func(i int) bool {
		return c.Entries[i].State.Compare(state) != Less
	})
}

func searchStates(states []State, state State) int {
	return sort.Search(len(states), func(i int) bool {
		return states[i].Compare(state) != Less
	})
}
