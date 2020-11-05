package logging_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/anz-bank/pkg/logging"
	"github.com/anz-bank/pkg/logging/codelinks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	buf := bytes.Buffer{}
	logger := logging.New(&buf)
	logger.Info().Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), `"level":"info"`)
}

func TestLogger_With(t *testing.T) {
	t.Parallel()
	buf := bytes.Buffer{}
	logger := logging.New(&buf)
	logger2 := logger.WithStr("key", "val")

	// Run first logger then second to check immutability of the first
	logger.Info().Msg("Hello World")
	assert.NotContains(t, buf.String(), `"key":"val"`)

	logger2.Info().Msg("Hello World")
	assert.Contains(t, buf.String(), `"key":"val"`)
}

func TestLogger_WithOutput(t *testing.T) {
	t.Parallel()
	buf := bytes.Buffer{}
	logger := logging.New(&buf)

	buf2 := bytes.Buffer{}
	logger2 := logger.WithOutput(&buf2)

	// Run first logger then second to check immutability of the first
	logger.Info().Msg("Hello World")
	assert.NotEmpty(t, buf.String())
	assert.Empty(t, buf2.String())

	logger2.Info().Msg("Hello World")
	assert.NotEmpty(t, buf2.String())
}

func TestLogger_WithLevel(t *testing.T) {
	t.Parallel()
	buf := bytes.Buffer{}
	logger := logging.New(&buf)
	logger2 := logger.WithLevel(logging.DebugLevel)

	// Run first logger then second to check immutability of the first
	logger.Debug().Msg("Hello World")
	assert.Empty(t, buf.String())

	logger2.Debug().Msg("Hello World")
	assert.NotEmpty(t, buf.String())
}

func TestLogger_WithCodeLinks(t *testing.T) {
	t.Parallel()
	buf := bytes.Buffer{}
	logger := logging.New(&buf)
	logger2 := logger.WithCodeLinks(true, codelinks.LocalLinker{})

	// Run first logger then second to check immutability of the first
	logger.Info().Msg("Hello World")
	assert.NotContains(t, buf.String(), "source_code")

	logger2.Info().Msg("Hello World")
	// We can't guarantee the link follows a particular path across multiple machines
	assert.Contains(t, buf.String(), `"source_code":`)
	assert.Contains(t, buf.String(), `logger_test.go:77`)
}

func TestLogger_LogFuncs(t *testing.T) {
	t.Parallel()

	t.Run("Info", func(t *testing.T) {
		t.Parallel()
		buf := bytes.Buffer{}
		logger := logging.New(&buf)
		logger.Info().Msg("Hello World")
		assert.Contains(t, buf.String(), "Hello World")
		assert.Contains(t, buf.String(), `"level":"info"`)
	})

	t.Run("Debug", func(t *testing.T) {
		t.Parallel()
		buf := bytes.Buffer{}
		logger := logging.New(&buf).WithLevel(logging.DebugLevel)
		logger.Debug().Msg("Hello World")
		assert.Contains(t, buf.String(), "Hello World")
		assert.Contains(t, buf.String(), `"level":"debug"`)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		buf := bytes.Buffer{}
		logger := logging.New(&buf)
		logger.Error(fmt.Errorf("Bad things")).Msg("Uh oh")
		assert.Contains(t, buf.String(), `"level":"error"`)
		assert.Contains(t, buf.String(), "Uh oh")
		assert.Contains(t, buf.String(), `"error":"Bad things"`)
	})
}
