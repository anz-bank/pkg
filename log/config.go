package log

type typeKey int

const (
	Formatter typeKey = iota
	verbosity
	OutSetter
)

type Config interface {
	TypeKey() interface{}
	Apply(logger Logger) error
}

type standardFormat struct{}
type jsonFormat struct{}

type verboseMode struct{ on bool }

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

func applyFormatter(formatter Config, logger Logger) error {
	return logger.(formattable).SetFormatter(formatter)
}
