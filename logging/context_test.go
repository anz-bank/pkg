package logging_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/anz-bank/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddLoggerToContext(t *testing.T) {
	ctx := logging.New(ioutil.Discard).With().ToContext(context.Background())
	require.NotPanics(t, func() {
		logging.Info(ctx).Msg("Hello World")
	})
}

func TestGetLoggerFromContext(t *testing.T) {
	logger := logging.New(ioutil.Discard)
	ctx := logger.ToContext(context.Background())
	assert.NotPanics(t, func() { _ = logging.FromContext(ctx) })
	// assert.Equal(t, logger, loggerFromCtx)
}

func TestLogFuncs(t *testing.T) {
	t.Run("Info", func(t *testing.T) {
		buf := bytes.Buffer{}
		logger := logging.New(&buf)
		ctx := logger.ToContext(context.Background())
		logging.Info(ctx).Msg("Hello World")
		assert.Contains(t, buf.String(), "Hello World")
	})
	t.Run("InfoUsesGlobalIfNoLogger", func(t *testing.T) {
		buf := bytes.Buffer{}
		logging.SetGlobalLogContext(logging.New(&buf).With())
		logging.Info(context.Background()).Msg("Hello World")
		assert.Contains(t, buf.String(), `"globalLogger":true`)
	})

	t.Run("Debug", func(t *testing.T) {
		buf := bytes.Buffer{}
		logger := logging.New(&buf).WithLevel(logging.DebugLevel)
		ctx := logger.ToContext(context.Background())
		logging.Debug(ctx).Msg("Hello World")
		assert.Contains(t, buf.String(), "Hello World")
	})
	t.Run("DebugUsesGlobalIfNoLogger", func(t *testing.T) {
		buf := bytes.Buffer{}
		logging.SetGlobalLogContext(logging.New(&buf).WithLevel(logging.DebugLevel).With())
		logging.Debug(context.Background()).Msg("Hello World")
		assert.Contains(t, buf.String(), `"globalLogger":true`)
	})

	t.Run("Error", func(t *testing.T) {
		buf := bytes.Buffer{}
		logger := logging.New(&buf)
		ctx := logger.ToContext(context.Background())
		logging.Error(ctx, fmt.Errorf("Bad things")).Msg("Hello World")
		assert.Contains(t, buf.String(), "Hello World")
		assert.Contains(t, buf.String(), "Bad things")
	})
	t.Run("ErrorUsesGlobalIfNoLogger", func(t *testing.T) {
		buf := bytes.Buffer{}
		logging.SetGlobalLogContext(logging.New(&buf).With())
		logging.Error(context.Background(), fmt.Errorf("Bad things")).Msg("Hello World")
		assert.Contains(t, buf.String(), `"globalLogger":true`)
	})
}
