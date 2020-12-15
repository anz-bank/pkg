package logging

import "os"

// gblLogContext is the global log context used when a log context can't be found in context
//
// Usage of the global log context is usually a bug but can be useful for logging in a
// function where the context doesn't matter or for quick testing/debugging purposes.
var gblContext = New(os.Stdout).WithBool("globalLogger", true).With()

// SetGlobalLogContext sets the global logger to the given logger
//
// Allows adding static application values to a logger
func SetGlobalLogContext(logCtx *Context) {
	gblContext = logCtx.WithBool("globalLogger", true)
}
