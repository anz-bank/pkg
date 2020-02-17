package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

func OnInfo(ctx context.Context, fields Fields) {
	addCanonicalFields(getCanonicalFields(ctx), fields)
}

func OnDebug(ctx context.Context, fields Fields) {
	addCanonicalFields(getCanonicalFields(ctx), fields)
}

func OnError(ctx context.Context, fields Fields) {
	addCanonicalFields(getCanonicalFields(ctx), fields)
}

func RegisterCanonicalListener(ctx context.Context, callback ...func(context.Context, Fields)) context.Context{
	callbacks := frozen.NewSetBuilder(len(callback))
	for _, c := range callback {
		callbacks.Add(c)
	}
	return context.WithValue(
		context.WithValue(ctx, canonicalFieldsKey{}, frozen.NewMapBuilder(0)),
		canonicalListenerKey{},
		callbacks.Finish(),
	)
}

func CanonicalInfo(ctx context.Context) {
	from(ctx, Fields{getCanonicalFields(ctx).Finish()}).Info(ctx)
}

func CanonicalDebug(ctx context.Context) {
	from(ctx, Fields{getCanonicalFields(ctx).Finish()}).Debug(ctx)
}

func CanonicalError(ctx context.Context, errMsg error) {
	from(ctx, Fields{getCanonicalFields(ctx).Finish()}).Error(errMsg)
}
