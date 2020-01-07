package log

import (
	"context"
	"testing"

	"github.com/anz-bank/pkg/log/loggers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCopiedLogger(t *testing.T) {
	t.Parallel()

	t.Run("Context has no logger", func(tt *testing.T) {
		tt.Parallel()

		require.Panics(tt, func() {
			getCopiedLogger(context.Background())
		})
	})
	t.Run("Context has a logger", func(tt *testing.T) {
		tt.Parallel()

		ctx := context.Background()
		ctx = WithLogger(ctx, loggers.NewStandardLogger())

		logger := getCopiedLogger(ctx)
		require.NotNil(t, logger)

		fromContext := ctx.Value(loggerKey).(loggers.Logger)
		assert.True(t, logger != fromContext)
	})
}
