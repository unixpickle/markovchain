package main

import "math/rand"

type NodeTransition struct {
	Weight int
	Node   *Node
}

type Node struct {
	LastWord    string
	Separator   string
	CurrentWord string
	Transitions []*NodeTransition
	TotalWeight int
}

func NewNode(last, sep, current string) *Node {
	return &Node{
		LastWord:    last,
		Separator:   sep,
		CurrentWord: current,
		Transitions: []*NodeTransition{},
	}
}

func (n *Node) AddTransition(node *Node) {
	n.TotalWeight++
	for _, transition := range n.Transitions {
		if transition.Node == node {
			transition.Weight++
			return
		}
	}
	n.Transitions = append(n.Transitions, &NodeTransition{
		Weight: 1,
		Node:   node,
	})
}

func (n *Node) RandomTransition() *Node {
	if n.TotalWeight == 0 {
		return nil
	}

	num := rand.Intn(n.TotalWeight)
	sum := 0
	for _, transition := range n.Transitions {
		if sum <= num && sum+transition.Weight > num {
			return transition.Node
		}
		sum += transition.Weight
	}
	panic("code should be unreachable")
}
