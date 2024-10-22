package config

import (
	"slices"
	"testing"

	"github.com/vic3lord/bufile/route"
)

func TestParse(t *testing.T) {
	cfg, err := Parse("testdata/config.json")
	if err != nil {
		t.Error(err)
	}

	ok := slices.ContainsFunc(cfg.Modules, func(m route.Module) bool {
		return m.URL == "buf.build/vic3lord/bufile"
	})
	if !ok {
		t.Error("expected module not found")
	}
}
