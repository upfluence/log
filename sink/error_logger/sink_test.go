package error_logger

import (
	"errors"
	"reflect"
	"testing"

	"github.com/upfluence/errors/reporter"
	"github.com/upfluence/log"
	"github.com/upfluence/log/logtest"
)

type report struct {
	err  error
	opts reporter.ReportOptions
}

type mockReporter struct {
	rs []report
}

func (mr *mockReporter) Close() error { return nil }

func (mr *mockReporter) Report(err error, opts reporter.ReportOptions) {
	mr.rs = append(mr.rs, report{err: err, opts: opts})
}

func TestSink(t *testing.T) {
	for _, tt := range []struct {
		opts []logtest.RecordOption

		reports []report
	}{
		{
			reports: []report{
				{
					err:  errors.New("default msg"),
					opts: reporter.ReportOptions{Depth: 2, Tags: map[string]interface{}{}},
				},
			},
		},
		{
			opts: []logtest.RecordOption{logtest.WithArgs(errors.New("err1"))},
			reports: []report{
				{
					err:  errors.New("err1"),
					opts: reporter.ReportOptions{Depth: 2, Tags: map[string]interface{}{}},
				},
			},
		},
		{
			opts: []logtest.RecordOption{
				logtest.WithErrors(errors.New("err1")),
				logtest.WithFields(log.Field("foo", "bar")),
			},
			reports: []report{
				{
					err: errors.New("err1"),
					opts: reporter.ReportOptions{
						Depth: 2,
						Tags:  map[string]interface{}{"foo": "bar"},
					},
				},
			},
		},
		{
			opts: []logtest.RecordOption{
				logtest.WithErrors(errors.New("err1"), errors.New("err2")),
			},
			reports: []report{
				{
					err:  errors.New("err1"),
					opts: reporter.ReportOptions{Depth: 2, Tags: map[string]interface{}{}},
				},
				{
					err:  errors.New("err2"),
					opts: reporter.ReportOptions{Depth: 2, Tags: map[string]interface{}{}},
				},
			},
		},
	} {
		var mr mockReporter

		err := WrapReporter(&mr, 1).Log(logtest.BuildRecord(tt.opts...))

		if err != nil {
			t.Errorf("Unexpected error on Log(): %+v", err)
		}

		if !reflect.DeepEqual(mr.rs, tt.reports) {
			t.Errorf("Unexpected reports: %+v [want: %+v]", mr.rs, tt.reports)
		}
	}
}
