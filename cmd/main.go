package main

import (
	"context"

	"github.com/valek177/auth/internal/app"
	"github.com/valek177/auth/internal/logger"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.FatalWithMsg("failed to init app: ", err)
	}

	err = a.Run(ctx)
	if err != nil {
		logger.FatalWithMsg("failed to run app: ", err)
	}
}
