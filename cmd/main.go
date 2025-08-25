package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"shot/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	application := app.New()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("fauiled to run application err: %s", err.Error())
	}
}
