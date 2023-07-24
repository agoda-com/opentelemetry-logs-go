package logs // Package logs import "github.com/agoda-com/opentelemetry-logs-go/otel/sdk/logs"
import (
	"context"
	"fmt"
	"github.com/agoda-com/opentelemetry-logs-go/otel/internal/global"
	"github.com/agoda-com/opentelemetry-logs-go/otel/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"sync"
	"sync/atomic"
)

const (
	defaultLoggerName = "gitlab.agodadev.io/devenv/devstack/pkg/otel/sdk/logger"
)

// loggerProviderConfig Configuration for OTEL extension for zap logger
type loggerProviderConfig struct {
	processors []LogRecordProcessor
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

// WithLogsProcessor will configure processor to process logs
func WithLogsProcessor(logsProcessor LogRecordProcessor) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg loggerProviderConfig) loggerProviderConfig {
		cfg.processors = append(cfg.processors, logsProcessor)
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

// LoggerProvider provide access to Logger. The API is not intended to be called by application developers directly.
// see https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#loggerprovider
type LoggerProvider struct {
	mu          sync.Mutex
	namedLogger map[instrumentation.Scope]*logger
	//cfg loggerProviderConfig

	logProcessors atomic.Pointer[logRecordProcessorStates]
	isShutdown    atomic.Bool

	// These fields are not protected by the lock mu. They are assumed to be
	// immutable after creation of the LoggerProvider.
	resource *resource.Resource
}

var _ logs.LoggerProvider = &LoggerProvider{}

func (lp LoggerProvider) Logger(name string, options ...logs.LoggerOption) logs.Logger {

	//TODO implement me
	panic("implement me")
}

var _ logs.LoggerProvider = &LoggerProvider{}

func NewLoggerProvider(opts ...LoggerProviderOption) *LoggerProvider {
	o := loggerProviderConfig{}

	// TODO: o = applyLoggerProviderEnvConfigs(o)

	for _, opt := range opts {
		o = opt.apply(o)
	}

	// TODO: o = ensureValidLoggerProviderConfig(o)

	lp := &LoggerProvider{
		namedLogger: make(map[instrumentation.Scope]*logger),
		resource:    o.resource,
	}

	global.Info("LoggerProvider created", "config", o)

	spss := make(logRecordProcessorStates, 0, len(o.processors))
	for _, sp := range o.processors {
		spss = append(spss, newLogsProcessorState(sp))
	}
	lp.logProcessors.Store(&spss)

	return lp

}

//func (lp LoggerProvider) Send(rol ReadWriteLogRecord) {
//
//	// add resource level attributes
//	rol.SetResource(lp.cfg.resource)
//
//	// process log
//	for _, p := range lp.cfg.processors {
//		p.OnEmit(rol)
//	}
//}

func (p *LoggerProvider) getLogsProcessors() logRecordProcessorStates {
	return *(p.logProcessors.Load())
}

func (p LoggerProvider) Shutdown(ctx context.Context) error {
	// This check prevents deadlocks in case of recursive shutdown.
	if p.isShutdown.Load() {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	// This check prevents calls after a shutdown has already been done concurrently.
	if !p.isShutdown.CompareAndSwap(false, true) { // did toggle?
		return nil
	}

	var retErr error
	for _, sps := range p.getLogsProcessors() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var err error
		sps.state.Do(func() {
			err = sps.lp.Shutdown(ctx)
		})
		if err != nil {
			if retErr == nil {
				retErr = err
			} else {
				// Poor man's list of errors
				retErr = fmt.Errorf("%v; %v", retErr, err)
			}
		}
	}
	p.logProcessors.Store(&logRecordProcessorStates{})
	return retErr

}

// ForceFlush immediately exports all logs that have not yet been exported for
// all the registered span processors.
func (p *LoggerProvider) ForceFlush(ctx context.Context) error {
	spss := p.getLogsProcessors()
	if len(spss) == 0 {
		return nil
	}

	for _, sps := range spss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := sps.lp.ForceFlush(ctx); err != nil {
			return err
		}
	}
	return nil
}
