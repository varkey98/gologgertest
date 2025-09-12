package main

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	shutdownFn := initOtelGo()
	defer shutdownFn()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, otelzap.NewCore("test-collector"))
	}))

	settings, err := initialize(logger)
	if err != nil {
		logger.Fatal("failed to initialize otel settings", zap.Error(err))
	}

	app, err := otelcol.NewCollector(settings)
	if err != nil {
		logger.Fatal("failed to create otel collector", zap.Error(err))
	}

	err = app.Run(context.Background())
	if err != nil {
		logger.Fatal("failed to start otel collector", zap.Error(err))
	}
}
