package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// Sentry initializes Sentry, for error reporting
func Sentry() error {
	dsn, err := GetEnvVar(DSNEnvVarName)
	if err != nil {
		return err
	}
	return sentry.Init(sentry.ClientOptions{Dsn: dsn})
}

// ListenAddress determines what port to listen on and falls back to a default
// port if the environment does not supply a port
func ListenAddress() string {
	port := os.Getenv(PortEnvVarName)
	if port == "" {
		port = DefaultPort
	}
	address := fmt.Sprintf(":%s", port)
	return address
}

// ErrorMap turns the supplied error into a map with "error" as the key
func ErrorMap(err error) map[string]string {
	errMap := make(map[string]string)
	errMap["error"] = err.Error()
	return errMap
}

// RequestDebugMiddleware dumps the incoming HTTP request to the log for inspection
func RequestDebugMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Errorf("Unable to read request body for debugging: error %#v", err)
				}
				req, err := httputil.DumpRequest(r, true)
				if err != nil {
					log.Errorf("Unable to dump cloned request for debugging: error %#v", err)
				}
				if IsDebug() {
					log.Printf("Raw request: %v", string(req))
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				next.ServeHTTP(w, r)
			},
		)
	}
}

// LogStartupError is used to e.g log fatal startup errors.
// It logs, attempts to report the error to StackDriver then panics/crashes.
func LogStartupError(ctx context.Context, err error) {
	errorClient := StackDriver(ctx)
	if err != nil {
		if errorClient != nil {
			errorClient.Report(errorreporting.Entry{Error: err})
		}
		log.WithFields(log.Fields{"error": err}).Error("Server startup error")
	}
}

// DecodeJSONToTargetStruct maps JSON from a HTTP request to a struct.
func DecodeJSONToTargetStruct(w http.ResponseWriter, r *http.Request, targetStruct interface{}) {
	err := json.NewDecoder(r.Body).Decode(targetStruct)
	if err != nil {
		WriteJSONResponse(w, ErrorMap(err), http.StatusBadRequest)
		return
	}
}

// ConvertStringToInt converts a supplied string value to an integer.
// It writes an error to the JSON response writer if the conversion fails.
func ConvertStringToInt(w http.ResponseWriter, val string) int {
	converted, err := strconv.Atoi(val)
	if err != nil {
		WriteJSONResponse(w, ErrorMap(err), http.StatusInternalServerError)
		return -1 // sentinel value
	}
	return converted
}

// StackDriver initializes StackDriver logging, error reporting, profiling etc
func StackDriver(ctx context.Context) *errorreporting.Client {
	// project setup
	projectID, err := GetEnvVar(GoogleCloudProjectIDEnvVarName)
	if err != nil {
		log.WithFields(log.Fields{
			"environment variable name": GoogleCloudProjectIDEnvVarName,
			"error":                     err,
		}).Error("Unable to determine the Google Cloud Project, can't set up StackDriver")
		return nil
	}

	// logging
	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.WithFields(log.Fields{
			"project ID": projectID,
			"error":      err,
		}).Error("Unable to initialize logging client")
		return nil
	}
	defer CloseStackDriverLoggingClient(loggingClient)

	// error reporting
	errorClient, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: AppName,
		OnError: func(err error) {
			log.WithFields(log.Fields{
				"project ID":   projectID,
				"service name": AppName,
				"error":        err,
			}).Info("Unable to initialize error client")
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to initialize error client")
		return nil
	}
	defer CloseStackDriverErrorClient(errorClient)

	// tracing
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"project ID": projectID,
			"error":      err,
		}).Info("Unable to initialize tracing")
		return errorClient // the error client is already initialized, return it
	}
	trace.RegisterExporter(exporter)

	// profiler
	err = profiler.Start(profiler.Config{
		Service:        AppName,
		ServiceVersion: AppVersion,
		ProjectID:      projectID,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"project ID":      projectID,
			"service name":    AppName,
			"service version": AppVersion,
			"error":           err,
		}).Info("Unable to initialize profiling")
		return errorClient // the error client is already initialized, return it
	}

	return errorClient
}

// WriteJSONResponse writes the content supplied via the `source` parameter to
// the supplied http ResponseWriter. The response is returned with the indicated
// status.
func WriteJSONResponse(w http.ResponseWriter, source interface{}, status int) {
	w.WriteHeader(status) // must come first...otherwise the first call to Write... sets an implicit 200
	content, errMap := json.Marshal(source)
	if errMap != nil {
		msg := fmt.Sprintf("error when marshalling %#v to JSON bytes: %#v", source, errMap)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, errMap = w.Write(content)
	if errMap != nil {
		msg := fmt.Sprintf(
			"error when writing JSON %s to http.ResponseWriter: %#v", string(content), errMap)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

// CloseStackDriverLoggingClient closes a StackDriver logging client and logs any arising error.
//
// It was written to be defer()'d in servrer initialization code.
func CloseStackDriverLoggingClient(loggingClient *logging.Client) {
	err := loggingClient.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to close StackDriver logging client")
	}
}

// CloseStackDriverErrorClient closes a StackDriver error client and logs any arising error.
//
// It was written to be defer()'d in servrer initialization code.
func CloseStackDriverErrorClient(errorClient *errorreporting.Client) {
	err := errorClient.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to close StackDriver error client")
	}
}
