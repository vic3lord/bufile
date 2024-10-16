package config

import (
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	cfg, err := Parse("testdata/config.json")
	if err != nil {
		t.Error(err)
	}

	ok := slices.Contains(cfg.Modules, "buf.build/vic3lord/bufile")
	if !ok {
		t.Error("expected module not found")
	}
}
