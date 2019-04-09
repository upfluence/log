package writer

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/upfluence/log"
	"github.com/upfluence/log/record"
)

type mockRecord struct {
	msg string

	t0  time.Time
	lvl record.Level

	fs   []record.Field
	errs []error

	args []interface{}
}

func (r mockRecord) Fields() []record.Field     { return r.fs }
func (r mockRecord) Errs() []error              { return r.errs }
func (mockRecord) ID() uint64                   { return 1 }
func (r mockRecord) Time() time.Time            { return r.t0 }
func (r mockRecord) Level() record.Level        { return r.lvl }
func (r mockRecord) WriteFormatted(w io.Writer) { io.WriteString(w, r.msg) }
func (r mockRecord) Args() []interface{}        { return r.args }

func TestLog(t *testing.T) {
	for _, tt := range []struct {
		f   Formatter
		r   record.Record
		err error
		out string
	}{
		{
			f:   NewFastFormatter(),
			r:   mockRecord{msg: "foo bar", lvl: record.Info},
			out: "[I 010101 00:00:00] foo bar\n",
		},
		{
			f: NewFastFormatter(),
			r: mockRecord{
				msg:  "foo bar",
				errs: []error{errors.New("foo bar")},
				lvl:  record.Info,
			},
			out: "[I 010101 00:00:00] foo bar [error: foo bar]\n",
		},
		{
			f: NewFastFormatter(),
			r: mockRecord{
				msg: "foo bar",
				fs:  []record.Field{log.Field("fiz", "buz")},
				lvl: record.Info,
			},
			out: "[I 010101 00:00:00] [fiz: buz] foo bar\n",
		},
	} {
		var buf bytes.Buffer

		if err := NewSink(tt.f, &buf).Log(tt.r); err != tt.err {
			t.Errorf("Log() = %v [ expected: %v ]", err, tt.err)
		}

		if res := buf.String(); res != tt.out {
			t.Errorf("Wrote: %q, expected: %q", res, tt.out)
		}
	}
}
