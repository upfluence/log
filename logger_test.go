package log

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/upfluence/log/record"
	"github.com/upfluence/log/sink/leveled"
	"github.com/upfluence/log/sink/writer"
)

func logBenchmark(b *testing.B, fn func(Logger, int), lb loggerBench) {
	var l = lb.logger()

	b.ResetTimer()

	b.Run(lb.bname(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn(l, i)
		}
	})
}

type staticFormatter struct{}

func (staticFormatter) Format(io.Writer, record.Record) error { return nil }

type staticSink struct{}

func (staticSink) Log(record.Record) error { return nil }

type loggerBench interface {
	bname() string
	logger() Logger
}

type loggerBenchUpf struct {
	name string
	opts []LoggerOption
}

func (b loggerBenchUpf) bname() string  { return "upfluence/log/" + b.name }
func (b loggerBenchUpf) logger() Logger { return NewLogger(b.opts...) }

type loggerBenchLogrus struct {
	name string

	l logrus.FieldLogger
}

func (b loggerBenchLogrus) bname() string  { return "logrus/" + b.name }
func (b loggerBenchLogrus) logger() Logger { return &logrusAdapter{lr: b.l} }

type logrusAdapter struct {
	Logger

	lr logrus.FieldLogger
}

func (la *logrusAdapter) Info(vs ...interface{}) { la.lr.Info(vs...) }
func (la *logrusAdapter) Infof(fmt string, vs ...interface{}) {
	la.lr.Infof(fmt, vs...)
}
func (la *logrusAdapter) WithField(f record.Field) SugaredLogger {
	return &logrusAdapter{lr: la.lr.WithField(f.GetKey(), f.GetValue())}
}

var benches = []loggerBench{
	loggerBenchUpf{
		name: "standard with stack trace",
		opts: []LoggerOption{WithSink(writer.NewStandardSink(ioutil.Discard))},
	},
	loggerBenchLogrus{
		name: "standard Text with no stacktrace",
		l: &logrus.Logger{
			Out:       ioutil.Discard,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		},
	},
	loggerBenchLogrus{
		name: "standard JSON with no stacktrace",
		l: &logrus.Logger{
			Out:       ioutil.Discard,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		},
	},
	loggerBenchLogrus{
		name: "leved",
		l: &logrus.Logger{
			Out:       ioutil.Discard,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.WarnLevel,
		},
	},
	loggerBenchUpf{
		name: "standard with no stack trace",
		opts: []LoggerOption{
			WithSink(writer.NewSink(writer.NewFastFormatter(), ioutil.Discard)),
		},
	},
	loggerBenchUpf{
		name: "static formatter",
		opts: []LoggerOption{
			WithSink(writer.NewSink(staticFormatter{}, ioutil.Discard)),
		},
	},
	loggerBenchUpf{
		name: "static sink",
		opts: []LoggerOption{WithSink(staticSink{})},
	},
	loggerBenchUpf{
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
		logBenchmark(
			b,
			func(l Logger, _ int) { l.Info("foo bar") },
			bench,
		)
	}
}

func BenchmarkFmtString(b *testing.B) {
	for _, bench := range benches {
		logBenchmark(
			b,
			func(l Logger, _ int) { l.Infof("foo bar %v", "foo") },
			bench,
		)
	}
}

func BenchmarkWithField(b *testing.B) {
	f := Field("foo", "bar")

	for _, bench := range benches {
		logBenchmark(
			b,
			func(l Logger, _ int) { l.WithField(f).Info("foo bar") },
			bench,
		)
	}
}

func TestFieldThreshold(t *testing.T) {
	var (
		s recordSink

		l = NewLogger(WithSink(&s), WithDefaultFieldThreshold(record.Notice))
	)

	l.WithField(Field("foo", 1.)).Debug("foo")
	assert.Equal(t, 0, len(s.r.Fields()))

	l.WithField(Field("foo", 1)).Notice("bar")
	assert.Equal(t, []record.Field{Field("foo", 1)}, s.r.Fields())
}

func TestErrorThreshold(t *testing.T) {
	var (
		s recordSink

		err = errors.New("mock")

		l = NewLogger(WithSink(&s), WithDefaultErrorThreshold(record.Info))
	)

	l.WithError(err).Debug("foo")
	assert.Equal(t, 0, len(s.r.Errs()))

	l.WithError(err).Info("bar")
	assert.Equal(t, []error{err}, s.r.Errs())
}
