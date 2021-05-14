package error_logger

import "github.com/upfluence/errors/reporter"

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
