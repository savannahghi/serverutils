package client

import (
	"errors"
	"net/http"
)

// BlowUpOnClose provides a closer that always returns an error, for testing purposes
type BlowUpOnClose struct{}

// Close on this mock always returns an error
func (rc BlowUpOnClose) Close() error { return errors.New("ka-boom") }

// Read on this mock always reads 0 bytes and returns no error
func (rc BlowUpOnClose) Read(_ []byte) (n int, err error) {
	return 0, nil
}

// BlowUpOnRead provides a reader that always returns an error, for testing purposes
type BlowUpOnRead struct{}

// Close on this mock always succeeds with no error
func (rc BlowUpOnRead) Close() error { return nil }

// Read on this mock always returns an error
func (rc BlowUpOnRead) Read(_ []byte) (n int, err error) { return 0, errors.New("boom") }

// MockHTTPTransportFunc defines the signature of a function that can be assigned to a HTTP client Transport value
type MockHTTPTransportFunc func(req *http.Request) (*http.Response, error)

// RoundTrip always applies the func on which it is a receiver
func (f MockHTTPTransportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// MockHTTPClient returns *http.Client with Transport replaced to avoid making real calls
func MockHTTPClient(fn MockHTTPTransportFunc) *http.Client {
	return &http.Client{
		Transport: MockHTTPTransportFunc(fn),
	}
}
