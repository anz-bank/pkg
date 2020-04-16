package log

import (
	"github.com/alecthomas/assert"
	"github.com/arr-ai/frozen"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLogrusPkgFormatters(t *testing.T) {
	t.Parallel()

	entry := &LogEntry{
		Time:    time.Now(),
		Message: "message",
		Data:    frozen.NewMap(frozen.KV("cat", "dog")),
		Verbose: true,
	}

	logger := logrus.New()
	formatter := standardFormat{}
	formatted, err := formatter.Format(entry)
	require.NoError(t, err)
	doubleWrappedFormatter := logrusFormatterToPkgFormatter{logger, pkgFormatterToLogrusFormatter{formatter}}
	doubleWrappedFormatted, err := doubleWrappedFormatter.Format(entry)
	require.NoError(t, err)
	assert.Equal(t, formatted, doubleWrappedFormatted)
}

func TestLogrusPkgEntry(t *testing.T) {
	t.Parallel()

	entry := &LogEntry{
		Time:    time.Now(),
		Message: "message",
		Data:    frozen.NewMap(frozen.KV("cat", "dog")),
		Caller:  CodeReference{"example.go", 123},
		Verbose: true,
	}

	logger := logrus.New()
	roundTripEntry := logrusEntryToPkgLogEntry(pkgLogEntryToLogrusEntry(logger, entry))
	assert.Equal(t, entry, roundTripEntry)
}
