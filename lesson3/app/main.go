package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	tracer, closer := initJaeger("example", logger)
	defer closer.Close()

	a := app{
		logger: logger,
		tracer: tracer,
	}

	a.logger.Debug("starting server ...")
	if err := a.Serve(); err != nil {
		a.logger.With(zap.Error(err)).Fatal("couldnt start server")
	}
}
