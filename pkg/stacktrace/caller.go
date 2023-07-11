package stacktrace

import (
	"runtime"

	"github.com/upfluence/log/internal/stacktrace"
)

func FindCaller(depth int, blacklist []string) *runtime.Frame {
	return stacktrace.FindCaller(
		1+depth,
		append(blacklist, "github.com/upfluence/log"),
	)
}
