package writer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/upfluence/log"
	"github.com/upfluence/log/logtest"
	"github.com/upfluence/log/record"
)

func TestLog(t *testing.T) {
	for _, tt := range []struct {
		f    Formatter
		opts []logtest.RecordOption
		err  error
		out  string
	}{
		{
			f: NewFastFormatter(),
			opts: []logtest.RecordOption{
				logtest.WithMessage("foo bar"),
				logtest.WithLevel(record.Info),
			},
			out: "[I 010101 00:00:00] foo bar\n",
		},
		{
			f: NewFastFormatter(),
			opts: []logtest.RecordOption{
				logtest.WithMessage("foo bar"),
				logtest.WithLevel(record.Info),
				logtest.WithErrors(errors.New("foo bar")),
			},
			out: "[I 010101 00:00:00] foo bar [error: foo bar]\n",
		},
		{
			f: NewFastFormatter(),
			opts: []logtest.RecordOption{
				logtest.WithMessage("foo bar"),
				logtest.WithLevel(record.Info),
				logtest.WithFields(log.Field("fiz", "buz")),
			},
			out: "[I 010101 00:00:00] [fiz: buz] foo bar\n",
		},
	} {
		var buf bytes.Buffer

		r := logtest.BuildRecord(tt.opts...)

		if err := NewSink(tt.f, &buf).Log(r); err != tt.err {
			t.Errorf("Log() = %v [ expected: %v ]", err, tt.err)
		}

		if res := buf.String(); res != tt.out {
			t.Errorf("Wrote: %q, expected: %q", res, tt.out)
		}
	}
}
