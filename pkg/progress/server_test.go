package progress_test

import (
	"github.com/kubernetes-simulator/simulator/pkg/progress"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NullLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	return log
}

func Test_ServeHTTP_GET_missing_scenario(t *testing.T) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.NewHTTPHandler(progress.NewLocalStateProvider(NullLogger()),
		NullLogger())

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request")
	assert.Equal(t, "Missing scenario name\n", rr.Body.String(), "Wrong message")
}

func Test_ServeHTTP_GET_no_progress(t *testing.T) {
	req, err := http.NewRequest("GET", "/?scenario=test", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.NewHTTPHandler(progress.NewLocalStateProvider(NullLogger()),
		NullLogger())

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "Expected 404 Not Found")
}

func Test_ServeHTTP_GET_with_progress(t *testing.T) {
	makeProgress("test-scenario")
	req, err := http.NewRequest("GET", "/?scenario=test-scenario", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.NewHTTPHandler(progress.NewLocalStateProvider(NullLogger()),
		NullLogger())

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK")
}

func Test_ServeHTTP_POST(t *testing.T) {
	req, err := http.NewRequest("POST", "/", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.NewHTTPHandler(progress.NewLocalStateProvider(NullLogger()), NullLogger())

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK")
}

func Test_ServeHTTP_invalid_method(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.NewHTTPHandler(progress.NewLocalStateProvider(NullLogger()),
		NullLogger())

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code,
		"Expected 405 Method Not Allowed")
}

func Test_ServeHTTP_POST_with_garbage(t *testing.T) {
	req, err := http.NewRequest("POST", "/", strings.NewReader("invalid"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.NewHTTPHandler(progress.NewLocalStateProvider(NullLogger()),
		NullLogger())

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request")
}
