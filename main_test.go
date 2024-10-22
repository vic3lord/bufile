package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/vic3lord/bufile/route"
)

func TestRunMain(t *testing.T) {
	if os.Getenv("BUF_TOKEN") == "" {
		t.Skip("skipping test; BUF_TOKEN not set")
	}
	var tests = []struct {
		name  string
		given route.Module
	}{
		{"correct mod", route.Module{URL: "buf.build/vic3lord/bufile"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var buf strings.Builder
			given := []route.Module{tt.given}
			err := run(ctx, given, &buf)
			if err != nil {
				t.Errorf("expected nil, actual %v", err)
			}
		})
	}
}
