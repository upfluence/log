package error_logger

import (
	"bytes"
	"errors"

	"github.com/upfluence/errors/reporter"
	"github.com/upfluence/log/record"
	"github.com/upfluence/log/sink"
)

type ErrorLogger interface {
	Capture(error, map[string]interface{}) error
}

type errorLoggerWrapper struct {
	el ErrorLogger
}

func (elw *errorLoggerWrapper) Report(err error, opts reporter.ReportOptions) {
	elw.el.Capture(err, opts.Tags)
}

func (elw *errorLoggerWrapper) Close() error { return nil }

type Sink struct {
	r reporter.Reporter
	d int
}

func NewSink(el ErrorLogger) sink.Sink {
	return WrapReporter(&errorLoggerWrapper{el: el}, 0)
}

func WrapReporter(r reporter.Reporter, depth int) sink.Sink {
	return &Sink{r: r, d: depth + 1}
}

func (s *Sink) Log(r record.Record) error {
	var (
		errs = r.Errs()
		tags = map[string]interface{}{}
	)

	if len(errs) == 0 {
		for _, arg := range r.Args() {
			if err, ok := arg.(error); ok {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) == 0 {
		var buf bytes.Buffer
		r.WriteFormatted(&buf)

		errs = []error{errors.New(buf.String())}
	}

	for _, f := range r.Fields() {
		tags[f.GetKey()] = f.GetValue()
	}

	for _, err := range errs {
		s.r.Report(err, reporter.ReportOptions{Tags: tags, Depth: s.d})
	}

	return nil
}
