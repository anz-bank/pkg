package log

type configKey int

const (
	formatter configKey = iota
	standardFormatter
	jsonFormatter
)

type config interface {
	getConfig() configKey
	getConfigType() configKey
}

// StandardFormatter sets the logger to log with a standard format
type StandardFormatter struct{}

// JSONFormatter sets the logger to log with a JSON format
type JSONFormatter struct{}

func (StandardFormatter) getConfigType() configKey { return formatter }
func (StandardFormatter) getConfig() configKey     { return standardFormatter }

func (JSONFormatter) getConfigType() configKey { return formatter }
func (JSONFormatter) getConfig() configKey     { return jsonFormatter }
