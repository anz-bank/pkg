package log

import "io"

type typeKey int
type internalTypeKey int

const (
	FormatterType typeKey = iota
)

const (
	verbosity internalTypeKey = iota
	output
	addHooks
	logCaller
)

type Config interface {
	TypeKey() interface{}
	Apply(logger Logger) error
}

type standardFormat struct{}
type jsonFormat struct{}

type verboseMode struct{ on bool }
type outputConfig struct{ writer io.Writer }
type addHooksConfig struct{ hooks []Hook }
type logCallerConfig struct{ on bool }

func NewStandardFormat() Config             { return standardFormat{} }
func (standardFormat) TypeKey() interface{} { return FormatterType }
func (sf standardFormat) Apply(logger Logger) error {
	return applyFormatter(sf, logger)
}

func NewJSONFormat() Config             { return jsonFormat{} }
func (jsonFormat) TypeKey() interface{} { return FormatterType }
func (jf jsonFormat) Apply(logger Logger) error {
	return applyFormatter(jf, logger)
}

func SetVerboseMode(on bool) Config {
	return verboseMode{on}
}

func (verboseMode) TypeKey() interface{} { return verbosity }

func (v verboseMode) Apply(logger Logger) error {
	return logger.(SettableVerbosity).SetVerbose(v.on)
}

func SetOutput(w io.Writer) Config {
	return outputConfig{w}
}

func (outputConfig) TypeKey() interface{} { return output }

func (o outputConfig) Apply(logger Logger) error {
	return logger.(SettableOutput).SetOutput(o.writer)
}

// AddHooks adds the given hooks to the logger.
func AddHooks(hooks ...Hook) Config {
	return addHooksConfig{hooks}
}

func (addHooksConfig) TypeKey() interface{} { return addHooks }

func (o addHooksConfig) Apply(logger Logger) error {
	return logger.(AddableHooks).AddHooks(o.hooks...)
}

// SetLogCaller sets whether or not a reference to the calling function is logged.
func SetLogCaller(on bool) Config {
	return logCallerConfig{on}
}

func (logCallerConfig) TypeKey() interface{} { return logCaller }

func (c logCallerConfig) Apply(logger Logger) error {
	return logger.(SettableLogCaller).SetLogCaller(c.on)
}

func applyFormatter(formatter Config, logger Logger) error {
	return logger.(Formattable).SetFormatter(formatter)
}

type forwardingHook struct {
	logger Logger
}

func (h *forwardingHook) OnLogged(entry *LogEntry) error {
	h.logger.(entryLogger).Log(entry)
	return nil
}

// NewForwardingHook returns a Hook that forwards all entries to the given Logger.
func NewForwardingHook(logger Logger) Hook {
	return &forwardingHook{logger}
}
