package main

import (
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Core = (*core)(nil)

type core struct {
	level    zapcore.Level
	delegate zapcore.Core
}

func NewZapCore(name string, provider log.LoggerProvider) zapcore.Core {
	return &core{
		level:    zapcore.InfoLevel,
		delegate: otelzap.NewCore(name, otelzap.WithLoggerProvider(provider)),
	}
}

func (c *core) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level) && c.delegate.Enabled(level)
}

func (c *core) With(fields []zapcore.Field) zapcore.Core {
	return &core{
		level:    c.level,
		delegate: c.delegate.With(fields),
	}
}

func (c *core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		ce.AddCore(ent, c)
	}
	return ce
}

func (c *core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// override to avoid this lock on every write
	// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/bridges/otelzap/core.go#L235
	// https://github.com/open-telemetry/opentelemetry-go/blob/main/sdk/log/provider.go#L124
	entry.LoggerName = ""
	return c.delegate.Write(entry, fields)
}

func (c *core) Sync() error {
	return c.delegate.Sync()
}
