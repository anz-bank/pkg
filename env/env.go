// Package env is a context-driven wrapper for access to environment variable.
// It allows substitution of mock environments via context.Context for testing
// and other purposes.
package env

import (
	"context"
	"os"
)

// Env models the envvar functions of the os package.
type Env interface {
	Clearenv()
	Environ() []string
	LookupEnv(key string) (string, bool)
	Setenv(key, value string) error
	Unsetenv(key string) error
}

type envKey struct{}

// Onto sets the context Env. Pass nil to revert to the default system env.
func Onto(ctx context.Context, env Env) context.Context {
	return context.WithValue(ctx, envKey{}, env)
}

// From gets the Env from the Context.
func From(ctx context.Context) Env {
	if env := ctx.Value(envKey{}); env != nil {
		return env.(Env)
	}
	return defaultEnv{}
}

// The following functions mirror their counterparts in the os package.

func Clearenv(ctx context.Context) {
	From(ctx).Clearenv()
}

func Environ(ctx context.Context) []string {
	return From(ctx).Environ()
}

func getenv(env Env, key string) string {
	if value, has := env.LookupEnv(key); has {
		return value
	}
	return ""
}

func ExpandEnv(ctx context.Context, s string) string {
	env := From(ctx)
	getenv := func(key string) string {
		return getenv(env, key)
	}
	return os.Expand(s, getenv)
}

func Getenv(ctx context.Context, key string) string {
	return getenv(From(ctx), key)
}

func LookupEnv(ctx context.Context, key string) (string, bool) {
	return From(ctx).LookupEnv(key)
}

func Setenv(ctx context.Context, key, value string) error {
	return From(ctx).Setenv(key, value)
}

func Unsetenv(ctx context.Context, key string) error {
	return From(ctx).Unsetenv(key)
}
