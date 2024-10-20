package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"dagger.io/dagger"
)

var task = flag.String("task", "", "task to run")

func Test(ctr *dagger.Container) *dagger.Container {
	return ctr.WithExec([]string{"go", "test", "-v", "./..."})
}

func Build(ctr *dagger.Container) *dagger.Container {
	return ctr.WithExec([]string{"go", "build", "."})
}

func All(ctr *dagger.Container) *dagger.Container {
	return ctr.
		With(Test).
		With(Build)
}

func Task(name string) func(ctr *dagger.Container) *dagger.Container {
	switch name {
	case "test":
		return Test
	case "build":
		return Build
	}
	return All
}

func Publish(dag *dagger.Client, ctr *dagger.Container) *dagger.Container {
	return dag.Container().
		From("gcr.io/distroless/static").
		WithFile("/bufile", ctr.File("/src/bufile")).
		WithEntrypoint([]string{"/bufile"})
}

func Base(dag *dagger.Client) *dagger.Container {
	return dag.Container().
		From("golang:1.23").
		WithMountedDirectory("/src", dag.Host().Directory(".")).
		WithWorkdir("/src").
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("/root/.cache/go-build")).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("/go/pkg/mod"))
}

func main() {
	flag.Parse()
	ctx := context.Background()
	dag, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		slog.Error("Connect dagger engine", slog.String("err", err.Error()))
		os.Exit(1)
	}

	ctr, err := Base(dag).
		With(Task(*task)).
		Sync(ctx)
	if err != nil {
		slog.Error("Sync container", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if *task != "" {
		slog.Info("[Dev mode] Image won't be published", slog.String("task", *task))
		return
	}

	_, err = Publish(dag, ctr).Sync(ctx)
	if err != nil {
		slog.Error("Publish container", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
