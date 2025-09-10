package main

import (
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	x = 5
)

func testLoggerRace() {
	initializeLogs()

	loggerProvider := global.GetLoggerProvider()

	core2 := NewZapCore("test", loggerProvider)
	logger, _ := zap.NewProduction()

	logger = zap.New(zapcore.NewTee(logger.Core(), core2))
	for range 10 {
		go func() {
			//record := log.Record{}
			//record.SetBody(log.StringValue("test"))
			//logger.Emit(context.Background(), record)
			logger.Info("test")
		}()
	}

	time.Sleep(30 * time.Second)
}
