package otelhealth

// nolint:lll
//go:generate go run -mod=mod github.com/golang/mock/mockgen -build_flags=-mod=mod -destination=testdata/mocks/int64_counter.go -package=mocks github.com/anz-bank/pkg/health/otelhealth Int64Counter
//go:generate go run -mod=mod github.com/golang/mock/mockgen -build_flags=-mod=mod -destination=testdata/mocks/int64_observer_result.go -package=mocks github.com/anz-bank/pkg/health/otelhealth Int64ObserverResult

import (
	"testing"

	"github.com/anz-bank/pkg/health"
	"github.com/anz-bank/pkg/health/otelhealth/testdata/mocks"
	otelAttribute "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
)

func TestRegisterWithValidValues(t *testing.T) {
	defer resetGlobals()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	versionCounter := mocks.NewMockInt64Counter(ctrl)

	mHandler = &metricHandler{
		meter:               global.Meter(""),
		metricCounters:      map[string]Int64Counter{"myapp_version": versionCounter},
		metricGaugeObserver: map[string]metric.Int64GaugeObserver{"myapp_ready": {}},
	}

	versionCounter.EXPECT().Add(gomock.Any(), gomock.Any(),
		CommitHash.String("1ee4e1f233caea38d6e331299f57dd86efb47361"),
		BuildLogURL.String("https://github.com/anz-bank/pkg/actions/runs/818181"),
		ContainerTag.String("gcr.io/google-containers/v1.0.0"),
		RepoURL.String("http://github.com/anz-bank/pkg"),
		Semver.String("v0.0.0"),
	)

	health.RepoURL = "http://github.com/anz-bank/pkg"
	health.CommitHash = "1ee4e1f233caea38d6e331299f57dd86efb47361"
	health.BuildLogURL = "https://github.com/anz-bank/pkg/actions/runs/818181"
	health.ContainerTag = "gcr.io/google-containers/v1.0.0"
	health.Semver = "v0.0.0"

	server, err := health.NewGRPCServer()
	require.NoError(t, err)

	err = Register(server.State, WithPrefix("myapp"))

	require.NoError(t, err)
}

func TestRegisterWithNewLabels(t *testing.T) {
	defer resetGlobals()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	versionCounter := mocks.NewMockInt64Counter(ctrl)

	mHandler = &metricHandler{
		meter:               global.Meter(""),
		metricCounters:      map[string]Int64Counter{"myapp_version": versionCounter},
		metricGaugeObserver: map[string]metric.Int64GaugeObserver{"myapp_ready": {}},
	}

	labels := map[otelAttribute.Key]otelAttribute.Value{
		otelAttribute.Key("foo"):  otelAttribute.Key("foo").String("bar").Value,
		otelAttribute.Key("test"): otelAttribute.Key("test").String("result").Value,
	}

	versionCounter.EXPECT().Add(gomock.Any(), gomock.Any(),
		CommitHash.String("1ee4e1f233caea38d6e331299f57dd86efb47361"),
		BuildLogURL.String("https://github.com/anz-bank/pkg/actions/runs/818181"),
		ContainerTag.String("gcr.io/google-containers/v1.0.0"),
		RepoURL.String("http://github.com/anz-bank/pkg"),
		Semver.String("v0.0.0"),
	)

	versionCounter.EXPECT().Add(gomock.Any(), gomock.Any(),
		otelAttribute.Key("foo").String("bar"))

	versionCounter.EXPECT().Add(gomock.Any(), gomock.Any(),
		otelAttribute.Key("test").String("result"))

	health.RepoURL = "http://github.com/anz-bank/pkg"
	health.CommitHash = "1ee4e1f233caea38d6e331299f57dd86efb47361"
	health.BuildLogURL = "https://github.com/anz-bank/pkg/actions/runs/818181"
	health.ContainerTag = "gcr.io/google-containers/v1.0.0"
	health.Semver = "v0.0.0"

	server, err := health.NewGRPCServer()
	require.NoError(t, err)

	err = Register(server.State, WithPrefix("myapp"), WithConstLabels(labels))

	require.NoError(t, err)
}

func TestRegisterWithUndefinedValues(t *testing.T) {
	defer resetGlobals()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	versionCounter := mocks.NewMockInt64Counter(ctrl)

	mHandler = &metricHandler{
		meter:               global.Meter(""),
		metricCounters:      map[string]Int64Counter{"myapp_version": versionCounter},
		metricGaugeObserver: map[string]metric.Int64GaugeObserver{"myapp_ready": {}},
	}

	versionCounter.EXPECT().Add(gomock.Any(), gomock.Any(),
		CommitHash.String("undefined"),
		BuildLogURL.String("undefined"),
		ContainerTag.String("undefined"),
		RepoURL.String("undefined"),
		Semver.String("undefined"),
	)

	s, err := health.NewState()
	require.NoError(t, err)

	err = Register(s, WithPrefix("myapp"))
	require.NoError(t, err)
}

func resetGlobals() {
	health.RepoURL = "undefined"
	health.CommitHash = "undefined"
	health.BuildLogURL = "undefined"
	health.ContainerTag = "undefined"
	health.Semver = "undefined"
}
