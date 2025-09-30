package serverutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/profiler"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Sentry initializes Sentry, for error reporting
func Sentry() error {
	dsn, err := GetEnvVar(DSNEnvVarName)
	if err != nil {
		return err
	}

	environment, err := GetEnvVar(Environment)
	if err != nil {
		return err
	}

	traceSampleRate, err := GetEnvVar(TraceSampleRateEnvVarName)
	if err != nil {
		return err
	}

	if traceSampleRate == "" {
		traceSampleRate = "1.0"
	}

	sampleRate, err := strconv.ParseFloat(traceSampleRate, 64)
	if err != nil {
		return err
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		EnableTracing:    true,
		TracesSampleRate: sampleRate,
	})
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
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Errorf("Unable to read request body for debugging: error %#v", err)
				}

				if IsDebug() {
					req, err := httputil.DumpRequest(r, true)
					if err != nil {
						log.Errorf("Unable to dump cloned request for debugging: error %#v", err)
					}

					log.Printf("Raw request: %v", string(req))
				}

				r.Body = io.NopCloser(bytes.NewBuffer(body))

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
// TODO: Move to common helpers
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
	exporter, err := trace.New(
		trace.WithProjectID(os.Getenv("GOOGLE_CLOUD_PROJECT")),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"project ID": projectID,
			"error":      err,
		}).Info("Unable to initialize tracing")

		return errorClient // the error client is already initialized, return it
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(AppName),
		semconv.ServiceVersionKey.String(AppVersion),
		attribute.String("environment", os.Getenv("ENVIRONMENT")),
	)

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

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
// TODO: Move to common helpers
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

// PrepareServer is the signature of a function that Knows how to prepare & initialise the server
type PrepareServer func(ctx context.Context, port int, allowedOrigins []string) *http.Server

func randomPort() int {
	/* #nosec G404 */
	rng := rand.New(rand.NewSource(time.Now().Unix()))

	minVal := 32768
	maxVal := 60999
	/* #nosec G404 */
	port := rng.Intn(maxVal-minVal+1) + minVal

	return port
}

// StartTestServer starts up test server
func StartTestServer(ctx context.Context, prepareServer PrepareServer, allowedOrigins []string) (*http.Server, string, error) {
	// prepare the server
	port := randomPort()
	srv := prepareServer(ctx, port, allowedOrigins)
	baseURL := fmt.Sprintf("http://localhost:%d", port)

	if srv == nil {
		return nil, "", fmt.Errorf("nil test server")
	}

	// set up the TCP listener
	// this is done early so that we are sure we can connect to the port in
	// the tests; backlogs will be sent to the listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil || l == nil {
		return nil, "", fmt.Errorf("unable to listen on port %d: %w", port, err)
	}

	if IsDebug() {
		log.Printf("LISTENING on port %d", port)
	}

	// start serving
	go func() {
		err := srv.Serve(l)
		if err != nil {
			if IsDebug() {
				log.Printf("serve error: %s", err)
			}
		}
	}()

	// the cleanup of this server (deferred shutdown) needs to occur in the
	// acceptance test that will use this
	return srv, baseURL, nil
}

// HealthStatusCheck endpoint to check if the server is working.
func HealthStatusCheck(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}
