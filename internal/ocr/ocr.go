package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"terminal/pkg/slice"
	"time"
)

type Client struct {
	tokens []string
}

type responseOCR struct {
	Results []parsedResultOCR `json:"ParsedResults"`
}

type parsedResultOCR struct {
	Text string `json:"ParsedText"`
	Err  string `json:"ErrorMessage"`
}

type response struct {
	text string
	err  error
}

func New(tokens []string) *Client {
	return &Client{tokens}
}

func (c *Client) ExtractWords(ctx context.Context, filepath string) ([]string, error) {
	limitation := 3 * time.Second
	ctx, cancel := context.WithTimeout(ctx, limitation)
	defer cancel()

	respch := make(chan response)

	go func() {
		text, err := c.extractTextFromImage(filepath)
		respch <- response{text, err}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("ocr: extracting text from image took too long")
		case resp := <-respch:
			words := findWords(resp.text)
			return words, resp.err
		}
	}
}

func (c *Client) extractTextFromImage(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	writer.WriteField("language", "eng")
	writer.WriteField("isOverlayRequired", "true")

	err = writer.Close()
	if err != nil {
		return "", err
	}

	endpoint := "https://api.ocr.space/Parse/Image"
	req, err := http.NewRequest("POST", endpoint, &requestBody)
	if err != nil {
		return "", err
	}
	token := slice.Choose(c.tokens)
	req.Header.Add("apikey", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response responseOCR
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Results) == 0 {
		return "", errors.New("ocr: no parsed results")
	}
	if response.Results[0].Err != "" {
		return response.Results[0].Text, fmt.Errorf("ocr: %s", response.Results[0].Err)
	}
	return response.Results[0].Text, nil
}

func findWords(text string) []string {
	words := make([]string, 0)

	lines := strings.Split(text, "\r\n")

	for _, line := range lines {
		word := strings.TrimSpace(line)

		if !isWord(word) {
			continue
		}

		words = append(words, word)
	}

	correctLength := findCorrectLength(words)
	words = filterWordsByLength(words, correctLength)

	return words
}

func findCorrectLength(words []string) int {
	lengthCount := make(map[int]int)
	for _, word := range words {
		length := len(word)
		lengthCount[length]++
	}

	maxCount := 0
	correctLength := 0
	for length, count := range lengthCount {
		if count > maxCount {
			maxCount = count
			correctLength = length
		}
	}

	return correctLength
}

func filterWordsByLength(words []string, correctLength int) []string {
	var filteredWords []string
	for _, word := range words {
		if len(word) == correctLength {
			filteredWords = append(filteredWords, word)
		}
	}
	return filteredWords
}

func isWord(word string) bool {
	// minimal word's length in TERMINAL is 4
	if len(word) < 4 {
		return false
	}

	// in TERMINAL we need only words, that contains only lowercase english letters
	for _, char := range word {
		if char < 'a' || char > 'z' {
			return false
		}
	}
	return true
}
