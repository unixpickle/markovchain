package markovchain

import (
	"math"
	"testing"
)

type testingState int

func (t testingState) Compare(t1 State) Comparison {
	n1 := t1.(testingState)
	if t < n1 {
		return Less
	} else if t > n1 {
		return Greater
	} else {
		return Equal
	}
}

func TestNewChainChan(t *testing.T) {
	states := make(chan State, 12)
	states <- testingState(2)
	states <- testingState(3)
	states <- testingState(2)
	states <- testingState(2)
	states <- testingState(5)
	states <- testingState(5)
	states <- testingState(4)
	states <- testingState(2)
	states <- testingState(5)
	states <- testingState(7)
	states <- testingState(1)
	states <- testingState(6)
	close(states)

	actual := NewChainChan(states)
	expected := &Chain{
		Entries: []*StateTransitions{
			&StateTransitions{
				State:         testingState(1),
				Targets:       []State{testingState(6)},
				Probabilities: []float64{1},
			},
			&StateTransitions{
				State:         testingState(2),
				Targets:       []State{testingState(2), testingState(3), testingState(5)},
				Probabilities: []float64{0.25, 0.25, 0.5},
			},
			&StateTransitions{
				State:         testingState(3),
				Targets:       []State{testingState(2)},
				Probabilities: []float64{1},
			},
			&StateTransitions{
				State:         testingState(4),
				Targets:       []State{testingState(2)},
				Probabilities: []float64{1},
			},
			&StateTransitions{
				State:         testingState(5),
				Targets:       []State{testingState(4), testingState(5), testingState(7)},
				Probabilities: []float64{1.0 / 3, 1.0 / 3, 1.0 / 3},
			},
			&StateTransitions{
				State: testingState(6),
			},
			&StateTransitions{
				State:         testingState(7),
				Targets:       []State{testingState(1)},
				Probabilities: []float64{1},
			},
		},
	}
	if len(actual.Entries) != len(expected.Entries) {
		t.Errorf("expected %d entries but got %d", len(expected.Entries),
			len(actual.Entries))
		return
	}
	for i, x := range expected.Entries {
		a := actual.Entries[i]
		if x.State != a.State {
			t.Errorf("entry %d: expected state %v but got %v", i, x.State, a.State)
		}
		if len(x.Targets) != len(a.Targets) ||
			len(x.Probabilities) != len(a.Probabilities) {
			t.Errorf("entry %d: expected %d,%d targets,probabilities but got %d,%d",
				i, len(x.Targets), len(x.Probabilities),
				len(a.Targets), len(a.Probabilities))
		} else {
			for j, expTarget := range x.Targets {
				actTarget := a.Targets[j]
				if expTarget != actTarget {
					t.Errorf("entry %d sub %d: expected target %v got %v", i, j,
						expTarget, actTarget)
				}
			}
			for j, expProb := range x.Probabilities {
				actProb := a.Probabilities[j]
				if math.Abs(expProb-actProb) > 1e-5 {
					t.Errorf("entry %d sub %d: expected probability %f got %f", i, j,
						expProb, actProb)
				}
			}
		}
	}
}
