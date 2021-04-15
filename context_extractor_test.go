package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/log/record"
)

type staticContextExtractor struct {
	f record.Field
}

func (sce staticContextExtractor) Extract(context.Context, record.Level) []record.Field {
	return []record.Field{sce.f}
}

type recordSink struct {
	r record.Record
}

func (rs *recordSink) Log(r record.Record) error {
	rs.r = r
	return nil
}

type subContext struct {
	context.Context
}

func TestMultiExtractor(t *testing.T) {
	var (
		s recordSink

		l = NewLogger(
			WithContextExtractor(staticContextExtractor{f: Field("foo", "bar")}),
			WithContextExtractor(staticContextExtractor{f: Field("buz", "bar")}),
			WithContextExtractor(staticContextExtractor{f: Field("biz", "bar")}),
			WithSink(&s),
		)
	)

	l.Error("foo")
	assert.Equal(t, 0, len(s.r.Fields()))

	l.WithContext(subContext{}).Error("buz")
	assert.Equal(t, 3, len(s.r.Fields()))
}

func TestLeveledContextExtractor(t *testing.T) {
	var (
		s recordSink

		l = NewLogger(
			WithContextExtractor(
				LeveledContextExtractor(
					CombineContextExtractors(
						ContextExtractorFunc(func(context.Context, record.Level) []record.Field {
							return []record.Field{Field("foo", "bar")}
						}),
						ContextExtractorFunc(func(context.Context, record.Level) []record.Field {
							return []record.Field{Field("biz", "buz")}
						}),
					),
					record.Notice,
				),
			),
			WithSink(&s),
		)
	)

	l.WithContext(subContext{}).Error("buz")
	assert.Equal(
		t,
		[]record.Field{Field("foo", "bar"), Field("biz", "buz")},
		s.r.Fields(),
	)

	l.WithContext(subContext{}).Info("buz")
	assert.Equal(t, 0, len(s.r.Fields()))
}
