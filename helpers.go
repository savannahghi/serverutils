package serverutils

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
)

// BoolEnv gets and parses a boolean environment variable
func BoolEnv(envVarName string) bool {
	envVar, err := GetEnvVar(envVarName)
	if err != nil {
		return false
	}

	val, err := strconv.ParseBool(envVar)
	if err != nil {
		return false
	}

	return val
}

// IsDebug returns true if debug has been turned on in the environment
func IsDebug() bool {
	return BoolEnv(DebugEnvVarName)
}

// IsRunningTests returns true if debug has been turned on in the environment
func IsRunningTests() bool {
	return BoolEnv(IsRunningTestsEnvVarName)
}

// GetEnvVar retrieves the environment variable with the supplied name and fails
// if it is not able to do so
func GetEnvVar(envVarName string) (string, error) {
	envVar := os.Getenv(envVarName)
	if envVar == "" {
		return "", fmt.Errorf("the environment variable '%s' is not set", envVarName)
	}

	return envVar, nil
}

// NewErrorResponseWriter returns an initialized ErrorResponseWriter
func NewErrorResponseWriter(err error) *ErrorResponseWriter {
	return &ErrorResponseWriter{
		err: err,
		rec: httptest.NewRecorder(),
	}
}

// ErrorResponseWriter is a http.ResponseWriter that always errors on attempted writes.
//
// It is necessary for tests.
type ErrorResponseWriter struct {
	err error
	rec *httptest.ResponseRecorder
}

// Header delegates reading of headers to the underlying response writer
func (w *ErrorResponseWriter) Header() http.Header {
	return w.rec.Header()
}

// Write always returns the supplied error on any attempt to write.
func (w *ErrorResponseWriter) Write([]byte) (int, error) {
	return 0, w.err
}

// WriteHeader delegates writing of headers to the underlying response writer
func (w *ErrorResponseWriter) WriteHeader(statusCode int) {
	w.rec.WriteHeader(statusCode)
}

// MustGetEnvVar returns the value of the environment variable with the indicated name or panics.
// It is intended to be used in the INTERNALS of the server when we can guarantee (through orderly
// coding) that the environment variable was set at server startup.
// Since the env is required, kill the app if the env is not set. In the event a variable is not super
// required, set a sensible default or don't call this method
func MustGetEnvVar(envVarName string) string {
	val, err := GetEnvVar(envVarName)
	if err != nil {
		msg := fmt.Sprintf("mandatory environment variable %s not found", envVarName)
		log.Panic(msg)
	}

	return val
}
