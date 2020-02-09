package log

type typeKey int

const (
	Formatter typeKey = iota
	LevelSetter
	OutSetter
)

type Config interface {
	TypeKey() interface{}
	Apply(logger Logger) error
}

type standardFormat struct{}
type jsonFormat struct{}

type errorLevel struct{}
type debugLevel struct{}
type infoLevel struct{}

type stderrOut struct{}
type stdOut struct{}
type bufferOut struct{}

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

func NewErrorLevel() Config             { return errorLevel{} }
func (errorLevel) TypeKey() interface{} { return LevelSetter }
func (el errorLevel) Apply(logger Logger) error {
	return applyLevel(el, logger)
}

func NewDebugLevel() Config             { return debugLevel{} }
func (debugLevel) TypeKey() interface{} { return LevelSetter }
func (dl debugLevel) Apply(logger Logger) error {
	return applyLevel(dl, logger)
}

func NewInfoLevel() Config             { return infoLevel{} }
func (infoLevel) TypeKey() interface{} { return LevelSetter }
func (il infoLevel) Apply(logger Logger) error {
	return applyLevel(il, logger)
}

func NewStderrOut() Config             { return stderrOut{} }
func (stderrOut) TypeKey() interface{} { return OutSetter }
func (se stderrOut) Apply(logger Logger) error {
	return applyOutput(se, logger)
}

func NewStdOut() Config             { return stdOut{} }
func (stdOut) TypeKey() interface{} { return OutSetter }
func (s stdOut) Apply(logger Logger) error {
	return applyOutput(s, logger)
}

func NewBufferOut() Config             { return bufferOut{} }
func (bufferOut) TypeKey() interface{} { return OutSetter }
func (b bufferOut) Apply(logger Logger) error {
	return applyOutput(b, logger)
}

func applyFormatter(formatter Config, logger Logger) error {
	return logger.(formattable).SetFormatter(formatter)
}

func applyLevel(level Config, logger Logger) error {
	return logger.(settableLevel).SetLevel(level)
}

func applyOutput(out Config, logger Logger) error {
	return logger.(settableOutput).SetOutput(out)
}
