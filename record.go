package log

import "github.com/upfluence/log/record"

var SkipFrame = nullField{}

type nullField struct{}

func (nullField) GetKey() string   { return "" }
func (nullField) GetValue() string { return "" }
func (nullField) Discard() bool    { return true }

type withFields struct {
	record.Context

	fields []record.Field
}

func (wf *withFields) Fields() []record.Field {
	return append(wf.Context.Fields(), wf.fields...)
}

type withErrors struct {
	record.Context

	errs []error
}

func (we *withErrors) Errs() []error {
	return append(we.Context.Errs(), we.errs...)
}

type nullContext struct{}

func (*nullContext) Fields() []record.Field { return nil }
func (*nullContext) Errs() []error          { return nil }
