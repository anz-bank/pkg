package log

import (
	"bytes"
	"context"
	"strconv"
	"testing"

	"github.com/anz-bank/pkg/log"
	"github.com/sirupsen/logrus"
)

func BenchmarkLog5Fields(b *testing.B) {
	runBenchmark(b, 5)
}

func BenchmarkLog10Fields(b *testing.B) {
	runBenchmark(b, 10)
}

func BenchmarkLog50Fields(b *testing.B) {
	runBenchmark(b, 50)
}

func BenchmarkLog100Fields(b *testing.B) {
	runBenchmark(b, 100)
}

func BenchmarkLog1000Fields(b *testing.B) {
	runBenchmark(b, 1000)
}

func BenchmarkWith(b *testing.B) {
	logger := log.NewNullLogger()
	ctx := log.With("key", "val").With("abc", 123).WithLogger(logger).Onto(context.Background())

	for i := 0; i < b.N; i++ {
		log.With("foo", "bar").
			With("hello", "world").
			With("myvar", "myval").
			With("this", "that").
			Onto(ctx)
	}
}

func BenchmarkLogrus(b *testing.B) {
	logger := logrus.New()
	logger.SetOutput(&bytes.Buffer{})
	logger.SetReportCaller(true)

	for i := 0; i < b.N; i++ {
		logger.Info("TestMsg")
	}
}

func BenchmarkLog(b *testing.B) {
	logger := log.NewNullLogger()
	ctx := log.With("x-user-id", "12344").
		With("x-trace-id", "acbdd").
		WithLogger(logger).
		Onto(context.Background())

	for i := 0; i < b.N; i++ {
		log.From(ctx).Info("TestMsg")
	}
}

func runBenchmark(b *testing.B, l int) {
	for i := 0; i < b.N; i++ {
		var f log.Fields
		for j := 0; j < l; j++ {
			f = f.With(strconv.Itoa(j), j)
		}
		f.WithLogger(log.NewNullLogger()).From(context.Background()).Info("test")
	}
}
