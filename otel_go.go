package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	x = 5
)

func initOtelGo() func() {

	logsOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint("localhost:5317"),
		otlploggrpc.WithInsecure(),
	}

	logsExporter, err := otlploggrpc.New(context.Background(), logsOpts...)
	if err != nil {
		log.Fatal(err)
	}

	var resourceAttrs []attribute.KeyValue
	logsBatchProcessor := sdklog.NewBatchProcessor(logsExporter)
	resourceAttrs = append(resourceAttrs, attribute.String("service.instance.id", uuid.NewString()))
	logsResource, err := resource.New(context.Background(), resource.WithAttributes(resourceAttrs...))
	if err != nil {
		log.Fatal(err)
	}

	loggerProvider := sdklog.NewLoggerProvider(sdklog.WithResource(logsResource), sdklog.WithProcessor(logsBatchProcessor))
	global.SetLoggerProvider(loggerProvider)
	return func() {
		err = loggerProvider.Shutdown(context.Background())
		if err != nil {
			log.Printf("an error while calling metrics provider shutdown: %v", err)
		}

		err := logsBatchProcessor.Shutdown(context.Background())
		if err != nil {
			log.Printf("an error while calling metrics reader shutdown: %v", err)
		}
	}
}
