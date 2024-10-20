package main

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestRunMain(t *testing.T) {
	if os.Getenv("BUF_TOKEN") == "" {
		t.Skip("skipping test; BUF_TOKEN not set")
	}
	var tests = []struct {
		name  string
		given []string
	}{
		{"correct mod", []string{"buf.build/vic3lord/bufile"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var buf strings.Builder
			err := run(ctx, tt.given, &buf)
			if err != nil {
				t.Errorf("expected nil, actual %v", err)
			}
		})
	}
}
