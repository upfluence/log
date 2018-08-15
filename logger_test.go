package log

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/upfluence/log/record"
	"github.com/upfluence/log/sink/leveled"
	"github.com/upfluence/log/sink/writer"
)

func logBenchmark(b *testing.B, fn func(Logger, int), opts []LoggerOption) {
	var l = NewLogger(opts...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fn(l, i)
	}
}

type staticFormatter struct{}

func (staticFormatter) Format(io.Writer, record.Record) error { return nil }

type staticSink struct{}

func (staticSink) Log(record.Record) error { return nil }

type loggerBench struct {
	name string
	opts []LoggerOption
}

var benches = []loggerBench{
	loggerBench{
		name: "standard with stack trace",
		opts: []LoggerOption{WithSink(writer.NewStandardSink(ioutil.Discard))},
	},
	loggerBench{
		name: "standard with no stack trace",
		opts: []LoggerOption{
			WithSink(writer.NewSink(writer.NewFastFormatter(), ioutil.Discard)),
		},
	},
	loggerBench{
		name: "static formatter",
		opts: []LoggerOption{
			WithSink(writer.NewSink(staticFormatter{}, ioutil.Discard)),
		},
	},
	loggerBench{
		name: "static sink",
		opts: []LoggerOption{WithSink(staticSink{})},
	},
	loggerBench{
		name: "leveled",
		opts: []LoggerOption{
			WithSink(
				leveled.NewSink(record.Notice, writer.NewStandardSink(ioutil.Discard)),
			),
		},
	},
}

func BenchmarkFixedString(b *testing.B) {
	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			logBenchmark(
				b,
				func(l Logger, _ int) { l.Info("foo bar") },
				bench.opts,
			)
		})
	}
}

func BenchmarkFmtString(b *testing.B) {
	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			logBenchmark(
				b,
				func(l Logger, _ int) { l.Infof("foo bar %v", "foo") },
				bench.opts,
			)
		})
	}
}
