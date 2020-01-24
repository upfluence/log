package stacktrace

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestWriteCaller(t *testing.T) {
	for _, tt := range []struct {
		name     string
		in       []string
		callerfn func(*testing.T, string)
	}{
		{
			name: "no black list",
			callerfn: func(t *testing.T, c string) {
				if c != "stacktrace_test.go:37" {
					t.Errorf("invalid caller: %q", c)
				}
			},
		},
		{
			name: "blacklist package",
			in:   []string{"github.com/upfluence/log/internal"},
			callerfn: func(t *testing.T, c string) {
				if cs := strings.Split(c, ":"); len(cs) != 2 || cs[0] != "testing.go" {
					t.Errorf("invalid caller: %q", c)
				}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			WriteCaller(&buf, tt.in)

			tt.callerfn(t, buf.String())
		})
	}
}

func BenchmarkWriteCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WriteCaller(ioutil.Discard, []string{})
	}
}
