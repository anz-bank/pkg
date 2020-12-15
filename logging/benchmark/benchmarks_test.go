package logging_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/anz-bank/pkg/logging"
	"github.com/anz-bank/pkg/logging/codelinks"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

func BenchmarkLogTime(b *testing.B) {
	ctx := logging.New(&bytes.Buffer{}).
		WithCodeLinks(true, codelinks.LocalLinker{}).
		WithStr("foo", "bar").
		ToContext(context.Background())
	for n := 0; n < b.N; n++ {
		logging.Info(ctx).Msg("TestMsg")
	}
	// BenchmarkLogTime-12    	  695059	      1595 ns/op	    1344 B/op	       5 allocs/op
}

func BenchmarkLogrusLogTime(b *testing.B) {
	logrus.SetOutput(&bytes.Buffer{})
	logrus.SetReportCaller(true)
	for n := 0; n < b.N; n++ {
		logrus.WithField("foo", "bar").Info("TestMsg")
	}
	// BenchmarkLogrusLogTime-12    	  210763	      6703 ns/op	    2234 B/op	      27 allocs/op
}

func BenchmarkLogTime_StdFields(b *testing.B) {
	ocFunc := logging.ContextFunc{
		Keys: []string{"x-b3-traceid"},
		Function: func(ctx context.Context, logCtx zerolog.Context) zerolog.Context {
			span := trace.FromContext(ctx)
			spanCtx := span.SpanContext()
			logCtx = logCtx.
				Str("x-b3-traceid", spanCtx.TraceID.String()).
				Str("x-b3-spanid", spanCtx.SpanID.String())
			return logCtx
		},
	}
	ctx := logging.New(&bytes.Buffer{}).
		WithCodeLinks(true, codelinks.LocalLinker{}).
		WithTimeDiff("time_since_benchmark_start", time.Now()).
		WithStr("service", "test-service").
		WithStr("version", "v1.0.0").
		WithStr("rpc", "/my.rpc").
		With(
			ocFunc,
		).ToContext(context.Background())
	ctx, _ = trace.StartSpan(ctx, "my-span")
	for n := 0; n < b.N; n++ {
		logging.Info(ctx).Msg("TestMsg")
	}
	// BenchmarkLogTime_StdFields-12    	  413539	      2474 ns/op	    1846 B/op	      11 allocs/op
}
