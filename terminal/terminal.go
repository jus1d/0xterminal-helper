package terminal

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrDifferentWordsLength = errors.New("terminal.Game.New(): words could not be different length")
)

type Game struct {
	Words    []string
	Attempts []*Attempt
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
		Words:    words,
		Attempts: make([]*Attempt, 0),
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
	for _, word := range g.Words {
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
	g.Words = updated
}

func (g *Game) IsFinished() bool {
	return len(g.Words) == 1
}

func (g *Game) PrintAvailableWords() {
	fmt.Printf("Available %d words:\n", len(g.Words))
	for i, word := range g.Words {
		fmt.Printf("#%d: %s\n", i, word)
	}
}

func (g *Game) sortWords() {
	sort.Slice(g.Words, func(i, j int) bool {
		return countWordSexyIndex(g.Words[i], g.Words) < countWordSexyIndex(g.Words[j], g.Words)
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
