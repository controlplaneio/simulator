package progress_test

import (
	"github.com/kubernetes-simulator/simulator/pkg/progress"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_ServeHTTP_GET_missing_scenario(t *testing.T) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.HTTPHandler{}

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request")
	assert.Equal(t, "Missing scenario name\n", rr.Body.String(), "Wrong message")
}

func Test_ServeHTTP_GET(t *testing.T) {
	req, err := http.NewRequest("GET", "/?scenario=test", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.HTTPHandler{}

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK")
}

func Test_ServeHTTP_POST(t *testing.T) {
	req, err := http.NewRequest("POST", "/", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.HTTPHandler{}

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK")
}

func Test_ServeHTTP_invalid_method(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hh := progress.HTTPHandler{}

	hh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code,
		"Expected 405 Method Not Allowed")
}
