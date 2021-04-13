package logtest

import (
	"io"
	"time"

	"github.com/upfluence/log/record"
)

type RecordOption func(*mockRecord)

type mockRecord struct {
	msg string
	id  uint64

	t0  time.Time
	lvl record.Level

	fs   []record.Field
	errs []error

	args []interface{}
}

func (r mockRecord) Fields() []record.Field     { return r.fs }
func (r mockRecord) Errs() []error              { return r.errs }
func (r mockRecord) ID() uint64                 { return r.id }
func (r mockRecord) Time() time.Time            { return r.t0 }
func (r mockRecord) Level() record.Level        { return r.lvl }
func (r mockRecord) WriteFormatted(w io.Writer) { io.WriteString(w, r.msg) }
func (r mockRecord) Args() []interface{}        { return r.args }

var defaultRecord = mockRecord{msg: "default msg", id: 1, lvl: record.Error}

func WithLevel(lvl record.Level) RecordOption {
	return func(r *mockRecord) { r.lvl = lvl }
}

func WithMessage(msg string) RecordOption {
	return func(r *mockRecord) { r.msg = msg }
}

func WithArgs(args ...interface{}) RecordOption {
	return func(r *mockRecord) { r.args = args }
}

func WithErrors(errs ...error) RecordOption {
	return func(r *mockRecord) { r.errs = append(r.errs, errs...) }
}

func WithFields(fs ...record.Field) RecordOption {
	return func(r *mockRecord) { r.fs = append(r.fs, fs...) }
}

func BuildRecord(opts ...RecordOption) record.Record {
	var r = defaultRecord

	for _, opt := range opts {
		opt(&r)
	}

	return r
}
