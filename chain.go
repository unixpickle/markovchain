package main

import "math/rand"

type Chain struct {
	Nodes map[string]*Node
	Words []string
}

func NewChain() *Chain {
	return &Chain{Nodes: map[string]*Node{}, Words: []string{}}
}

func (c *Chain) AddWord(word string) *Node {
	if node, ok := c.Nodes[word]; ok {
		return node
	}
	if word != "," && word != "?" && word != ";" && word != "!" {
		c.Words = append(c.Words, word)
	}
	node := NewNode(word)
	c.Nodes[word] = node
	return node
}

func (c *Chain) RandomNode() *Node {
	if len(c.Nodes) == 0 {
		panic("no nodes to randomly pick from")
	}
	idx := rand.Intn(len(c.Words))
	word := c.Words[idx]
	return c.Nodes[word]
}

func (c *Chain) RandomSequence() []string {
	res := make([]string, 0)
	n := c.RandomNode()
	for {
		res = append(res, n.Word)
		if n.Word == "." || n.Word == "?" || n.Word == "!" {
			break
		}
		next := n.RandomTarget()
		n = c.Nodes[next]
	}
	return res
}

type Node struct {
	Word        string
	Targets     map[string]int
	TotalWeight int
}

func NewNode(word string) *Node {
	return &Node{Word: word, Targets: map[string]int{}}
}

func (n *Node) AddTarget(word string) {
	n.Targets[word] = n.Targets[word] + 1
	n.TotalWeight++
}

func (n *Node) RandomTarget() string {
	// If there are no target words, all we can do is loop.
	if n.TotalWeight == 0 {
		return n.Word
	}

	num := rand.Intn(n.TotalWeight)
	sum := 0
	for word, weight := range n.Targets {
		if sum <= num && sum+weight > num {
			return word
		}
		sum += weight
	}
	panic("code should be unreachable")
}
