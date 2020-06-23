package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/anz-bank/pkg/log"
	"github.com/stretchr/testify/require"
)


func TestLogCallerFromLog(t *testing.T) {
	ctx := getContextWithLogger()
	buffer := bytes.Buffer{}
	ctx = log.WithConfigs(log.SetOutput(&buffer)).Onto(ctx)

	log.Info(ctx, "Info")
	require.False(t, strings.Contains(buffer.String(), "logentryCaller_test.go")) //don't log caller

	//set log caller
	buffer.Reset()
	ctx = log.WithConfigs(log.SetLogCaller(true)).Onto(ctx)
	log.Info(ctx, "Info")
	require.True(t, strings.Contains(buffer.String(), "logentryCaller_test.go")) //log caller
}

func TestLogCallerFromField(t *testing.T) {
	ctx := getContextWithLogger()
	buffer := bytes.Buffer{}
	ctx = log.WithConfigs(log.SetOutput(&buffer)).Onto(ctx)

	fields := log.With("hello", "world")
	ctx = fields.Onto(ctx)
	fields.Info(ctx, "Info")
	require.False(t, strings.Contains(buffer.String(), "logentryCaller_test.go")) //don't log caller

	//set log caller
	buffer.Reset()
	ctx = log.WithConfigs(log.SetLogCaller(true)).Onto(ctx)
	fields.Info(ctx, "Info")
	require.True(t, strings.Contains(buffer.String(), "logentryCaller_test.go")) //log caller
}

func getContextWithLogger() context.Context {
	ctx := context.Background()
	logger := log.NewStandardLogger()
	ctx = log.WithLogger(logger).Onto(ctx)
	return ctx
}
