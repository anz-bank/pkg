package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	// Do not parallelise.

	ctx := context.Background()

	// Don't test Clearenv. It's hazardous.
	assert.NotEmpty(t, Environ(ctx))

	assert.NotEmpty(t, ExpandEnv(ctx, "$PATH"))
	assert.NotEmpty(t, ExpandEnv(ctx, "${PATH}"))
	assert.Equal(t, "NOT_THERE = ", ExpandEnv(ctx, "NOT_THERE = ${NOT_THERE}"))

	assert.NotEmpty(t, Getenv(ctx, "PATH"))
	assert.Empty(t, Getenv(ctx, "THIS_ENVVAR_SHOULD_NOT_EXIST"))

	user, has := LookupEnv(ctx, "PATH")
	if assert.True(t, has) {
		assert.NotEmpty(t, user)
	}
	shouldNotExist, has := LookupEnv(ctx, "THIS_ENVVAR_SHOULD_NOT_EXIST")
	assert.False(t, has, shouldNotExist)

	fbb := "FOO_BAR_BAZ"
	shouldNotExist, has = LookupEnv(ctx, fbb)
	assert.False(t, has, shouldNotExist)
	if assert.NoError(t, Setenv(ctx, fbb, "42")) {
		assert.Equal(t, "42", Getenv(ctx, fbb))
		value, has := LookupEnv(ctx, fbb)
		if assert.True(t, has, shouldNotExist) {
			assert.Equal(t, "42", value)
		}
	}
	if assert.NoError(t, Setenv(ctx, fbb, "")) {
		assert.Equal(t, "", Getenv(ctx, fbb))
		value, has := LookupEnv(ctx, fbb)
		if assert.True(t, has, shouldNotExist) {
			assert.Empty(t, value)
		}
	}
	if assert.NoError(t, Unsetenv(ctx, fbb)) {
		assert.Equal(t, "", Getenv(ctx, fbb))
		shouldNotExist, has := LookupEnv(ctx, fbb)
		assert.False(t, has, shouldNotExist)
	}
}
