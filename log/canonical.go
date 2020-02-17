package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

func OnLog(ctx context.Context, fields Fields) {
	addCanonicalFields(getCanonicalFields(ctx), fields)
}

func CanonicalLog(ctx context.Context) {
	from(ctx, Fields{getCanonicalFields(ctx).Finish()}).Info()
}

func StartCanonicalLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, canonicalFieldsKey{}, frozen.NewMapBuilder(0))
}
