package markovchain

import (
	"math/rand"
	"sort"

	"github.com/unixpickle/splaytree"
)

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

// Sample selects a random target state using the
// target probabilities.
// It returns nil if no targets are present.
func (s *StateTransitions) Sample() State {
	if len(s.Probabilities) == 0 {
		return nil
	}

	num := rand.Float64()
	var sum float64
	for i, prob := range s.Probabilities {
		sum += prob
		if sum >= num {
			return s.Targets[i]
		}
	}

	// Deal with rounding error edge cases.
	return s.Targets[len(s.Targets)-1]
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
	entries splaytree.Tree
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
	res.lookupOrCreate(lastState)
	res.normalizeProbabilities()

	return res
}

// Lookup finds the transition information going out
// of the given state.
// If no entry exists, nil is returned.
func (c *Chain) Lookup(state State) *StateTransitions {
	searchNode := &treeNode{trans: &StateTransitions{State: state}}
	node := c.entries.Root
	for node != nil {
		comp := node.Value.Compare(searchNode)
		if comp == 0 {
			return node.Value.(*treeNode).trans
		} else if comp == -1 {
			node = node.Right
		} else if comp == 1 {
			node = node.Left
		}
	}
	return nil
}

// Iterate iterates through all of the transitions in
// the table (in order of their state).
// If f returns false, iteration terminates early.
func (c *Chain) Iterate(f func(s *StateTransitions) bool) {
	iterateTree(c.entries.Root, f)
}

func (c *Chain) addEntry(oldState, newState State) {
	entry := c.lookupOrCreate(oldState)
	targetIdx := entry.registerTarget(newState)
	entry.Probabilities[targetIdx]++
}

func (c *Chain) normalizeProbabilities() {
	c.Iterate(func(entry *StateTransitions) bool {
		var sum float64
		for _, x := range entry.Probabilities {
			sum += x
		}
		for i := range entry.Probabilities {
			entry.Probabilities[i] /= sum
		}
		return true
	})
}

func (c *Chain) lookupOrCreate(state State) *StateTransitions {
	entry := c.Lookup(state)
	if entry == nil {
		entry = &StateTransitions{State: state}
		c.entries.Insert(&treeNode{trans: entry})
	}
	return entry
}

func searchStates(states []State, state State) int {
	return sort.Search(len(states), func(i int) bool {
		return states[i].Compare(state) != Less
	})
}

func iterateTree(node *splaytree.Node, f func(s *StateTransitions) bool) bool {
	if node == nil {
		return true
	}
	if !iterateTree(node.Left, f) {
		return false
	}
	if !f(node.Value.(*treeNode).trans) {
		return false
	}
	return iterateTree(node.Right, f)
}

type treeNode struct {
	trans *StateTransitions
}

func (t *treeNode) Compare(t1 splaytree.Value) int {
	comparison := t.trans.State.Compare(t1.(*treeNode).trans.State)
	if comparison == Less {
		return -1
	} else if comparison == Greater {
		return 1
	} else {
		return 0
	}
}
