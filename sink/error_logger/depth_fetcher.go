package error_logger

import "github.com/upfluence/log/internal/stacktrace"

var defaultBlacklist = []string{"github.com/upfluence/log"}

type depthFetcher interface {
	fetch() int
}

type staticDepthFetcher int

func (sdf staticDepthFetcher) fetch() int { return int(sdf) }

type blacklistDepthFetcher []string

func (bdf blacklistDepthFetcher) fetch() int {
	return stacktrace.FrameDepth(bdf) - 1
}
