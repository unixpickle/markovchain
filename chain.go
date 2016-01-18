package main

import "math/rand"

type NodeInfoKey struct {
	Source      string
	Transition  string
	Destination string
}

type Chain struct {
	Nodes          map[NodeInfoKey]*Node
	StartNodes     map[*Node]bool
	StartNodesList []*Node
}

func NewChain() *Chain {
	return &Chain{
		Nodes:          map[NodeInfoKey]*Node{},
		StartNodes:     map[*Node]bool{},
		StartNodesList: []*Node{},
	}
}

func (c *Chain) AddSentence(s Sentence) {
	if len(s) == 0 || len(s[0].Words) == 0 {
		panic("empty sentence or clause")
	}
	transition := ""
	lastWord := ""
	var lastNode *Node
	for _, clause := range s {
		for _, word := range clause.Words {
			node := c.makeNode(NodeInfoKey{
				Source:      lastWord,
				Transition:  transition,
				Destination: word,
			})
			if lastNode != nil {
				lastNode.AddTransition(node)
			} else {
				if !c.StartNodes[node] {
					c.StartNodes[node] = true
					c.StartNodesList = append(c.StartNodesList, node)
				}
			}
			lastNode = node
			lastWord = word
			transition = ""
		}
		transition = clause.Terminator
	}
	lastNode.AddTransition(c.stopNode(lastNode, transition))
}

func (c *Chain) RandomSentence(minLen int) Sentence {
	for {
		res := Sentence{}
		clause := Clause{}
		wordCount := 0
		node := c.randomStart()
		for node != nil {
			wordCount++
			if node.Separator != "" {
				clause.Terminator = node.Separator
				res = append(res, clause)
				clause = Clause{}
			}
			if node.CurrentWord != "" {
				clause.Words = append(clause.Words, node.CurrentWord)
			}
			node = node.RandomTransition()
		}
		if clause.Words != nil {
			res = append(res, clause)
		}
		if wordCount >= minLen {
			return res
		}
	}
}

func (c *Chain) randomStart() *Node {
	idx := rand.Intn(len(c.StartNodesList))
	return c.StartNodesList[idx]
}

func (c *Chain) makeNode(info NodeInfoKey) *Node {
	if n, ok := c.Nodes[info]; ok {
		return n
	}
	node := NewNode(info.Source, info.Transition, info.Destination)
	c.Nodes[info] = node
	return node
}

func (c *Chain) stopNode(lastNode *Node, sep string) *Node {
	return c.makeNode(NodeInfoKey{
		Source:      lastNode.CurrentWord,
		Transition:  sep,
		Destination: "",
	})
}
