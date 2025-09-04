package main

import (
	"context"
	"log"

	"github.com/kiryu2k/onlyfriends-auth-service/config"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/app"
	"github.com/pkg/errors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "load config"))
	}
	ctx := context.Background()
	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(errors.WithMessage(err, "run app"))
	}
}
