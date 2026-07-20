package store

import "testing"

func TestClassifyDiagnosticUsesSafeCategories(t *testing.T) {
	failure := "task_failed"
	tests := []struct {
		state, event string
		code         *string
		phase, want  string
	}{
		{"pending", "", nil, "queued", ""},
		{"running", "JOB_STARTED", nil, "starting", ""},
		{"running", "TASK_STARTED", nil, "executing", ""},
		{"failed", "HOST_UNREACHABLE", nil, "complete", "host_unreachable"},
		{"failed", "HOST_FAILED", &failure, "complete", "task_failed"},
	}
	for _, test := range tests {
		phase, code := classifyDiagnostic(test.state, test.event, test.code)
		if phase != test.phase || code != test.want {
			t.Fatalf("classify(%q,%q)=(%q,%q), want (%q,%q)", test.state, test.event, phase, code, test.phase, test.want)
		}
	}
}
