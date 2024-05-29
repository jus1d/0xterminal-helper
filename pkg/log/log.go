package log

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO(#12): migrate to log/slog with pretty logging for local environment

type Record struct {
	Timestamp time.Time
	Content   string
	Level     string
	Error     error
	Attrs     []Attr
}

type Attr struct {
	Key   string
	Value string
}

func Info(content string, attrs ...Attr) {
	timestamp := time.Now()
	record := Record{
		Timestamp: timestamp,
		Content:   content,
		Level:     "INFO",
		Attrs:     attrs,
	}

	record.Write(os.Stdout)
}

func Warn(content string, attrs ...Attr) {
	record := Record{
		Timestamp: time.Now(),
		Content:   content,
		Level:     "WARN",
		Attrs:     attrs,
	}
	record.Write(os.Stdout)
}

func Error(content string, err error, attrs ...Attr) {
	record := Record{
		Timestamp: time.Now(),
		Content:   content,
		Level:     "ERROR",
		Error:     err,
		Attrs:     attrs,
	}
	record.Write(os.Stdout)
}

func Fatal(content string, err error, attrs ...Attr) {
	record := Record{
		Timestamp: time.Now(),
		Content:   content,
		Level:     "FATAL",
		Error:     err,
		Attrs:     attrs,
	}
	record.Write(os.Stdout)
	os.Exit(1)
}

func (r *Record) Write(out *os.File) {
	json := fmt.Sprintf(`{"time":"%s","level":"%s","message":"%s"`,
		r.Timestamp.Format("2006-01-02T15:04:05.999999999Z"), r.Level, r.Content)

	if r.Error != nil {
		json += fmt.Sprintf(`,"error":"%s"`, r.Error.Error())
	}
	for _, attr := range r.Attrs {
		json += fmt.Sprintf(`,"%s":%s`, attr.Key, strings.Replace(attr.Value, "\n", "\\n", -1))
	}
	json += `}`
	fmt.Fprintln(out, json)
}

func WithString(key string, value string) Attr {
	return Attr{
		Key:   key,
		Value: fmt.Sprintf(`"%s"`, value),
	}
}

func WithInt(key string, value int) Attr {
	return Attr{
		Key:   key,
		Value: strconv.Itoa(value),
	}
}

func WithInt64(key string, value int64) Attr {
	return Attr{
		Key:   key,
		Value: strconv.FormatInt(value, 10),
	}
}
