package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Parallel()

	ctx := Onto(context.Background(), NewMap(map[string]string{
		"CROMULENCE": "cromulent",
		"PATH":       "/bin:/usr/bin:/usr/local/bin",
	}))

	// Don't test Clearenv. It's hazardous.
	assert.NotEmpty(t, Environ(ctx))

	assert.Equal(t, "CROMULENCE = cromulent", ExpandEnv(ctx, "CROMULENCE = $CROMULENCE"))
	assert.Equal(t, "CROMULENCE = cromulent", ExpandEnv(ctx, "CROMULENCE = ${CROMULENCE}"))
	assert.Equal(t, "NOT_THERE = ", ExpandEnv(ctx, "NOT_THERE = ${NOT_THERE}"))

	assert.Equal(t, "cromulent", Getenv(ctx, "CROMULENCE"))
	assert.Empty(t, Getenv(ctx, "THIS_ENVVAR_SHOULD_NOT_EXIST"))

	user, has := LookupEnv(ctx, "CROMULENCE")
	if assert.True(t, has) {
		assert.Equal(t, "cromulent", user)
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

	Clearenv(ctx)
	assert.Empty(t, Environ(ctx))
}

func TestTwoMaps(t *testing.T) {
	t.Parallel()

	ctx1 := Onto(context.Background(), NewMap(map[string]string{
		"NAME": "Mog",
	}))

	ctx2 := Onto(context.Background(), NewMap(map[string]string{
		"COLOR": "blue",
	}))

	assert.Equal(t, "Mog", Getenv(ctx1, "NAME"))
	assert.Empty(t, Getenv(ctx2, "NAME"))
	assert.Empty(t, Getenv(ctx1, "COLOR"))
	assert.Equal(t, "blue", Getenv(ctx2, "COLOR"))

	assert.NoError(t, Unsetenv(ctx1, "NAME"))
	assert.NoError(t, Unsetenv(ctx2, "NAME"))
	assert.Empty(t, Getenv(ctx1, "NAME"))
	assert.NotEmpty(t, Getenv(ctx2, "COLOR"))
}
