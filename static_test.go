package base

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSS(t *testing.T) {
	hf := CSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
	assert.Equal(t, rw.Body.Bytes(), []byte(baseCSSTemplate))
}

func TestVisitCSS(t *testing.T) {
	hf := VisitCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
	assert.Equal(t, rw.Body.Bytes(), []byte(visitCSSTemplate))
}

func TestProfileCSS(t *testing.T) {
	hf := ProfileCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
	assert.Equal(t, rw.Body.Bytes(), []byte(profileCSSTemplate))
}

func TestHistoryCSS(t *testing.T) {
	hf := HistoryCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
	assert.Equal(t, rw.Body.Bytes(), []byte(historyCSSTemplate))
}

func TestInvalidCSS(t *testing.T) {
	hf := InvalidCSS()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	hf(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)
	assert.Equal(t, rw.Body.Bytes(), []byte(invalidCSSTemplate))
}
