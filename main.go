package main

import (
	"context"
	"log"
	"os"

	"github.com/vic3lord/bufile/route"
)

func main() {
	err := route.Generate(context.Background(), os.Stdout)
	if err != nil {
		log.Fatalf("failed to generate routes: %v", err)
	}
}
