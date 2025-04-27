package main

import (
	"context"
	"github.com/Ippolid/auth/internal/app"
	"log"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
