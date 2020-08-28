package base_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestCSS(t *testing.T) {
	hf := base.CSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
}

func TestVisitCSS(t *testing.T) {
	hf := base.VisitCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
}

func TestProfileCSS(t *testing.T) {
	hf := base.ProfileCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
}

func TestHistoryCSS(t *testing.T) {
	hf := base.HistoryCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
}

func TestInvalidCSS(t *testing.T) {
	hf := base.InvalidCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
}
