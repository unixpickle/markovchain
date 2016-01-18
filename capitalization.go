package main

type CapCounts struct {
	AllCaps  int
	SomeCaps int
	NoCaps   int
}

type CapTracker struct {
	CountMap map[string]*CapCounts
}

func NewCapTracker() *CapTracker {
	return &CapTracker{CountMap: map[string]*CapCounts{}}
}

func (c *CapTracker) AddSentence(s Sentence) {
	for i, clause := range s {
		for j, word := range clause.Words {
			if i == 0 && j == 0 && word.Capitalization == SomeCapital {
				continue
			}
			entry, ok := c.CountMap[word.Text]
			if !ok {
				entry = &CapCounts{}
				c.CountMap[word.Text] = entry
			}
			switch word.Capitalization {
			case AllCapital:
				entry.AllCaps++
			case SomeCapital:
				entry.SomeCaps++
			case NoCapital:
				entry.NoCaps++
			}
		}
	}
}

func (c *CapTracker) FixSentence(s Sentence) Sentence {
	res := Sentence{}
	for i, clause := range s {
		newClause := Clause{Terminator: clause.Terminator}
		for j, word := range clause.Words {
			caps := c.capsForWord(word.Text)
			if caps == NoCapital && i == 0 && j == 0 {
				caps = SomeCapital
			}
			newWord := Word{Text: word.Text, Capitalization: caps}
			newClause.Words = append(newClause.Words, newWord)
		}
		res = append(res, newClause)
	}
	return res
}

func (c *CapTracker) capsForWord(word string) Capitalization {
	if entry, ok := c.CountMap[word]; ok {
		if entry.AllCaps > entry.SomeCaps && entry.AllCaps > entry.NoCaps {
			return AllCapital
		} else if entry.SomeCaps > entry.AllCaps && entry.SomeCaps > entry.NoCaps {
			return SomeCapital
		} else {
			return NoCapital
		}
	}
	return NoCapital
}
