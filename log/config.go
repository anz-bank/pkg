package log

import "io"

type typeKey int
type internalTypeKey int

const (
	Formatter typeKey = iota
)

const (
	verbosity internalTypeKey = iota
	output
)

type Config interface {
	TypeKey() interface{}
	Apply(logger Logger) error
}

type standardFormat struct{}
type jsonFormat struct{}

type verboseMode struct{ on bool }
type outputConfig struct{ writer io.Writer }

func NewStandardFormat() Config             { return standardFormat{} }
func (standardFormat) TypeKey() interface{} { return Formatter }
func (sf standardFormat) Apply(logger Logger) error {
	return applyFormatter(sf, logger)
}

func NewJSONFormat() Config             { return jsonFormat{} }
func (jsonFormat) TypeKey() interface{} { return Formatter }
func (jf jsonFormat) Apply(logger Logger) error {
	return applyFormatter(jf, logger)
}

func SetVerboseMode(on bool) Config {
	return verboseMode{on}
}

func (verboseMode) TypeKey() interface{} { return verbosity }

func (v verboseMode) Apply(logger Logger) error {
	return logger.(settableVerbosity).SetVerbose(v.on)
}

func SetOutput(w io.Writer) Config {
	return outputConfig{w}
}

func (outputConfig) TypeKey() interface{} { return output }

func (o outputConfig) Apply(logger Logger) error {
	return logger.(settableOutput).SetOutput(o.writer)
}

func applyFormatter(formatter Config, logger Logger) error {
	return logger.(formattable).SetFormatter(formatter)
}
