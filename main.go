package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/vic3lord/bufile/config"
	"github.com/vic3lord/bufile/route"
)

var (
	configFile = flag.String("config", "bufile.json", "Path to config file")
	apply      = flag.Bool("apply", false, "Apply the generated routes")
)

func run(ctx context.Context, cfg config.Config, w io.Writer) error {
	var errs error
	opts := route.Options{
		IncludeServiceName: cfg.IncludeServiceName,
	}
	for _, mod := range cfg.Modules {
		err := route.Generate(ctx, mod, w, opts)
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

	var sb strings.Builder
	ctx := context.Background()
	if err := run(ctx, cfg, &sb); err != nil {
		slog.Error("Generate routes", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if *apply {
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			slog.Error("Apply routes", slog.String("output", string(out)), slog.String("err", err.Error()))
			os.Exit(1)
		}
		fmt.Print(string(out))
		os.Exit(0)
	}
	fmt.Print(sb.String())
}
