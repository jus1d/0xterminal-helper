package terminal

import (
	"errors"
	"fmt"
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

	return &Game{
		Words:    words,
		Attempts: make([]*Attempt, 0),
	}, nil
}

func (g *Game) CommitAttempt(attempt Attempt) {
	g.Attempts = append(g.Attempts, &attempt)
}

func (g *Game) UpdateWords() {
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

func compareWords(a string, b string, expected int) bool {
	actual := 0
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			actual++
		}
	}
	return expected == actual
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
