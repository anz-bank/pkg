package ochealth

import (
	"context"
	"testing"

	"github.com/anz-bank/pkg/health"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	"go.opencensus.io/metric/metricproducer"
)

// metrics captures the global metrics by being a metricexport.Exporter and
// recording the []metricdata.Metric passed to it. The contents of the slice
// are not copied, so it should be used right away and not concurrently.
// It is only for the tests in this file.
type metrics struct {
	data []*metricdata.Metric
}

// readMetrics exports the current OpenCensus metrics from the producers in
// the global manager.
func readMetrics() *metrics {
	r := metricexport.NewReader()
	m := &metrics{}
	r.ReadAndExport(m)
	return m
}

// ExportMetrics implements metricexport.Exporter.
func (m *metrics) ExportMetrics(ctx context.Context, data []*metricdata.Metric) error {
	m.data = data
	return nil
}

// requireValue requires that a metric with a given name exists with the given value.
func (m *metrics) requireValue(t *testing.T, metric string, value interface{}) {
	for _, m := range m.data {
		if m.Descriptor.Name == metric {
			require.Equal(t, value, m.TimeSeries[0].Points[0].Value)
			return
		}
	}
	require.Fail(t, "metric not present", metric)
}

// requireLabelValue requires that a metric exists with a given label and label value.
func (m *metrics) requireLabelValue(t *testing.T, metric, label, value string) {
	for _, m := range m.data {
		if m.Descriptor.Name == metric {
			require.Equal(t, value, m.TimeSeries[0].LabelValues[0].Value)
			require.True(t, m.TimeSeries[0].LabelValues[0].Present)
			return
		}
	}
	require.Fail(t, "metric not present", metric)
}

func TestRegister(t *testing.T) {
	s, err := health.NewState()
	require.NoError(t, err)

	err = Register(s)
	require.NoError(t, err)
	defer deleteAllProducers()

	m := readMetrics()
	m.requireValue(t, "anz_health_ready", int64(0))
	m.requireValue(t, "anz_health_version", int64(1))
	m.requireLabelValue(t, "anz_health_version", "build_log_url", health.BuildLogURL)
	m.requireLabelValue(t, "anz_health_version", "commit_hash", health.CommitHash)
	m.requireLabelValue(t, "anz_health_version", "container_tag", health.ContainerTag)
	m.requireLabelValue(t, "anz_health_version", "repo_url", health.RepoURL)
	m.requireLabelValue(t, "anz_health_version", "semver", health.Semver)

	s.SetReady(true)
	m = readMetrics()
	m.requireValue(t, "anz_health_ready", int64(1))
}

func TestRegisterNoPrefix(t *testing.T) {
	s, err := health.NewState()
	require.NoError(t, err)

	err = Register(s, WithPrefix(""))
	require.NoError(t, err)
	defer deleteAllProducers()

	m := readMetrics()
	m.requireValue(t, "ready", int64(0))
	m.requireValue(t, "version", int64(1))
}

func TestRegisterCustomPrefix(t *testing.T) {
	s, err := health.NewState()
	require.NoError(t, err)

	err = Register(s, WithPrefix("xplore"))
	require.NoError(t, err)
	defer deleteAllProducers()

	m := readMetrics()
	m.requireValue(t, "xplore_ready", int64(0))
	m.requireValue(t, "xplore_version", int64(1))
}

func TestRegisterConstLabels(t *testing.T) {
	s, err := health.NewState()
	require.NoError(t, err)

	labels := map[metricdata.LabelKey]metricdata.LabelValue{
		{Key: "foo"}: metricdata.NewLabelValue("bar"),
	}

	// Remove require.Panics and return when
	// https://github.com/census-instrumentation/opencensus-go/pull/1221
	// is merged and released.
	require.Panics(t, func() { _ = Register(s, WithPrefix(""), WithConstLabels(labels)) })
	return

	//nolint:govet
	// unreachable code
	err = Register(s, WithPrefix(""), WithConstLabels(labels))
	require.NoError(t, err)

	m := readMetrics()
	m.requireValue(t, "ready", int64(0))
	m.requireLabelValue(t, "ready", "foo", "bar")
	m.requireValue(t, "version", int64(1))
	m.requireLabelValue(t, "version", "foo", "bar")
	m.requireLabelValue(t, "version", "build_log_url", health.BuildLogURL)
	m.requireLabelValue(t, "version", "commit_hash", health.CommitHash)
	m.requireLabelValue(t, "version", "container_tag", health.ContainerTag)
	m.requireLabelValue(t, "version", "repo_url", health.RepoURL)
	m.requireLabelValue(t, "version", "semver", health.Semver)
}

func TestRegisterErrors(t *testing.T) {
	// We cannot directly test `Register` for errors as we need the
	// metric.Registry it creates to inject conflicting metrics to
	// generate an error. So call the addMetrics and descendent functions
	// directly.

	s, err := health.NewState()
	require.NoError(t, err)

	r := metric.NewRegistry()
	_, err = r.AddFloat64Gauge("anz_health_ready")
	require.NoError(t, err)
	err = addMetrics(r, newRegisterOptions(), s)
	require.Error(t, err)

	r = metric.NewRegistry()
	_, err = r.AddFloat64Gauge("anz_health_version")
	require.NoError(t, err)
	err = addMetrics(r, newRegisterOptions(), s)
	require.Error(t, err)
}

func deleteAllProducers() {
	for _, p := range metricproducer.GlobalManager().GetAll() {
		metricproducer.GlobalManager().DeleteProducer(p)
	}
}
