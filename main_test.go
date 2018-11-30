package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/zap-api/app"
	"gitlab.com/zap-api/config"
)

var a app.App

func init() {
	os.Setenv("HOST", ":8080")
	os.Setenv("ZAP_PROPERTIES_ENDPOINT", "http://grupozap-code-challenge.s3-website-us-east-1.amazonaws.com/sources/source-2.json")
	config := config.GetConfig()
	a.Initialize(config)
}

// TestEmptyHeader tests the API for an empty Header
func TestEmptyHeader(t *testing.T) {

	req, err := http.NewRequest("GET", "/properties", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, response.Body.String(), `{"error":"Source not accepted."}`)
}

// TestWrongHeader Tests the API for a wrong Source (Company)
func TestWrongHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/properties", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("source", "xxx")

	response := executeRequest(req)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, response.Body.String(), `{"error":"Source not accepted."}`)
}

// TestZAP Tests for the correct response
func TestZAP(t *testing.T) {
	req, err := http.NewRequest("GET", "/properties", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("source", "zap")

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body.String())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.GetAllProperties)
	handler.ServeHTTP(rr, req)
	return rr
}
