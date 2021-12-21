// Package otelhealth exports health.State data via OpenTelemetry.
//
// otelhealth allows you to take the state from a health server and export the
// data in it via OpenTelemetry metrics. For example, to track the readiness of
// your server over time you can use OpenTelemetry to export via Prometheus metrics
// where you can graph and alert on your service not being ready. It also
// allows you to export version information so you can easily correlate
// version changes with changes in other metrics making it easier to identify
// regressions.
//
// The servers in github.com/anz-bank/pkg/health all have a State field. This
// can be passed to the Register function in this package to have that state
// provided as OpenTelemetry metrics.
//
// Two metrics are published by this package:
//   anz_health_ready
//   anz_health_version
//
// The "ready" metric tracks the real-time value of the Ready field in the
// State, exporting false as 0 and true as 1.
//
// The "version" metric is constant - it is exported as 1 and never changes.
// It is the labels on this that are the part of interest. The fields of the
// Version in the State are exported as labels on the metric. The current
// labels are "build_log_url", "commit_hash", "container_tag", "semver" and
// "repo_url".
//
// The prefix "anz_health" is configurable with the WithPrefix Option that can
// be passed to Register. The prefix can be removed entirely by using the empty
// string as a prefix.
//
package otelhealth

import (
	"context"
	"sync"

	otelAttribute "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"

	"github.com/anz-bank/pkg/health"
)

const (
	// CommitHash is the git hash.
	CommitHash = otelAttribute.Key("commit_hash")
	// BuildLogURL is the url for the build log.
	BuildLogURL = otelAttribute.Key("build_log_url")
	// ContainerTag is the tag of container.
	ContainerTag = otelAttribute.Key("container_tag")
	// RepoURL is the url of the github repository.
	RepoURL = otelAttribute.Key("repo_url")
	// Semver is the version.
	Semver = otelAttribute.Key("semver")
)

type registerOptions struct {
	metricPrefix string
	constLabels  map[otelAttribute.Key]otelAttribute.Value
}

type metricHandler struct {
	meter               metric.Meter
	metricCounters      map[string]Int64Counter
	metricValueObserver map[string]metric.Int64ValueObserver
}

var mHandler *metricHandler
var mHandlerMutex sync.Mutex

func newMetricHandler() {
	mHandlerMutex.Lock()
	defer mHandlerMutex.Unlock()
	mHandler = &metricHandler{
		meter:               global.Meter(""),
		metricCounters:      make(map[string]Int64Counter),
		metricValueObserver: make(map[string]metric.Int64ValueObserver),
	}
}

// Option is used to configure the registration of the health state. The
// With* functions should be used to obtain options for passing to Register.
type Option func(*registerOptions)

// WithPrefix returns an option that sets the prefix on the metric name
// published by this package. The default is "anz_health" which can be
// overidden with this option. The empty string removes the prefix. Any
// other value will be used as the prefix with an underscore separating
// the prefix from the base metric name.
func WithPrefix(metricPrefix string) Option {
	return func(ro *registerOptions) {
		ro.metricPrefix = metricPrefix
	}
}

// WithConstLabels returns an Option that adds additional labels with constant
// values to the metrics published by this package.
func WithConstLabels(labels map[otelAttribute.Key]otelAttribute.Value) Option {
	return func(ro *registerOptions) {
		for k, v := range labels {
			ro.constLabels[k] = v
		}
	}
}

// Register creates a metrics registry with values from the given health.State.
// The metrics registry tracks changes to the non-constant values
// in the health.State.
func Register(s *health.State, options ...Option) error {
	ro := newRegisterOptions(options...)

	if mHandler == nil {
		newMetricHandler()
	}

	ctx := context.Background()

	if err := addMetrics(ctx, ro, s); err != nil {
		return err
	}

	return nil
}

func newRegisterOptions(options ...Option) *registerOptions {
	ro := &registerOptions{
		metricPrefix: "anz_health",
		constLabels:  map[otelAttribute.Key]otelAttribute.Value{},
	}
	for _, option := range options {
		option(ro)
	}
	if ro.metricPrefix != "" {
		ro.metricPrefix += "_"
	}
	return ro
}

func addMetrics(ctx context.Context, ro *registerOptions, s *health.State) error {
	if err := addReadyMetric(ctx, ro, s); err != nil {
		return err
	}
	if err := addVersionMetric(ctx, ro, s); err != nil {
		return err
	}
	return nil
}

func addReadyMetric(ctx context.Context, ro *registerOptions, s *health.State) error {
	var err error
	mHandlerMutex.Lock()
	defer mHandlerMutex.Unlock()
	readyPrefix := ro.metricPrefix + "ready"
	int64Observer := mHandler.metricValueObserver[readyPrefix]
	if int64Observer == (metric.Int64ValueObserver{}) {
		int64Observer, err = mHandler.meter.NewInt64ValueObserver(readyPrefix, observer(s))
		if err != nil {
			return err
		}
	}
	mHandler.metricValueObserver[readyPrefix] = int64Observer

	return nil
}

func observer(s *health.State) func(ctx context.Context, result metric.Int64ObserverResult) {
	return func(ctx context.Context, result metric.Int64ObserverResult) {
		var isReady int64
		if s.IsReady() {
			isReady = 1
		}

		result.Observe(isReady)
	}
}

func addVersionMetric(ctx context.Context, ro *registerOptions, s *health.State) error {
	var err error
	mHandlerMutex.Lock()
	defer mHandlerMutex.Unlock()
	versionPrefix := ro.metricPrefix + "version"
	int64Counter := mHandler.metricCounters[versionPrefix]
	if int64Counter == nil {
		int64Counter, err = mHandler.meter.NewInt64Counter(versionPrefix)
		if err != nil {
			return err
		}
	}

	int64Counter.Add(ctx, int64(1),
		CommitHash.String(s.Version.CommitHash),
		BuildLogURL.String(s.Version.BuildLogUrl),
		ContainerTag.String(s.Version.ContainerTag),
		RepoURL.String(s.Version.RepoUrl),
		Semver.String(s.Version.Semver),
	)

	for k, v := range ro.constLabels {
		int64Counter.Add(ctx, int64(1), k.String(v.AsString()))
	}

	mHandler.metricCounters[versionPrefix] = int64Counter

	return nil
}
