## errs

Convenience functions for dealing with Go errors taken
from https://github.com/anzx/vcards. There are two main
features of this package that distinguish it from the
standard Go error handling:
1. It uses a tree structure to store the errors (see `multiErr`).
This allows error from multiple paths to be returned in a way that
maintains the originating path information.
2. Allows the error message to be captured but avoid wrapping the
error through `NoWrap`. This allows errors to be logged without
exposing implementation details to code.
See https://blog.golang.org/go1.13-errors

### New

Creating a new error that wraps one or more other errors.

Example usage:

```Go
var ErrMyError = fmt.Errorf("my error")

resp, err := Function(ctx, req)
if err != nil {
    return nil, errs.New(e.ErrMyError, err)
}
```

`New` honours `NoWrap` (see below).

### NoWrap

To mark an error as not for wrapping but for the error message to
be captured:
```Go
var ErrMyError = fmt.Errorf("my error")

resp, err := Function(ctx, req)
if err != nil {
    return nil, errs.New(e.ErrMyError, NoWrap(err))
}
```

### Errorf

As per `fmt.Errorf` but honours `NoWrap` markings.

### CloseIgnoreErr

Use `errs.CloseIgnoreErr(c)` instead of:
```Go
//nolint:errcheck
func myfunc(c io.Closer) {
    c.Close()
}
```

or:

```Go
func myfunc(c io.Closer) {
    defer func() {_ = c.Close()}()
}
```
