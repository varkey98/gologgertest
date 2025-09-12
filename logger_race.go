package main

import (
	"github.com/Traceableai/goagent"
	"github.com/Traceableai/goagent/config"
	v1 "github.com/hypertrace/agent-config/gen/go/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	x = 5
)

func testLoggerRace() {
	cfg := config.Load()
	//cfg.Tracing.Reporting.TraceReporterType = v1.TraceReporterType_LOGGING
	cfg.Tracing.Telemetry = &v1.Telemetry{
		Logs: &v1.LogsExport{
			Enabled: v1.Bool(true),
			Level:   v1.LogLevel_LOG_LEVEL_INFO,
		},
	}
	goagent.Init(cfg)
	core2 := goagent.NewZapCore("test", &v1.LogsExport{
		Enabled: v1.Bool(true),
		Level:   v1.LogLevel_LOG_LEVEL_INFO,
	})
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
