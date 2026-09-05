package middleware

import "testing"

func TestShouldNotifyDatabaseActivity(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "/api/status", want: false},
		{path: "/", want: false},
		{path: "/assets/index.js", want: false},
		{path: "/api/status/test", want: true},
		{path: "/api/log/", want: true},
		{path: "/v1/chat/completions", want: true},
		{path: "/v1beta/models/gemini:generateContent", want: true},
		{path: "/pg/chat/completions", want: true},
		{path: "/mj/task/1/fetch", want: true},
		{path: "/suno/mj/task/1/fetch", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := shouldNotifyDatabaseActivity(tt.path); got != tt.want {
				t.Fatalf("shouldNotifyDatabaseActivity(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
