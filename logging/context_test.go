package logging_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/anz-bank/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddLoggerToContext(t *testing.T) {
	ctx := logging.New(io.Discard).With().ToContext(context.Background())
	require.NotPanics(t, func() {
		logging.Info(ctx).Msg("Hello World")
	})
}

func TestGetLoggerFromContext(t *testing.T) {
	logger := logging.New(io.Discard)
	ctx := logger.ToContext(context.Background())
	assert.NotPanics(t, func() { _ = logging.FromContext(ctx) })
	// assert.Equal(t, logger, loggerFromCtx)
}

func TestLogFuncInfo(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf)
	ctx := logger.ToContext(context.Background())
	logging.Info(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
}

func TestLogFuncInfoUsesGlobalIfNoLogger(t *testing.T) {
	buf := bytes.Buffer{}
	logging.SetGlobalLogContext(logging.New(&buf).With())
	logging.Info(context.Background()).Msg("Hello World")
	assert.Contains(t, buf.String(), `"globalLogger":true`)
}

func TestLogFuncDebug(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
}

func TestLogFuncDebugUsesGlobalIfNoLogger(t *testing.T) {
	buf := bytes.Buffer{}
	logging.SetGlobalLogContext(logging.New(&buf).WithLevel(logging.DebugLevel).With())
	logging.Debug(context.Background()).Msg("Hello World")
	assert.Contains(t, buf.String(), `"globalLogger":true`)
}

func TestLogFuncError(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf)
	ctx := logger.ToContext(context.Background())
	logging.Error(ctx, fmt.Errorf("Bad things")).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "Bad things")
}

func TestLogFuncErrorUsesGlobalIfNoLogger(t *testing.T) {
	buf := bytes.Buffer{}
	logging.SetGlobalLogContext(logging.New(&buf).With())
	logging.Error(context.Background(), fmt.Errorf("Bad things")).Msg("Hello World")
	assert.Contains(t, buf.String(), `"globalLogger":true`)
}

func TestLogContextFromContextWithDur(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	loggerContext := logging.ContextFromContext(ctx)
	loggerContext = loggerContext.WithDur("sale_duration", 72*time.Hour)
	ctx = loggerContext.ToContext(ctx)
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "sale_duration")
	assert.Contains(t, buf.String(), "259200000")
}

func TestLogContextFromContextWithTime(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	loggerContext := logging.ContextFromContext(ctx)
	loggerContext = loggerContext.WithTime("sale_start", time.Date(2020, time.November, 23, 0, 0, 0, 0, time.UTC))
	ctx = loggerContext.ToContext(ctx)
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "sale_start")
	assert.Contains(t, buf.String(), "2020-11-23T00:00:00Z")
}

func TestLogContextFromContextWithStr(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	loggerContext := logging.ContextFromContext(ctx)
	loggerContext = loggerContext.WithStr("fruit", "banana")
	ctx = loggerContext.ToContext(ctx)
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "fruit")
	assert.Contains(t, buf.String(), "banana")
}

func TestLogContextFromContextWithInt(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	loggerContext := logging.ContextFromContext(ctx)
	loggerContext = loggerContext.WithInt("price", 100)
	ctx = loggerContext.ToContext(ctx)
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "price")
	assert.Contains(t, buf.String(), "100")
}

func TestLogContextFromContextWithDict(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	loggerContext := logging.ContextFromContext(ctx)
	loggerContext = loggerContext.WithDict("specials", logging.Dict().
		Str("bundle", "bunch").
		Int("price", 299))
	ctx = loggerContext.ToContext(ctx)
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "specials")
	assert.Contains(t, buf.String(), "bundle")
	assert.Contains(t, buf.String(), "bunch")
	assert.Contains(t, buf.String(), "price")
	assert.Contains(t, buf.String(), "299")
}

func TestLogContextFromContextWithArray(t *testing.T) {
	buf := bytes.Buffer{}
	logger := logging.New(&buf).WithLevel(logging.DebugLevel)
	ctx := logger.ToContext(context.Background())
	loggerContext := logging.ContextFromContext(ctx)
	loggerContext = loggerContext.WithStr("fruit", "banana")
	loggerContext = loggerContext.WithInt("price", 100)
	loggerContext = loggerContext.WithDict("specials", logging.Dict().
		Str("bundle", "bunch").
		Int("price", 299))
	loggerContext = loggerContext.WithArray("sales", logging.Array().Int(3).Int(2).Int(1).Int(9).Int(0))
	ctx = loggerContext.ToContext(ctx)
	logging.Debug(ctx).Msg("Hello World")
	assert.Contains(t, buf.String(), "Hello World")
	assert.Contains(t, buf.String(), "sales")
	assert.Contains(t, buf.String(), "3")
	assert.Contains(t, buf.String(), "2")
	assert.Contains(t, buf.String(), "1")
	assert.Contains(t, buf.String(), "9")
	assert.Contains(t, buf.String(), "0")
}
