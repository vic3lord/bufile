package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/vic3lord/bufile/config"
	"github.com/vic3lord/bufile/route"
)

var (
	configFile = flag.String("config", "bufile.json", "Path to config file")
)

func run(ctx context.Context, mods []route.Module, w io.Writer) error {
	var errs error
	for _, mod := range mods {
		err := route.Generate(ctx, mod, w)
		if err != nil {
			moderr := fmt.Errorf("generate routes for module %q: %w", mod, err)
			errs = errors.Join(errs, moderr)
		}
	}
	return errs

}

func main() {
	flag.Parse()
	cfg, err := config.Parse(*configFile)
	if err != nil {
		slog.Error("Parse config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()
	if err := run(ctx, cfg.Modules, os.Stdout); err != nil {
		slog.Error("Generate routes", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
