package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/assert"
)

// $ curl localhost:8000/generate
// {"code":"400","error":"No credentials provided"}
func TestGenerateHandlerNoCredentials(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/generate", nil)
	w := httptest.NewRecorder()
	GenerateHashHandler(w, r)

	type Response struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}

	text := w.Body.Bytes()
	data := Response{}
	json.Unmarshal(text, &data)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, data.Code, "400")
	assert.Equal(t, data.Error, "No credentials provided")
}

// $ curl -u test:test localhost:8000/generate
// {"code":"200","hash":"$2a$10$GuNajf3reF3NgXoWy5VIFuxBvxILSqvdvRn1WgnTvpd30T9W0ARNS","status":"OK"}
func TestGenerateHandlerWithCredentials(t *testing.T) {

	r, _ := http.NewRequest("GET", "/generate", nil)
	r.SetBasicAuth("test", "test")
	w := httptest.NewRecorder()

	GenerateHashHandler(w, r)

	type Response struct {
		Code   string `json:"code"`
		Status string `json:"status"`
		Hash   string `json:"hash"`
	}

	text := w.Body.Bytes()
	data := Response{}
	json.Unmarshal(text, &data)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, data.Code, "200")
	assert.Equal(t, data.Status, "OK")
	assert.Equal(t, len(data.Hash), 60)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResourceHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/resource", nil)
	w := httptest.NewRecorder()
	ResourceHandler(w, r)

	type Response struct {
		Code   string `json:"code"`
		Data   string `json:"data"`
		Status string `json:"status"`
	}

	text := w.Body.Bytes()
	data := Response{}
	json.Unmarshal(text, &data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, data.Code, "200")
	assert.Equal(t, data.Status, "OK")
	assert.Equal(t, data.Data, "Here is some data you have access to")
}

func TestAuthMiddleware(t *testing.T) {
	tokens["demo"] = "demo"
	handlerToTest := basicAuthMiddleware(ResourceHandler)
	r := httptest.NewRequest("GET", "/data/demo", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Token", "demo")
	r = mux.SetURLVars(r, map[string]string{"resource": "data", "name": "demo"})
	handlerToTest.ServeHTTP(w, r)

	type Response struct {
		Code   string `json:"code"`
		Data   string `json:"data"`
		Status string `json:"status"`
	}

	text := w.Body.Bytes()
	data := Response{}
	json.Unmarshal(text, &data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, data.Code, "200")
	assert.Equal(t, data.Status, "OK")
	assert.Equal(t, data.Data, "Here is some data you have access to")
}

func TestAuthMiddlewareNotAuthenticated(t *testing.T) {
	tokens["demo"] = "demo"
	handlerToTest := basicAuthMiddleware(ResourceHandler)
	r := httptest.NewRequest("GET", "/data/demo", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Token", "baddemo")
	r = mux.SetURLVars(r, map[string]string{"resource": "data", "name": "demo"})
	handlerToTest.ServeHTTP(w, r)

	type Response struct {
		Code   string `json:"code"`
		Data   string `json:"data"`
		Status string `json:"status"`
	}

	text := w.Body.Bytes()
	data := Response{}
	json.Unmarshal(text, &data)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, data.Code, "403")
	assert.Equal(t, data.Status, "Forbidden")
	assert.Equal(t, data.Data, "")
}

func TestAuthMiddlewareNotAuthorized(t *testing.T) {
	tokens["demo"] = "demo"
	handlerToTest := basicAuthMiddleware(ResourceHandler)
	r := httptest.NewRequest("GET", "/baddata/demo", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Token", "demo")
	r = mux.SetURLVars(r, map[string]string{"resource": "baddata", "name": "demo"})
	handlerToTest.ServeHTTP(w, r)

	type Response struct {
		Code   string `json:"code"`
		Data   string `json:"data"`
		Status string `json:"status"`
	}

	text := w.Body.Bytes()
	data := Response{}
	json.Unmarshal(text, &data)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, data.Code, "401")
	assert.Equal(t, data.Status, "Unauthorized")
	assert.Equal(t, data.Data, "")
}