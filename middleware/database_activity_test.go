package middleware

import "testing"

func TestShouldNotifyDatabaseActivity(t *testing.T) {
	tests := []struct {
		path     string
		routeTag string
		want     bool
	}{
		{path: "/api/status", routeTag: "api", want: false},
		{path: "/", want: false},
		{path: "/assets/index.js", want: false},
		{path: "/api/status/test", routeTag: "api", want: true},
		{path: "/api/log/", routeTag: "api", want: true},
		{path: "/v1/chat/completions", routeTag: "relay", want: true},
		{path: "/v1beta/models/gemini:generateContent", routeTag: "relay", want: true},
		{path: "/pg/chat/completions", routeTag: "relay", want: true},
		{path: "/mj/task/1/fetch", routeTag: "relay", want: true},
		{path: "/suno/mj/task/1/fetch", routeTag: "relay", want: true},
		{path: "/kling/task/1/fetch", routeTag: "relay", want: true},
		{path: "/vendor/jobs", routeTag: "relay", want: true},
		{path: "/legacy/dashboard", routeTag: "old_api", want: true},
		{path: "/api/not-found", want: false},
		{path: "/v1/not-found", want: false},
		{path: "/api/status/", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := shouldNotifyDatabaseActivity(tt.path, tt.routeTag); got != tt.want {
				t.Fatalf("shouldNotifyDatabaseActivity(%q, %q) = %v, want %v", tt.path, tt.routeTag, got, tt.want)
			}
		})
	}
}
