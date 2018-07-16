package log

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/upfluence/log/record"
	"github.com/upfluence/log/sink/writer"
)

func logBenchmark(b *testing.B, fn func(Logger, int), opts ...LoggerOption) {
	var l = NewLogger(opts...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fn(l, i)
	}
}

func BenchmarkFixedString(b *testing.B) {
	logBenchmark(
		b,
		func(l Logger, _ int) { l.Info("foo bar") },
		WithSink(writer.NewStandardSink(ioutil.Discard)),
	)
}

func BenchmarkFmtString(b *testing.B) {
	logBenchmark(
		b,
		func(l Logger, _ int) { l.Infof("foo bar %v", errors.New("foo")) },
		WithSink(writer.NewStandardSink(ioutil.Discard)),
	)
}

type staticFormatter struct{}

func (staticFormatter) Format(r record.Record) string { return "" }

func BenchmarkStaticFormatter(b *testing.B) {
	logBenchmark(
		b,
		func(l Logger, _ int) { l.Info("foo bar") },
		WithSink(writer.NewSink(staticFormatter{}, ioutil.Discard)),
	)
}
