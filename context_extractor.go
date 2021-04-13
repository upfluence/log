package log

import (
	"context"

	"github.com/upfluence/log/record"
)

type ContextExtractor interface {
	Extract(context.Context, record.Level) []record.Field
}

type noopExtractor struct{}

func (noopExtractor) Extract(context.Context, record.Level) []record.Field {
	return nil
}

type multiContextExtractor []ContextExtractor

func (ces multiContextExtractor) Extract(ctx context.Context, lvl record.Level) []record.Field {
	var fs []record.Field

	for _, ce := range ces {
		fs = append(fs, ce.Extract(ctx, lvl)...)
	}

	return fs
}

func mergeContextExtractor(lce, rce ContextExtractor) ContextExtractor {
	switch tlce := lce.(type) {
	case noopExtractor:
		return rce
	case multiContextExtractor:
		return multiContextExtractor(append(tlce, rce))
	default:
		return multiContextExtractor([]ContextExtractor{lce, rce})
	}
}

type ContextExtractorFunc func(context.Context, record.Level) []record.Field

func (fn ContextExtractorFunc) Extract(ctx context.Context, lvl record.Level) []record.Field {
	return fn(ctx, lvl)
}

type leveledContextExtractor struct {
	ce  ContextExtractor
	lvl record.Level
}

func LeveledContextExtractor(ce ContextExtractor, lvl record.Level) ContextExtractor {
	return &leveledContextExtractor{ce: ce, lvl: lvl}
}

func (lce *leveledContextExtractor) Extract(ctx context.Context, lvl record.Level) []record.Field {
	if lvl < lce.lvl {
		return nil
	}

	return lce.ce.Extract(ctx, lvl)
}
