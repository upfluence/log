package writer

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/upfluence/log"
	"github.com/upfluence/log/record"
)

const (
	stdFmt  = "[%s %s %s:%d] %s%s%s"
	dateFmt = "020106 15:04:05"
)

var levelPrettifier = map[record.Level]string{
	record.Debug:   "D",
	record.Info:    "I",
	record.Notice:  "N",
	record.Warning: "W",
	record.Error:   "E",
	record.Fatal:   "F",
}

type formatter struct {
	calldepth int
}

type discarder interface {
	Discard() bool
}

func (f *formatter) formatFields(fs []record.Field) string {
	if len(fs) == 0 {
		return ""
	}

	var res string

	for _, f := range fs {
		d, ok := f.(discarder)

		if ok && d.Discard() {
			continue
		}

		res += fmt.Sprintf("[%s: %s]", f.GetKey(), f.GetValue())
	}

	return res + " "
}

func (f *formatter) formatErrs(errs []error) string {
	if len(errs) == 0 {
		return ""
	}

	var res string

	for _, err := range errs {
		res += fmt.Sprintf("[error: %v]", err)
	}

	return " " + res
}

func (f *formatter) Format(r record.Record) string {
	var depth = f.calldepth + 1

	for _, f := range r.Fields() {
		if f == log.SkipFrame {
			depth++
		}
	}

	_, file, line, ok := runtime.Caller(depth)

	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	return fmt.Sprintf(
		stdFmt,
		levelPrettifier[r.Level()],
		r.Time().Format(dateFmt),
		file,
		line,
		f.formatFields(r.Fields()),
		r.Formatted(),
		f.formatErrs(r.Errs()),
	)
}
