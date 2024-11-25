package route

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	if os.Getenv("BUF_TOKEN") == "" {
		t.Skip("skipping test; BUF_TOKEN not set")
	}
	var tests = []struct {
		name  string
		given Module
	}{
		{"correct mod", Module{URL: "buf.build/vic3lord/bufile"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var buf strings.Builder
			err := Generate(ctx, tt.given, &buf, Options{})
			if err != nil {
				t.Errorf("expected nil, actual %v", err)
			}
		})
	}
}
