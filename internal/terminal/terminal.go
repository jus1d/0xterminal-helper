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
	ErrInsufficientWords    = errors.New("terminal.Game.New(): insufficient words list")
)

type Game struct {
	initialWords   []string
	availableWords []string
	attempts       []*attempt
}

type attempt struct {
	Word           string
	GuessedLetters int
}

func New(words []string) (*Game, error) {
	if !isWordsEqualLength(words) {
		return nil, ErrDifferentWordsLength
	}

	if len(words) < 6 {
		return nil, ErrInsufficientWords
	}

	for i := range words {
		words[i] = strings.TrimSpace(strings.ToLower(words[i]))
	}

	game := &Game{
		initialWords:   words,
		availableWords: words,
		attempts:       make([]*attempt, 0),
	}
	game.sortWordsBySexyIndex()

	return game, nil
}

func Attempt(word string, guessedLetters int) attempt {
	return attempt{
		Word:           word,
		GuessedLetters: guessedLetters,
	}
}

func (g *Game) Words() []string {
	return g.initialWords
}

func (g *Game) AvailableWords() []string {
	return g.availableWords
}

func (g *Game) Target() string {
	if len(g.availableWords) != 1 {
		return ""
	}
	return g.availableWords[0]
}

func (g *Game) Attempts() int {
	n := 0
	for _, attempt := range g.attempts {
		if attempt.GuessedLetters != len(g.initialWords[0]) {
			n++
		}
	}
	return n + 1
}

func (g *Game) SubmitAttempt(attempt attempt) {
	g.attempts = append(g.attempts, &attempt)
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
	for _, word := range g.availableWords {
		fits := true
		for _, tried := range g.attempts {
			if !compareWordsMatchedLetters(word, tried.Word, tried.GuessedLetters) {
				fits = false
				break
			}
		}
		if fits {
			updated = append(updated, word)
		}
	}
	g.availableWords = updated
}

func (g *Game) sortWordsBySexyIndex() {
	sort.Slice(g.availableWords, func(i, j int) bool {
		return countWordSexyIndex(g.availableWords[i], g.availableWords) < countWordSexyIndex(g.availableWords[j], g.availableWords)
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
