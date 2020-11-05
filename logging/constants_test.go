package logging_test

import (
	"testing"

	logging "github.com/anz-bank/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	info := logging.MustParseLevel("info")
	assert.Equal(t, logging.InfoLevel, info)

	debug := logging.MustParseLevel("debug")
	assert.Equal(t, logging.DebugLevel, debug)

	errorLevel := logging.MustParseLevel("error")
	assert.Equal(t, logging.ErrorLevel, errorLevel)

	assert.Panics(t, func() { _ = logging.MustParseLevel("unknown") })
}

func TestLevel_String(t *testing.T) {
	assert.Equal(t, "info", logging.InfoLevel.String())
}
