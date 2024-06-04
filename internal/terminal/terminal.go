package terminal

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"terminal/pkg/slice"
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
	if !isWordsEqualLength(words) {
		return nil, ErrDifferentWordsLength
	}

	for i := range words {
		words[i] = strings.TrimSpace(strings.ToLower(words[i]))
	}

	game := &Game{
		InitialWords:   words,
		AvailableWords: words,
		Attempts:       make([]*Attempt, 0),
	}
	game.sortWordsBySexyIndex()

	return game, nil
}

func (g *Game) SubmitAttempt(attempt Attempt) {
	g.Attempts = append(g.Attempts, &attempt)
	g.updateWords()
	g.sortWordsBySexyIndex()
}

func RemoveTrashFromWordList(words []string) []string {
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
	return slice.Unique(cleaned)
}

func ComputeWordsHash(words []string) string {
	sort.Strings(words)

	var builder strings.Builder
	for _, word := range words {
		builder.WriteString(word)
		builder.WriteString(" ")
	}

	checksum := sha256.Sum256([]byte(builder.String()))

	return hex.EncodeToString(checksum[:])
}

func (g *Game) updateWords() {
	updated := make([]string, 0)
	for _, word := range g.AvailableWords {
		fits := true
		for _, tried := range g.Attempts {
			if !compareWordsMatchedLetters(word, tried.Word, tried.GuessedLetters) {
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

func (g *Game) sortWordsBySexyIndex() {
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

func compareWordsMatchedLetters(a string, b string, expected int) bool {
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

func isWordsEqualLength(words []string) bool {
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
