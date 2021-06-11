package go_utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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
		envErrMsg := fmt.Sprintf("the environment variable '%s' is not set", envVarName)
		return "", fmt.Errorf(envErrMsg)
	}
	return envVar, nil
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
		log.Panicf(msg)
		os.Exit(1)
	}
	return val
}

// CloseRespBody closes the body of the supplied HTTP response
func CloseRespBody(resp *http.Response) {
	if resp != nil {
		err := resp.Body.Close()
		if err != nil {
			log.Println("unable to close response body for request made to ", resp.Request.RequestURI)
		}
	}
}

// GenerateRandomWithNDigits - given a digit generate random numbers
func GenerateRandomWithNDigits(numberOfDigits int) (string, error) {
	rangeEnd := int64(math.Pow10(numberOfDigits) - 1)
	value, _ := rand.Int(rand.Reader, big.NewInt(rangeEnd))
	return strconv.FormatInt(value.Int64(), 10), nil
}

// ExtractBearerToken gets a bearer token from an Authorization header.
// This is expected to contain a Firebase idToken prefixed with "Bearer "
func ExtractBearerToken(r *http.Request) (string, error) {
	return ExtractToken(r, "Authorization", "Bearer")
}

// ExtractToken extracts a token with the specified prefix from the specified header
func ExtractToken(r *http.Request, header string, prefix string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("nil request")
	}
	if r.Header == nil {
		return "", fmt.Errorf("no headers, can't extract bearer token")
	}
	authHeader := r.Header.Get(header)
	if authHeader == "" {
		return "", fmt.Errorf("expected an `%s` request header", header)
	}
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("the `Authorization` header contents should start with `Bearer`")
	}
	tokenOnly := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	return tokenOnly, nil
}

// GenerateRandomEmail allows us to get "unique" emails while still keeping
// one main be.well@bewell.co.ke email account
func GenerateRandomEmail() string {
	return fmt.Sprintf("be.well+%v@bewell.co.ke", time.Now().Unix())
}

// MergeURLValues merges > 1 url.Values into one
func MergeURLValues(values ...url.Values) url.Values {
	merged := url.Values{}
	for _, value := range values {
		for k, v := range value {
			merged[k] = v
		}
	}
	return merged
}

// GetAPIPaginationParams composes pagination parameters for use by a REST API that uses
// offset based pagination
func GetAPIPaginationParams(pagination *PaginationInput) (url.Values, error) {
	if pagination == nil {
		return url.Values{}, nil
	}

	// Treat first or last, when set, literally as page sizes
	// We intentionally "demote" `last`; if both `first` and `last` are specified,
	// `first` will supersede `last`
	var err error
	pageSize := DefaultRESTAPIPageSize
	if pagination.Last > 0 {
		pageSize = pagination.Last
	}
	if pagination.First > 0 {
		pageSize = pagination.First
	}

	// For these "pass through APIs", "after" and "before" should be parseable as ints
	// (literal offsets).
	// We intentionally demote `before` i.e if both `before` and `after` are set,
	// `after` will supersede `before`
	offset := 0
	if pagination.Before != "" {
		offset, err = strconv.Atoi(pagination.Before)
		if err != nil {
			return url.Values{}, fmt.Errorf("expected `before` to be parseable as an int; got %s", pagination.Before)
		}
	}
	if pagination.After != "" {
		offset, err = strconv.Atoi(pagination.After)
		if err != nil {
			return url.Values{}, fmt.Errorf("expected `after` to be parseable as an int; got %s", pagination.After)
		}
	}
	page := int(offset/pageSize) + 1 // page numbers are one based
	values := url.Values{}
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("page_size", fmt.Sprintf("%d", pageSize))
	return values, nil
}
