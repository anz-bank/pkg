package log

type configKey struct{}

const (
	formatter = iota
	standardFormatter
	jsonFormatter
)

type config interface {
	getConfig() int
	getConfigType() int
}

// StandardFormatter sets the logger to log with a standard format
type StandardFormatter struct{}

// JSONFormatter sets the logger to log with a JSON format
type JSONFormatter struct{}

func (StandardFormatter) getConfigType() int { return formatter }
func (StandardFormatter) getConfig() int     { return standardFormatter }

func (JSONFormatter) getConfigType() int { return formatter }
func (JSONFormatter) getConfig() int     { return jsonFormatter }
