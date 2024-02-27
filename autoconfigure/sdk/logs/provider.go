package logs

import (
	"context"
	"errors"
	"github.com/kudarap/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/kudarap/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	"github.com/kudarap/opentelemetry-logs-go/internal/global"
	sdk "github.com/kudarap/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
)

// loggerProviderConfig Configuration for Logger Provider
type loggerProviderConfig struct {
	processors []sdk.LogRecordProcessor
	// resource contains attributes representing an entity that produces telemetry.
	resource *resource.Resource
}

// LoggerProviderOption configures a LoggerProvider.
type LoggerProviderOption interface {
	apply(loggerProviderConfig) loggerProviderConfig
}
type loggerProviderOptionFunc func(loggerProviderConfig) loggerProviderConfig

func (fn loggerProviderOptionFunc) apply(cfg loggerProviderConfig) loggerProviderConfig {
	return fn(cfg)
}

// WithLogRecordProcessors will configure processor to process logs
func WithLogRecordProcessors(logsProcessors []sdk.LogRecordProcessor) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg loggerProviderConfig) loggerProviderConfig {
		cfg.processors = logsProcessors
		return cfg
	})
}

// WithResource will configure OTLP logger with common resource attributes.
//
// Parameters:
// r (*resource.Resource) list of resources will be added to every log as resource level tags
func WithResource(r *resource.Resource) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg loggerProviderConfig) loggerProviderConfig {
		var err error
		cfg.resource, err = resource.Merge(resource.Environment(), r)
		if err != nil {
			otel.Handle(err)
		}
		return cfg
	})
}

func applyLoggerProviderExporterEnvConfigs(ctx context.Context, cfg loggerProviderConfig) loggerProviderConfig {

	// if processors already defined explicitly - skip env configuration
	if cfg.processors != nil {
		return cfg
	}

	exporters, isProvided := exportersFromEnv()
	if isProvided == false || len(exporters) == 0 {
		exporters = []string{logsExporterOTLP}
	}

	// currently values are hardcoded, but subject to be extracted into global map to LoggerProviderOption
	// to support custom exporters
	for _, exporter := range exporters {
		switch exporter {
		case logsExporterNone:
		case logsExporterOTLP:
			otlpExporter, err := otlplogs.NewExporter(ctx)
			if err != nil {
				global.Error(err, "Can't instantiate otlp exporter")
			}
			cfg.processors = append(cfg.processors, sdk.NewBatchLogRecordProcessor(otlpExporter))
		case logsExporterLogging:
			sdtoutExporter, err := stdoutlogs.NewExporter()
			if err != nil {
				global.Error(err, "Can't instantiate logging exporter")
			}
			cfg.processors = append(cfg.processors, sdk.NewSimpleLogRecordProcessor(sdtoutExporter))
		default:
			err := errors.New("Exporter " + exporter + " is not supported")
			if err != nil {
				global.Error(err, "Can't instantiate "+exporter+" exporter")
			}
		}
	}
	return cfg
}

// NewLoggerProvider will autoconfigure exporters and create logger provider
func NewLoggerProvider(ctx context.Context, opts ...LoggerProviderOption) *sdk.LoggerProvider {

	o := loggerProviderConfig{}

	for _, opt := range opts {
		o = opt.apply(o)
	}

	// apply exporter env options after as should not instantiate exporters if they will be overridden
	o = applyLoggerProviderExporterEnvConfigs(ctx, o)

	var sdkOptions []sdk.LoggerProviderOption

	for _, processor := range o.processors {
		sdkOptions = append(sdkOptions, sdk.WithLogRecordProcessor(processor))
	}
	sdkOptions = append(sdkOptions, sdk.WithResource(o.resource))

	return sdk.NewLoggerProvider(sdkOptions...)
}
