package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/vic3lord/bufile/config"
	"github.com/vic3lord/bufile/route"
)

var (
	configFile = flag.String("config", "bufile.json", "Path to config file")
)

func main() {
	flag.Parse()
	cfg, err := config.Parse(*configFile)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	ctx := context.Background()
	for _, mod := range cfg.Modules {
		err = route.Generate(ctx, mod, os.Stdout)
		if err != nil {
			l := slog.With(
				slog.String("module", mod),
				slog.String("err", err.Error()),
			)
			l.Error("failed to generate routes")
		}
	}
}
