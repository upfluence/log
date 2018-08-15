package stacktrace

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestWriteCaller(t *testing.T) {
	for _, tCase := range []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "no black list",
			out:  "stacktrace_test.go:28",
		},
		{
			name: "blacklist package",
			in:   []string{"github.com/upfluence/log/internal"},
			out:  "testing.go:777",
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			WriteCaller(buf, tCase.in)

			if res := buf.String(); tCase.out != res {
				t.Errorf("Wrong result: %v [ instead of: %v ]", res, tCase.out)
			}
		})
	}
}

func BenchmarkWriteCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WriteCaller(ioutil.Discard, []string{})
	}
}
