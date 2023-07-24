package logs

import "go.opentelemetry.io/otel/attribute"

// LoggerConfig is a group of options for a Logger.
type LoggerConfig struct {
	instrumentationVersion string
	// Schema URL of the telemetry emitted by the Tracer.
	schemaURL string
	attrs     attribute.Set
}

// InstrumentationVersion returns the version of the library providing instrumentation.
func (t *LoggerConfig) InstrumentationVersion() string {
	return t.instrumentationVersion
}

// InstrumentationAttributes returns the attributes associated with the library
// providing instrumentation.
func (t *LoggerConfig) InstrumentationAttributes() attribute.Set {
	return t.attrs
}

// SchemaURL returns the Schema URL of the telemetry emitted by the Tracer.
func (t *LoggerConfig) SchemaURL() string {
	return t.schemaURL
}

// NewLoggerConfig applies all the options to a returned LoggerConfig.
func NewLoggerConfig(options ...LoggerOption) LoggerConfig {
	var config LoggerConfig
	for _, option := range options {
		config = option.apply(config)
	}
	return config
}

// LoggerOption applies an option to a LoggerConfig.
type LoggerOption interface {
	apply(LoggerConfig) LoggerConfig
}

type loggerOptionFunc func(LoggerConfig) LoggerConfig

func (fn loggerOptionFunc) apply(cfg LoggerConfig) LoggerConfig {
	return fn(cfg)
}
