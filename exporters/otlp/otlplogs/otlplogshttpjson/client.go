/*
Copyright Agoda Services Co.,Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package otlplogshttpjson

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	internal "github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/internal"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/internal/retry"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/otlpconfig"
	"github.com/golang/protobuf/jsonpb"
	"go.opentelemetry.io/otel"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const contentTypeJson = "application/json"

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		return w
	},
}

// Keep it in sync with golang's DefaultTransport from net/http! We
// have our own copy to avoid handling a situation where the
// DefaultTransport is overwritten with some different implementation
// of http.RoundTripper or it's modified by other package.
var ourTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

type client struct {
	name        string
	cfg         otlpconfig.SignalConfig
	generalCfg  otlpconfig.Config
	requestFunc retry.RequestFunc
	client      *http.Client
	stopCh      chan struct{}
	stopOnce    sync.Once
	marshaller  jsonpb.Marshaler
}

var _ otlplogs.Client = (*client)(nil)

// NewClient creates a new HTTP logs client.
func NewClient(opts ...Option) otlplogs.Client {
	cfg := otlpconfig.NewHTTPConfig(asHTTPOptions(opts)...)

	httpClient := &http.Client{
		Transport: ourTransport,
		Timeout:   cfg.Logs.Timeout,
	}
	if cfg.Logs.TLSCfg != nil {
		transport := ourTransport.Clone()
		transport.TLSClientConfig = cfg.Logs.TLSCfg
		httpClient.Transport = transport
	}

	stopCh := make(chan struct{})
	return &client{
		name:        "logs",
		cfg:         cfg.Logs,
		generalCfg:  cfg,
		requestFunc: cfg.RetryConfig.RequestFunc(evaluate),
		stopCh:      stopCh,
		client:      httpClient,
		marshaller:  jsonpb.Marshaler{OrigName: true},
	}
}

// Start does nothing in a HTTP client.
func (d *client) Start(ctx context.Context) error {
	// nothing to do
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// Stop shuts down the client and interrupt any in-flight request.
func (d *client) Stop(ctx context.Context) error {
	d.stopOnce.Do(func() {
		close(d.stopCh)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// retryableError represents a request failure that can be retried.
type retryableError struct {
	throttle int64
}

// evaluate returns if err is retry-able. If it is and it includes an explicit
// throttling delay, that delay is also returned.
func evaluate(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}

	rErr, ok := err.(retryableError)
	if !ok {
		return false, 0
	}

	return true, time.Duration(rErr.throttle)
}

func (d *client) contextWithStop(ctx context.Context) (context.Context, context.CancelFunc) {
	// Unify the parent context Done signal with the client's stop
	// channel.
	ctx, cancel := context.WithCancel(ctx)
	go func(ctx context.Context, cancel context.CancelFunc) {
		select {
		case <-ctx.Done():
			// Nothing to do, either cancelled or deadline
			// happened.
		case <-d.stopCh:
			cancel()
		}
	}(ctx, cancel)
	return ctx, cancel
}

func (d *client) newRequest(body []byte) (request, error) {
	u := url.URL{Scheme: d.getScheme(), Host: d.cfg.Endpoint, Path: d.cfg.URLPath}
	r, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return request{Request: r}, err
	}

	r.Header.Set("User-Agent", otlpconfig.GetUserAgentHeader())

	for k, v := range d.cfg.Headers {
		r.Header.Set(k, v)
	}
	r.Header.Set("Content-Type", contentTypeJson)

	req := request{Request: r}
	switch Compression(d.cfg.Compression) {
	case NoCompression:
		r.ContentLength = (int64)(len(body))
		req.bodyReader = bodyReader(body)
	case GzipCompression:
		// Ensure the content length is not used.
		r.ContentLength = -1
		r.Header.Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		var b bytes.Buffer
		gz.Reset(&b)

		if _, err := gz.Write(body); err != nil {
			return req, err
		}
		// Close needs to be called to ensure body if fully written.
		if err := gz.Close(); err != nil {
			return req, err
		}

		req.bodyReader = bodyReader(b.Bytes())
	}

	return req, nil
}

// bodyReader returns a closure returning a new reader for buf.
func bodyReader(buf []byte) func() io.ReadCloser {
	return func() io.ReadCloser {
		return io.NopCloser(bytes.NewReader(buf))
	}
}

func (d *client) getScheme() string {
	if d.cfg.Insecure {
		return "http"
	}
	return "https"
}

// request wraps an http.Request with a resettable body reader.
type request struct {
	*http.Request

	// bodyReader allows the same body to be used for multiple requests.
	bodyReader func() io.ReadCloser
}

// reset reinitializes the request Body and uses ctx for the request.
func (r *request) reset(ctx context.Context) {
	r.Body = r.bodyReader()
	r.Request = r.Request.WithContext(ctx)
}

// newResponseError returns a retryableError and will extract any explicit
// throttle delay contained in headers.
func newResponseError(header http.Header) error {
	var rErr retryableError
	if s, ok := header["Retry-After"]; ok {
		if t, err := strconv.ParseInt(s[0], 10, 64); err == nil {
			rErr.throttle = t
		}
	}
	return rErr
}

func (e retryableError) Error() string {
	return "retry-able request failure"
}

func (d *client) UploadLogs(ctx context.Context, protoLogs []*logspb.ResourceLogs) error {

	// Export the logs using the OTLP logs exporter client
	exportLogs := &collogspb.ExportLogsServiceRequest{
		ResourceLogs: protoLogs,
	}

	// Serialize the OTLP logs payload
	rawRequest, _ := d.marshaller.MarshalToString(exportLogs)
	ctx, cancel := d.contextWithStop(ctx)
	defer cancel()

	request, err := d.newRequest([]byte(rawRequest))
	if err != nil {
		return err
	}

	return d.requestFunc(ctx, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request.reset(ctx)
		resp, err := d.client.Do(request.Request)
		if err != nil {
			return err
		}

		if resp != nil && resp.Body != nil {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					otel.Handle(err)
				}
			}()
		}

		switch resp.StatusCode {
		case http.StatusOK:
			// Success, do not retry.
			// Read the partial success message, if any.
			var respData bytes.Buffer
			if _, err := io.Copy(&respData, resp.Body); err != nil {
				return err
			}

			if respData.Len() != 0 {
				var response collogspb.ExportLogsServiceResponse
				if err := jsonpb.UnmarshalString(respData.String(), &response); err != nil {
					return err
				}

				// TODO: partialsuccess can't be handled properly by OTEL as current otlp.internal.PartialSuccess is custom
				// need to have that interface in official OTEL otlp.internal package
				if response.PartialSuccess != nil {
					msg := response.PartialSuccess.GetErrorMessage()
					n := response.PartialSuccess.GetRejectedLogRecords()
					if n != 0 || msg != "" {
						err := internal.LogRecordPartialSuccessError(n, msg)
						otel.Handle(err)
					}
				}
			}
			return nil

		case http.StatusTooManyRequests, http.StatusServiceUnavailable:
			// Retry-able failures.  Drain the body to reuse the connection.
			if _, err := io.Copy(io.Discard, resp.Body); err != nil {
				otel.Handle(err)
			}
			return newResponseError(resp.Header)
		default:
			return fmt.Errorf("failed to send to %s: %s", request.URL, resp.Status)
		}
	})
}

// MarshalLog is the marshaling function used by the logging system to represent this Client.
func (d *client) MarshalLog() interface{} {
	return struct {
		Type     string
		Endpoint string
		Insecure bool
	}{
		Type:     "otlphttphttpjson",
		Endpoint: d.cfg.Endpoint,
		Insecure: d.cfg.Insecure,
	}
}
