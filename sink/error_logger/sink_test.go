package error_logger

import (
	"testing"

	"github.com/upfluence/log/record"
)

type mockLogger struct {
	calls []struct {
		err  error
		tags map[string]interface{}
	}
}

func TestSink_Log(t *testing.T) {
	tests := []struct {
		name    string
		eLogger ErrorLogger
		record  record
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sink{
				eLogger: tt.fields.eLogger,
			}
			if err := s.Log(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Sink.Log() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
