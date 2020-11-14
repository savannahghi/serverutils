package base_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

const (
	testUserEmail     = "automated.test.bewell.ediproxy.main@healthcloud.co.ke"
	listPayersQuery   = `{"query": "query ListPayers() { payers { id, name, slade_code } }"}`
	listContactsQuery = `{"query": "query ListContacts() { contacts { id, contact, contact_type, is_home, is_personal, is_work } }"}`
)

func TestSentry(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "default case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := base.Sentry(); (err != nil) != tt.wantErr {
				t.Errorf("Sentry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrorMap(t *testing.T) {
	err := fmt.Errorf("test error")
	errMap := base.ErrorMap(err)
	if errMap["error"] == "" {
		t.Errorf("empty error key in errMap")
	}
	if errMap["error"] != "test error" {
		t.Errorf("expected the error value to be '%s', got '%s'", "test error", errMap["error"])
	}
}

func TestRequestDebugMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mw := base.RequestDebugMiddleware()
	h := mw(next)

	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw, req)

	rw1 := httptest.NewRecorder()
	reader1 := ioutil.NopCloser(bytes.NewBuffer([]byte("will be closed")))
	err := reader1.Close()
	assert.Nil(t, err)
	req1 := httptest.NewRequest(http.MethodPost, "/", reader1)
	h.ServeHTTP(rw1, req1)
}

func TestLogStartupError(t *testing.T) {
	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
				err: fmt.Errorf("this is a test error"),
			},
		},
		{
			name: "nil error",
			args: args{
				ctx: context.Background(),
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.LogStartupError(tt.args.ctx, tt.args.err)
		})
	}
}

func TestDecodeJSONToTargetStruct(t *testing.T) {
	type target struct {
		A string `json:"a,omitempty"`
	}
	targetStruct := target{}

	type args struct {
		w            http.ResponseWriter
		r            *http.Request
		targetStruct interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good case",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/",
					bytes.NewBuffer([]byte(
						"{\"a\":\"1\"}",
					)),
				),
				targetStruct: &targetStruct,
			},
		},
		{
			name: "invalid / failed decode",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/",
					nil,
				),
				targetStruct: &targetStruct,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.DecodeJSONToTargetStruct(tt.args.w, tt.args.r, tt.args.targetStruct)
		})
	}
}

func Test_convertStringToInt(t *testing.T) {
	tests := map[string]struct {
		val                string
		rw                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"successful_conversion": {
			val:                "768",
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 200,
		},
		"failed_conversion": {
			val:                "not an int",
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "{\"error\":\"strconv.Atoi: parsing \\\"not an int\\\": invalid syntax\"}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			base.ConvertStringToInt(tc.rw, tc.val)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
			assert.Equal(t, tc.expectedResponse, tc.rw.Body.String())
		})
	}
}

func Test_StackDriver_Setup(t *testing.T) {
	errorClient := base.StackDriver(context.Background())
	err := fmt.Errorf("test error")
	if errorClient != nil {
		errorClient.Report(errorreporting.Entry{
			Error: err,
		})
	}
}

func TestStackDriver(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base.StackDriver(tt.args.ctx)
			assert.NotNil(t, got)
		})
	}
}

func TestWriteJSONResponse(t *testing.T) {
	unmarshallable := make(chan string) // can't be marshalled to JSON
	errReq := base.NewErrorResponseWriter(fmt.Errorf("ka-boom"))

	type args struct {
		w      http.ResponseWriter
		source interface{}
		status int
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "happy case - successful write",
			args: args{
				w:      httptest.NewRecorder(),
				source: map[string]string{"test_key": "test_value"},
				status: http.StatusOK,
			},
			wantStatus: 200,
		},
		{
			name: "unmarshallable content",
			args: args{
				w:      httptest.NewRecorder(),
				source: unmarshallable,
				status: http.StatusInternalServerError,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "closed recorder",
			args: args{
				w:      errReq,
				source: map[string]string{"test_key": "test_value"},
				status: http.StatusOK,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.WriteJSONResponse(tt.args.w, tt.args.source, tt.args.status)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if ok {
				assert.NotNil(t, rec)
				assert.Equal(t, tt.wantStatus, rec.Code)
			}
			if !ok {
				rec, ok := tt.args.w.(*base.ErrorResponseWriter)
				assert.True(t, ok)
				assert.NotNil(t, rec)
			}
		})
	}
}

func Test_closeStackDriverLoggingClient(t *testing.T) {
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	loggingClient, err := logging.NewClient(context.Background(), projectID)
	assert.Nil(t, err)

	type args struct {
		loggingClient *logging.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case - successful",
			args: args{
				loggingClient: loggingClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.CloseStackDriverLoggingClient(tt.args.loggingClient)
		})
	}
}

func Test_closeStackDriverErrorClient(t *testing.T) {
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	errorClient, err := errorreporting.NewClient(context.Background(), projectID, errorreporting.Config{
		ServiceName: base.AppName,
		OnError: func(err error) {
			log.WithFields(log.Fields{
				"project ID":   projectID,
				"service name": base.AppName,
				"error":        err,
			}).Info("Unable to initialize error client")
		},
	})
	assert.Nil(t, err)

	type args struct {
		errorClient *errorreporting.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case - successful",
			args: args{
				errorClient: errorClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.CloseStackDriverErrorClient(tt.args.errorClient)
		})
	}
}

func TestGetGraphQLHeaders(t *testing.T) {

	authenticatedContext, bearerToken := base.GetAuthenticatedContextAndBearerToken(t)
	authHeader := fmt.Sprintf("Bearer %s", bearerToken)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "context with no authorization header",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "context with authorization header",
			args: args{
				ctx: authenticatedContext,
			},
			want: req.Header{
				"Accept":        "application/json",
				"Content-Type":  "application/json",
				"Authorization": authHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetGraphQLHeaders(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGraphQLHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}

func TestGetBearerTokenHeader(t *testing.T) {

	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Good Test Case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetBearerTokenHeader(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerTokenHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEqual(t, got, "")
		})
	}
}

func TestStartTestServer(t *testing.T) {

	ctx := context.Background()
	srv, baseURL, serverErr := base.StartTestServer(ctx, healthCheckServer, []string{
		"http://localhost:5000",
	})
	if serverErr != nil {
		t.Errorf("Unable to start test server %s", serverErr)
		return
	}
	defer srv.Close()
	if srv == nil {
		t.Errorf("nil test server %s", serverErr)
		return
	}
	if baseURL == "" {
		t.Errorf("empty base url %s", serverErr)
		return
	}
}

func healthCheckRouter() (*mux.Router, error) {
	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error

	r.Use(base.RequestDebugMiddleware())
	r.Path("/health").HandlerFunc(base.HealthStatusCheck)

	return r, nil
}

func healthCheckServer(ctx context.Context, port int, allowedOrigins []string) *http.Server {
	// start up the router
	r, err := healthCheckRouter()
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// start the server
	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)
	h = handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)
	h = handlers.CombinedLoggingHandler(os.Stdout, h)
	h = handlers.ContentTypeHandler(h, "application/json")
	srv := &http.Server{
		Handler: h,
		Addr:    addr,
	}
	log.Infof("Server running at port %v", addr)
	return srv

}
