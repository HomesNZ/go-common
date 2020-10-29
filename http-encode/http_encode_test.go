package http_encode

import (
	"encoding/json"
	"errors"
	"github.com/HomesNZ/go-common/logger"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// should set 200 status and correct header
func TestEncodeOKResponse(t *testing.T) {
	logger := logger.Init(
		logger.Level("info"),
	).WithField("service", "example")

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	func(w http.ResponseWriter, r *http.Request) {
		EncodeOKResponse(logger, w, "test")
	}(w, req)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
}

// should set 300 status and correct header
func TestEncodeResponse(t *testing.T) {
	logger := logger.Init(
		logger.Level("info"),
	).WithField("service", "example")

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	func(w http.ResponseWriter, r *http.Request) {
		EncodeResponse(logger, w, http.StatusMultipleChoices, "test")
	}(w, req)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, 300)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
}

// should set 500 status, correct content-type and return default err msg
func TestEncodeErrorResponseDefault(t *testing.T) {
	logger := logger.Init(
		logger.Level("info"),
	).WithField("service", "example")

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	func(w http.ResponseWriter, r *http.Request) {
		EncodeErrorResponse(logger, w, errors.New("Error !!!"))
	}(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, resp.StatusCode, 500)
	assert.Equal(t, string(body), "{\"error\":\"Something went wrong\"}\n")
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
}

// should set 400 status, correct content-type and return default err msg
func TestEncodeErrorResponseSyntaxError(t *testing.T) {
	logger := logger.Init(
		logger.Level("info"),
	).WithField("service", "example")

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(`{Username": "test"}`))
	w := httptest.NewRecorder()

	func(w http.ResponseWriter, r *http.Request) {
		var i interface{}
		err := json.NewDecoder(r.Body).Decode(&i)
		EncodeErrorResponse(logger, w, err)
	}(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, resp.StatusCode, 400)
	assert.Equal(t, string(body), "{\"error\":\"invalid character 'U' looking for beginning of object key string\"}\n")
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
}

//TODO: add tests for validation.Errors and schema.MultiError
