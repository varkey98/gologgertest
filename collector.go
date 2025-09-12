package main

import (
	"os"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/exporter/nopexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initialize(logger *zap.Logger) (otelcol.CollectorSettings, error) {

	cfgProviderSettings := otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs: []string{getConfigFlag(os.Args[1:])},
			ProviderFactories: []confmap.ProviderFactory{
				fileprovider.NewFactory(),
				envprovider.NewFactory(),
			},
			ConverterSettings: confmap.ConverterSettings{Logger: logger},
		},
	}

	settings := otelcol.CollectorSettings{
		BuildInfo: component.BuildInfo{
			Command:     "test-collector",
			Description: "Test Collector",
			Version:     "1.0.0",
		},
		Factories:              getFactories,
		ConfigProviderSettings: cfgProviderSettings,
		LoggingOptions: []zap.Option{
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return logger.Core()
			}),
		},
	}

	return settings, nil
}

func getFactories() (otelcol.Factories, error) {
	receivers, err := otelcol.MakeFactoryMap[receiver.Factory](otlpreceiver.NewFactory())
	if err != nil {
		return otelcol.Factories{}, err
	}

	exporters, err := otelcol.MakeFactoryMap[exporter.Factory](
		debugexporter.NewFactory(),
		nopexporter.NewFactory())
	if err != nil {
		return otelcol.Factories{}, err
	}

	return otelcol.Factories{
		Receivers: receivers,
		Exporters: exporters,
	}, nil
}

func getConfigFlag(args []string) string {
	for i, v := range args {
		if v == "--config" {
			if len(args) > i+1 {
				return args[i+1]
			} else {
				return ""
			}
		} else if strings.HasPrefix(v, "--config") {
			return v[len("--config="):]
		}
	}

	return ""
}
