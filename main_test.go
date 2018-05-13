package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var h = buildRouter()

func TestAuthorization(t *testing.T) {
	content := strings.NewReader("test payload")
	// Good authorization
	req := httptest.NewRequest("POST", "/objects/auth", content)
	req.Header.Set("Authorization", defaultToken)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code, got %v want %v",
			status, http.StatusNoContent)
	}

	// Bad authorization
	req = httptest.NewRequest("POST", "/objects/auth", content)
	req.Header.Set("Authorization", "apple_juice")
	rr = httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong code, got %v want %v",
			status, http.StatusUnauthorized)
	}

	// No authorization
	req = httptest.NewRequest("POST", "/objects/auth", content)
	rr = httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong code, got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestObject(t *testing.T) {
	p := "test payload"
	content := strings.NewReader(p)

	// Create
	req := httptest.NewRequest("POST", "/objects/test", content)
	req.Header.Set("Authorization", defaultToken)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong code, got %v want %v",
			status, http.StatusNoContent)
	}

	// Get
	req = httptest.NewRequest("GET", "/objects/test", nil)
	rr = httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong code, got %v want %v",
			status, http.StatusOK)
	}
	buf, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if c := string(buf); c != p {
		t.Errorf("handler returned invalid payload, got %v want %v",
			c, p)
	}

	// Ensure GET deleted
	req = httptest.NewRequest("GET", "/objects/test", nil)
	rr = httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong code, got %v want %v",
			status, http.StatusNotFound)
	}

	// Ensure POST checks for content
	req = httptest.NewRequest("POST", "/objects/test", nil)
	req.Header.Set("Authorization", defaultToken)
	rr = httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong code, got %v want %v",
			status, http.StatusBadRequest)
	}
}
