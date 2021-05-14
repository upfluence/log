package stacktrace

import (
	"io"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const size = 64

var (
	stacktracePool = &sync.Pool{
		New: func() interface{} {
			return make([]uintptr, size)
		},
	}

	semicolon     = []byte(":")
	defaultCaller = []byte("???:0")
)

func WriteCaller(w io.Writer, blacklist []string) {
	pcs := stacktracePool.Get().([]uintptr)
	defer stacktracePool.Put(pcs)

	runtime.Callers(2, pcs)

	frames := runtime.CallersFrames(pcs)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if shouldSkip(&frame, blacklist) {
			continue
		}

		io.WriteString(w, filepath.Base(frame.File))
		w.Write(semicolon)
		io.WriteString(w, strconv.Itoa(frame.Line))
		return
	}

	w.Write(defaultCaller)
}

func FrameDepth(blacklist []string) int {
	var i = 1

	if len(blacklist) == 0 {
		return i
	}

	pcs := stacktracePool.Get().([]uintptr)
	defer stacktracePool.Put(pcs)

	runtime.Callers(2, pcs)

	frames := runtime.CallersFrames(pcs)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if !shouldSkip(&frame, blacklist) {
			break
		}

		i++
	}

	return i
}

func shouldSkip(f *runtime.Frame, paths []string) bool {
	for _, p := range paths {
		if strings.Contains(f.Function, p) {
			return true
		}
	}

	return false
}
