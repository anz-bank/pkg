// Package ochealth exports health.State data via OpenCensus.
//
// ochealth allows you to take the state from a health server and export the
// data in it via OpenCensus metrics. For example, to track the readiness of
// your server over time you can use OpenCensus to export via Prometheus metrics
// where you can graph and alert on your service not being ready. It also
// allows you to export version information so you can easily correlate
// version changes with changes in other metrics making it easier to identify
// regressions.
//
// The servers in github.com/anz-bank/pkg/health all have a State field. This
// can be passed to the Register function in this package to have that state
// provided as OpenCensus metrics.
//
// Two metrics are published by this package:
// - anz_health_ready
// - anz_health_version
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
// Additional labels can be added to the metrics with the WithConstLabels Option.
// NOTE: This currently panics due to a bug in the OpenCensus Go library. The
// panic will be removed when
// https://github.com/census-instrumentation/opencensus-go/pull/1221
// is merged and released.
package ochealth

import (
	"github.com/anz-bank/pkg/health"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricproducer"
)

type registerOptions struct {
	metricPrefix string
	constLabels  map[metricdata.LabelKey]metricdata.LabelValue
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
func WithConstLabels(labels map[metricdata.LabelKey]metricdata.LabelValue) Option {
	panic("Broken until https://github.com/census-instrumentation/opencensus-go/pull/1221 is merged and released")
	//nolint:govet
	// unreachable code
	return func(ro *registerOptions) {
		for k, v := range labels {
			ro.constLabels[k] = v
		}
	}
}

// Register creates a metrics registry with values from the given health.State
// and adds it to the global metric producer to be available to all the metric
// exporters. The metrics registry tracks changes to the non-constant values
// in the health.State.
func Register(s *health.State, options ...Option) error {
	ro := newRegisterOptions(options...)
	r := metric.NewRegistry()
	if err := addMetrics(r, ro, s); err != nil {
		return err
	}
	metricproducer.GlobalManager().AddProducer(r)
	return nil
}

func newRegisterOptions(options ...Option) *registerOptions {
	ro := &registerOptions{
		metricPrefix: "anz_health",
		constLabels:  map[metricdata.LabelKey]metricdata.LabelValue{},
	}
	for _, option := range options {
		option(ro)
	}
	if ro.metricPrefix != "" {
		ro.metricPrefix += "_"
	}
	return ro
}

func addMetrics(r *metric.Registry, ro *registerOptions, s *health.State) error {
	if err := addReadyMetric(r, ro, s); err != nil {
		return err
	}
	if err := addVersionMetric(r, ro, s); err != nil {
		return err
	}
	return nil
}

func addReadyMetric(r *metric.Registry, ro *registerOptions, s *health.State) error {
	g, err := r.AddInt64DerivedGauge(ro.metricPrefix+"ready",
		metric.WithDescription("Readiness state of server"),
		metric.WithUnit(metricdata.UnitDimensionless),
		metric.WithConstLabel(ro.constLabels))
	if err != nil {
		return err
	}
	return g.UpsertEntry(func() int64 {
		if s.IsReady() {
			return 1
		}
		return 0
	})
}

func addVersionMetric(r *metric.Registry, ro *registerOptions, s *health.State) error {
	labels := map[metricdata.LabelKey]metricdata.LabelValue{
		// note: commit_hash and build_log_url are deliberately out of order as
		// the metrics library sorts the labels and we want to discover any
		// regression introduced with maintaining the key/value mapping through
		// our tests.
		{Key: "commit_hash"}:   metricdata.NewLabelValue(s.Version.CommitHash),
		{Key: "build_log_url"}: metricdata.NewLabelValue(s.Version.BuildLogUrl),
		{Key: "container_tag"}: metricdata.NewLabelValue(s.Version.ContainerTag),
		{Key: "repo_url"}:      metricdata.NewLabelValue(s.Version.RepoUrl),
		{Key: "semver"}:        metricdata.NewLabelValue(s.Version.Semver),
	}
	for k, v := range ro.constLabels {
		labels[k] = v
	}
	g, err := r.AddInt64Gauge(ro.metricPrefix+"version",
		metric.WithDescription("Version information"),
		metric.WithUnit(metricdata.UnitDimensionless),
		metric.WithConstLabel(labels))
	if err != nil {
		return err
	}
	entry, err := g.GetEntry()
	if err != nil {
		return err
	}
	entry.Set(1)
	return nil
}
