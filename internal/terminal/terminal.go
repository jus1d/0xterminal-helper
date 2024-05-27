package terminal

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrDifferentWordsLength = errors.New("terminal.Game.New(): words could not be different length")
)

type Game struct {
	InitialWords   []string
	AvailableWords []string
	Attempts       []*Attempt
}

type Attempt struct {
	Word           string
	GuessedLetters int
}

func New(words []string) (*Game, error) {
	if !checkWordsLength(words) {
		return nil, ErrDifferentWordsLength
	}

	game := &Game{
		InitialWords:   words,
		AvailableWords: words,
		Attempts:       make([]*Attempt, 0),
	}
	game.sortWords()

	return game, nil
}

func (g *Game) CommitAttempt(attempt Attempt) {
	g.Attempts = append(g.Attempts, &attempt)
	g.updateWords()
	g.sortWords()
}

func (g *Game) updateWords() {
	updated := make([]string, 0)
	for _, word := range g.AvailableWords {
		fits := true
		for _, tried := range g.Attempts {
			if !compareWords(word, tried.Word, tried.GuessedLetters) {
				fits = false
				break
			}
		}
		if fits {
			updated = append(updated, word)
		}
	}
	g.AvailableWords = updated
}

func (g *Game) IsFinished() bool {
	return len(g.AvailableWords) == 1
}

func (g *Game) PrintAvailableWords() {
	fmt.Printf("Available %d words:\n", len(g.AvailableWords))
	for i, word := range g.AvailableWords {
		fmt.Printf("#%d: %s\n", i, word)
	}
}

func RemoveTrashFromWordsList(words []string) []string {
	cleaned := make([]string, 0)
	for _, word := range words {
		if strings.Contains(word, "NaN") {
			continue
		}
		if word[0] == '(' || word[0] == '{' || word[0] == '[' {
			continue
		}
		cleaned = append(cleaned, word)
	}
	return unique(cleaned)
}

func unique(words []string) []string {
	wordMap := make(map[string]bool)
	var uniqueWords []string

	for _, word := range words {
		if !wordMap[word] {
			wordMap[word] = true
			uniqueWords = append(uniqueWords, word)
		}
	}

	return uniqueWords
}

func (g *Game) sortWords() {
	sort.Slice(g.AvailableWords, func(i, j int) bool {
		return countWordSexyIndex(g.AvailableWords[i], g.AvailableWords) < countWordSexyIndex(g.AvailableWords[j], g.AvailableWords)
	})
}

func countWordSexyIndex(target string, words []string) int {
	matches := make(map[int]int, 0) // key: matches in word, value: words with this amount of matches
	for i := 0; i < len(words); i++ {
		if words[i] == target {
			continue
		}

		matches[countMatchedLetters(target, words[i])]++
	}

	sum := 0
	max := 0
	for _, v := range matches {
		sum += v
		if v > max {
			max = v
		}
	}
	if len(matches) == 0 {
		return 0
	}
	average := sum / len(matches)
	sexyIndex := max - average

	return sexyIndex
}

func compareWords(a string, b string, expected int) bool {
	return expected == countMatchedLetters(a, b)
}

func countMatchedLetters(a, b string) int {
	value := 0
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			value++
		}
	}
	return value
}

func checkWordsLength(words []string) bool {
	if len(words) == 0 {
		return true
	}

	expectedLength := len(words[0])
	for _, word := range words[1:] {
		if len(word) != expectedLength {
			return false
		}
	}

	return true
}
