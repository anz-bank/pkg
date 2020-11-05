package logging

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// Level types the level construct
//
// Mirrors zerolog.Level
type Level uint32

// zerolog level constants included here. Use these to avoid importing logrus directly
// Not all levels are included, only those used by this library
const (
	DebugLevel Level = Level(zerolog.DebugLevel)
	InfoLevel  Level = Level(zerolog.InfoLevel)
	ErrorLevel Level = Level(zerolog.ErrorLevel)
)

func (l Level) String() string {
	return zerolog.Level(l).String()
}

// MustParseLevel panics if input is not a valid level string
func MustParseLevel(s string) Level {
	l, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return l
}

// ParseLevel parses a level string and returns the equivalent level value
//
// Valid levels are info, debug, error. NOT CASE SENSITIVE
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "error":
		return ErrorLevel, nil
	default:
		return 0, fmt.Errorf("unrecognised log level %s", s)
	}
}
