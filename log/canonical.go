package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

// OnLog is a callback that is called when fields and configurations are added to the logger.
func OnLog(ctx context.Context, fields Fields) {
	if fields := ctx.Value(canonicalFieldsKey{}); fields == nil {
		Info(ctx, "Canonical fields have not been set up properly, fields will not be collected")
		return
	}
	addCanonicalFields(getCanonicalFields(ctx), fields)
}

// CanonicalLog logs all the collected fields from previous logs in non-verbose mode.
func CanonicalLog(ctx context.Context) {
	from(ctx, Fields{getCanonicalFields(ctx).Finish()}).Info()
}

// WithCanonicalLog creates a canonical fields collector.
func WithCanonicalLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, canonicalFieldsKey{}, frozen.NewMapBuilder(0))
}
