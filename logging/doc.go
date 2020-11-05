/*
Package logging handles your application logs. Logging is built upon the very
popular 'zerolog' package and extends it by making it easier to manage logging
context.

This module provides a thin wrapper on top of zerolog that allows you to set up
contextual logging, and code links We expose some of the zerolog API for adding
fields to logs, so it can help to understand zerolog.

https://github.com/rs/zerolog

Logging API

Logging is built around the zerolog.Event type. To log something, invoke a log
event through one of the base log functions or their equivalents on a logger
object.

 logging.Info
 logging.Debug
 logging.Error

Any fields defined in context will already be added to this event, but you can
add more. See the zerolog.Event docs or examples. To close off the event and
log a message, use the Msg method. A complete example might look like.

 logging.Error(ctx, err).             // Invoke a log event
   Str("failure_type", "validation"). // Add a field
   Msg("Request failed")              // execute log with msg

JSON logs

By default the logging package produces logs in JSON format. Every piece of log
information has an associated key, and you can define nested objects and arrays

 {
   "time": "timestamp",
   "level": "info",
   "message": "Hello World",
   "service": {
     "name": "service",
     "version": "v1.0.0"
   },
   "tags": ["t1", "t2"]
 }

Contextual logging

Contextual logging the the key feature this library provides over zerolog. It
refers to the ability to create loggers with values takes dynamically from
context. The key to doing this is using this libraries Context object.

 type Context struct {
 	logger *Logger
 	funcs  []ContextFunc
 }

This object acts as a go-between for loggers and context. Create a logging
context by adding ContextFuncs to a logger. You can then either add it to a
real context or create a new logger with values loaded from context.

 logCtx := logger.With(ContextFunc...)
 ctx = logCtx.ToContext(ctx)      // adds logging context to context
 logger = logCtx.FromContext(ctx) // loads values from context to create a logger with context values set

ToContext allows you to use the package log functions later without worrying
about passing values in or calling the correct function to retrieve values
before your log.

 ctx = logCtx.ToContext(ctx)
 logging.Info(ctx).Msg("Hello World")

Duplicate fields

Zerolog provides zero protection from duplicate fields. If you are not careful,
you may accidentally run into this problem. This package provides one level of
protection from field duplication, identifiers for context funcs. Each context
func has an identifier that the logger uses to discover duplicate.

 logger := logger.With(Str("foo", "bar"), Str("foo", "baz"))
 // {"foo","bar"} no duplicate

This does NOT completely protect from duplication. Fields added after context
information has been collected CAN still result in duplications

 logger := logger.With(Str("foo", "bar"))
 logger.Info(ctx).Str("foo", "baz").Msg("Hello World")
 // {"foo":"bar","foo":"baz","msg":"Hello World"}

Custom context funcs CAN also still accidentally duplicate fields if they set
the same fields but have different function identifiers. There are two rules to
help avoid this

1. All common packages that define context funcs MUST document the fields they
set.

 // MyContextFunc sets fields to log xyz
 //
 // Field abc - short description
 //
 // Field xyz - short description
 // ...

2. Make sure you context funcs declare the same keys it intends to set, its
easy for typos to creep in.
*/
package logging
