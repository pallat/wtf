// Package httpclient provides functions for creating http.Client
// and calling other APIs
package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pallat/wtf/serror"
)

const (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
)

// OptionFunc defines the function to update *http.Request
type OptionFunc func(*http.Request, ...context.Context)

type transport struct {
	*http.Transport
	options []OptionFunc
}

// NewHTTPClient returns *http.Client instance with the provided options
// and default timeout 5s.
func NewHTTPClient(options ...OptionFunc) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	options = append(options, jsonOption)
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: transport{
			Transport: t,
			options:   options,
		},
	}
}

// NewHTTPClientWithCA returns *http.Client instance with TLS configuration
// and the provided options and default timeout 5s.
// certPool is the root CA certificates to verify the server's certificate.
func NewHTTPClientWithCA(certPool *x509.CertPool, options ...OptionFunc) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	t.TLSClientConfig = &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}
	options = append(options, jsonOption)
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: transport{
			Transport: t,
			options:   options,
		},
	}
}

// Get is a helper function to make a Get request
// and decode the response into the given type.
func Get[RES any](ctx context.Context, client *http.Client, url string) (response Response[RES], err error) {
	return do[any, RES](ctx, client, http.MethodGet, url, nil)
}

// Post is a helper function to make a Post request
// and decode the response into the given type.
func Post[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ) (response Response[RES], err error) {
	return do[REQ, RES](ctx, client, http.MethodPost, url, payload)
}

type Response[T any] struct {
	Code     int
	Response T
}

func do[REQ, RES any](ctx context.Context, client *http.Client, method, url string, payload REQ) (response Response[RES], err error) {
	req, err := newRequest(ctx, client, method, url, payload)
	if err != nil {
		return response, err
	}

	return doRequest[RES](client, req)
}

// DoRequest sends an HTTP request via client given
// and returns the Response with an HTTP response it gets
func DoRequest[RES any](client *http.Client, req *http.Request) (response Response[RES], err error) {
	return doRequest[RES](client, req)
}

func doRequest[RES any](client *http.Client, req *http.Request) (response Response[RES], err error) {
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		err = serror.WrapSkip(err, 3)
		return
	}
	defer resp.Body.Close()

	var v RES
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		err = serror.WrapSkip(err, 3)
		return
	}

	response = Response[RES]{
		Code:     resp.StatusCode,
		Response: v,
	}
	return
}

// NewRequest returns *http.Request
func NewRequest(ctx context.Context, client *http.Client, method, url string, payload any) (*http.Request, error) {
	return newRequest(ctx, client, method, url, payload)
}

func newRequest(ctx context.Context, client *http.Client, method, url string, payload any) (*http.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&payload); err != nil {
		return nil, serror.WrapSkip(err, 3)
	}

	var req *http.Request
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, serror.WrapSkip(err, 3)
	}

	if t, ok := client.Transport.(transport); ok {
		for _, option := range t.options {
			option(req, ctx)
		}
	}

	return req.WithContext(ctx), nil
}
