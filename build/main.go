package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"dagger.io/dagger"
)

const (
	kubectlURL = "https://dl.k8s.io/release/v1.31.0/bin/linux/amd64/kubectl"
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

func Publish(ctx context.Context, dag *dagger.Client, ctr *dagger.Container) (string, error) {
	user := cmp.Or(os.Getenv("GITHUB_ACTOR"), "github-actions")
	pass := dag.SetSecret("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	registry := "ghcr.io"
	image := fmt.Sprintf("%s/%s", registry, cmp.Or(os.Getenv("GITHUB_REPOSITORY"), "bufile"))
	return dag.Container().
		From("gcr.io/distroless/static").
		WithFile("/bufile", ctr.File("/src/bufile")).
		WithFile("/usr/bin/kubectl", dag.HTTP(kubectlURL), dagger.ContainerWithFileOpts{Permissions: 0750}).
		WithEntrypoint([]string{"/bufile"}).
		WithRegistryAuth(registry, user, pass).
		Publish(ctx, image)
}

func Base(dag *dagger.Client) *dagger.Container {
	return dag.Container().
		From("golang:1.24").
		WithEnvVariable("CGO_ENABLED", "0").
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
		fmt.Println(err)
		os.Exit(1)
	}

	ctr, err := Base(dag).
		With(Task(*task)).
		Sync(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *task != "" {
		slog.Info("[Dev mode] Image won't be published", slog.String("task", *task))
		return
	}

	ref, err := Publish(ctx, dag, ctr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	slog.Info("Image published", slog.String("ref", ref))
}
