package logtest

import (
	"context"
	"testing"

	"github.com/upfluence/log"
	"github.com/upfluence/log/record"
	"github.com/upfluence/log/sink/writer"
)

func WrapTestingLogger(t testing.TB) log.Logger {
	return &testingLogger{
		t: t,
		Logger: log.NewLogger(
			log.WithSink(
				writer.NewSink(writer.NewFastFormatter(), &testingWriter{t: t}),
			),
		),
	}
}

type testingWriter struct {
	t testing.TB
}

func (tw *testingWriter) Write(p []byte) (int, error) {
	tw.t.Log(string(p))
	return len(p), nil
}

type testingLogger struct {
	log.Logger
	t testing.TB
}

func (tl *testingLogger) WithField(f record.Field) log.SugaredLogger {
	return &testingLogger{Logger: tl.Logger.WithField(f), t: tl.t}
}

func (tl *testingLogger) WithFields(fs ...record.Field) log.SugaredLogger {
	return &testingLogger{Logger: tl.Logger.WithFields(fs...), t: tl.t}
}
func (tl *testingLogger) WithContext(ctx context.Context) log.SugaredLogger {
	return &testingLogger{Logger: tl.Logger.WithContext(ctx), t: tl.t}
}

func (tl *testingLogger) WithError(err error) log.SugaredLogger {
	return &testingLogger{Logger: tl.Logger.WithError(err), t: tl.t}
}

func (tl *testingLogger) Fatal(vs ...interface{}) { tl.t.Fatal(vs...) }

func (tl *testingLogger) Fatalf(fmt string, vs ...interface{}) {
	tl.t.Fatalf(fmt, vs...)
}
