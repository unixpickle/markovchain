package main

import "strings"

var SentenceTerminators = []string{".", "?", "!"}
var ClauseTerminators = []string{"--", ",", ";", ":", ")", "\""}
var ClauseIntroducers = []string{"\"", "("}
var Titles = []string{"Mr.", "Mrs.", "Dr.", "Ms.", "M."}

type Capitalization int

const (
	NoCapital Capitalization = iota
	SomeCapital
	AllCapital
)

type Word struct {
	Text           string
	Capitalization Capitalization
}

type Clause struct {
	Terminator string
	Words      []Word
}

type Sentence []Clause

func TokenizeText(text string) []Sentence {
	res := []Sentence{}

	text = normalizeUnicodeSymbols(text)
	text = strings.Replace(text, "--", " -- ", -1)
	tokens := strings.Fields(text)

	var sentence Sentence
	var clause Clause
	for _, token := range tokens {
		bareToken := stripPunctuation(token)
		word := Word{strings.ToLower(bareToken), wordCapitalization(bareToken)}
		if introducesClause(bareToken) && clause.Words != nil {
			// TODO: use terminators for introductions.
			sentence = append(sentence, clause)
			clause = Clause{}
		}
		clause.Words = append(clause.Words, word)
		if flag, terminator := terminatesClause(bareToken); flag {
			clause.Terminator = terminator
			sentence = append(sentence, clause)
			clause = Clause{}
		} else if flag, terminator := terminatesSentence(bareToken); flag {
			clause.Terminator = terminator
			sentence = append(sentence, clause)
			clause = Clause{}
			res = append(res, sentence)
			sentence = nil
		}
	}
	if clause.Words != nil {
		sentence = append(sentence, clause)
	}
	if sentence != nil {
		res = append(res, sentence)
	}
	return res
}

func normalizeUnicodeSymbols(text string) string {
	replacements := []string{
		"\u201c", "\"",
		"\u201d", "\"",
		"\u2018", "'",
		"\u2019", "'",
		"\u2010", "--",
		"\u2011", "--",
		"\u2012", "--",
		"\u2013", "--",
		"\u2014", "--",
		"\u2015", "--",
	}
	for i := 0; i < len(replacements); i += 2 {
		text = strings.Replace(text, replacements[i], replacements[i+1], -1)
	}
	return text
}

func isTitle(text string) bool {
	for _, title := range Titles {
		if title == text {
			return true
		}
	}
	return false
}

func stripPunctuation(word string) string {
	if isTitle(word) {
		return word
	}

	terminators := []string{}
	terminators = append(terminators, SentenceTerminators...)
	terminators = append(terminators, ClauseTerminators...)
	for _, terminator := range terminators {
		if strings.HasSuffix(word, terminator) {
			word = word[:len(word)-len(terminator)]
		}
	}
	for _, intro := range ClauseIntroducers {
		if strings.HasPrefix(word, intro) {
			word = word[len(intro):]
		}
	}
	return word
}

func wordCapitalization(word string) Capitalization {
	if strings.ToUpper(word) == word {
		return AllCapital
	} else if strings.ToLower(word) == word {
		return NoCapital
	} else {
		return SomeCapital
	}
}

func terminatesClause(word string) (bool, string) {
	for _, terminator := range ClauseTerminators {
		if strings.HasSuffix(word, terminator) {
			return true, terminator
		}
	}
	return false, ""
}

func terminatesSentence(word string) (bool, string) {
	if isTitle(word) {
		return false, ""
	}
	for _, terminator := range SentenceTerminators {
		if strings.HasSuffix(word, terminator) {
			return true, terminator
		}
	}
	return false, ""
}

func introducesClause(word string) bool {
	for _, intro := range ClauseIntroducers {
		if strings.HasPrefix(word, intro) {
			return true
		}
	}
	return false
}
