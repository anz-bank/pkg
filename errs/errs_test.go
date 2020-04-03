package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorfIs(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")
	err3 := fmt.Errorf("error 3")
	err := Errorf("%v: %v", err1, err2)

	require.Equal(t, "error 1: error 2", err.Error())
	require.True(t, errors.Is(err, err1))
	require.True(t, errors.Is(err, err2))
	require.False(t, errors.Is(err, err3))
}

type errType1 string
type errType2 string
type errType3 string

func (e errType1) Error() string { return string(e) }
func (e errType2) Error() string { return string(e) }
func (e errType3) Error() string { return string(e) }

func TestErrorfAs(t *testing.T) {
	err1 := errType1("error 1")
	err2 := errType2("error 2")
	err := Errorf("%v: %v", err1, err2)

	require.Equal(t, "error 1: error 2", err.Error())
	var target1 errType1
	require.True(t, errors.As(err, &target1))
	require.Equal(t, err1, target1)
	var target2 errType2
	require.True(t, errors.As(err, &target2))
	require.Equal(t, err2, target2)
	var target3 errType3
	require.False(t, errors.As(err, &target3))
}

func TestUnwrapAlwaysNil(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err := Errorf("wrapping %v", err1)

	require.Equal(t, "wrapping error 1", err.Error())
	require.True(t, errors.Is(err, err1))
	require.Nil(t, errors.Unwrap(err))
}

func TestNoWrap(t *testing.T) {
	err1 := errType1("error 1")
	err2 := errType2("error 2")
	err := Errorf("%v: %v", NoWrap(err1), err2)

	require.Equal(t, "error 1: error 2", err.Error())
	require.False(t, errors.Is(err, err1))
	require.True(t, errors.Is(err, err2))
	var target1 errType1
	require.False(t, errors.As(err, &target1))
	var target2 errType2
	require.True(t, errors.As(err, &target2))
	require.Equal(t, err2, target2)
}

func TestNew(t *testing.T) {
	require.Nil(t, New())

	err1 := errType1("error 1")
	require.Equal(t, err1, New(err1))

	err2 := errType2("error 2")
	err := New(err1, err2)
	require.Error(t, err)
	require.Equal(t, "error 1: error 2", err.Error())

	err3 := errType2("error 3")
	err = New(err1, err2, err3)
	require.Error(t, err)
	require.Equal(t, "error 1: error 2: error 3", err.Error())
}

type closeChecker bool

func (cc *closeChecker) Close() error { *cc = true; return nil }

func TestCloseIgnoreErr(t *testing.T) {
	var cc closeChecker
	CloseIgnoreErr(&cc)
	require.True(t, bool(cc))
}
