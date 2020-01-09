package log

import (
	"context"
	"testing"

	"github.com/anz-bank/pkg/log/loggers"
	"github.com/anz-bank/pkg/log/testutil"
	"github.com/marcelocantos/frozen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWithLogger(t *testing.T) {
	t.Parallel()

	logger := loggers.NewMockLogger()
	setMockCopyAssertion(logger)
	With(context.Background(), logger)

	require.True(t, logger.AssertExpectations(t))
}

func TestWithFields(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		c := c
		t.Run("TestWithFields"+" "+c.Name, func(tt *testing.T) {
			tt.Parallel()

			logger := loggers.NewMockLogger()
			setLogSpecificFieldAssertion(logger, c.Fields)

			ctx := context.WithValue(context.Background(), loggerKey, logger)
			WithFields(ctx, testutil.ConvertToGoMap(c.Fields))
			logger = ctx.Value(loggerKey).(*loggers.MockLogger)

			assert.True(tt, logger.AssertExpectations(tt))
		})
	}
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name   string
		fields []Fields
	}{
		{
			name:   "Empty fields",
			fields: []Fields{},
		},
		{
			name: "Single fields only",
			fields: []Fields{
				NewField("test", 1),
				NewField("test again", 2),
				NewField("another test", 3),
			},
		},
		{
			name: "Multiple fields only",
			fields: []Fields{
				FromMap(
					map[string]interface{}{
						"test":  '1',
						"test2": 2,
					},
				),
				FromMap(
					map[string]interface{}{
						"test3": 3,
						"test4": 4,
					},
				),
			},
		},
		{
			name: "Mixed Fields",
			fields: []Fields{
				NewField("test", 1),
				FromMap(
					map[string]interface{}{
						"test3": 3,
						"test4": 4,
					},
				),
				NewField("another test", 2),
			},
		},
	}

	for _, ts := range tests {
		ts := ts
		t.Run(ts.name, func(tt *testing.T) {
			tt.Parallel()

			logger := loggers.NewMockLogger()
			setLogSpecificFieldAssertion(logger, fromFields(ts.fields))

			ctx := context.WithValue(context.Background(), loggerKey, logger)
			From(ctx, ts.fields...)
			assert.True(tt, logger.AssertExpectations(tt))
		})
	}
}

func setLogSpecificFieldAssertion(logger *loggers.MockLogger, fields frozen.Map) {
	logger.On("Copy").Return(logger)
	if fields.Count() != 0 {
		logger.On(
			"PutFields",
			mock.MatchedBy(
				func(arg frozen.Map) bool {
					return fields.Equal(arg)
				},
			),
		).Return(logger)
	}
}

func setMockCopyAssertion(logger *loggers.MockLogger) {
	// set to return the same logger for testing purposes, in real case it will return
	// a copied logger. Tests that use these usually are not checked for their return value
	// as the return value is mocked
	logger.On("Copy").Return(logger)
}
