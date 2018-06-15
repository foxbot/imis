package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

var h = buildRouter()

// copy the const into a var so we can make a pointer to when invoking request()
var dt = defaultToken

func TestAuthorization(t *testing.T) {
	content := strings.NewReader("test payload")
	// Good authorization
	rr := request("POST", "/auth", content, &dt, nil)
	assertCode(t, rr, http.StatusNoContent, "auth:good")

	// Bad authorization
	token := "apple_juice"
	rr = request("POST", "/auth", content, &token, nil)
	assertCode(t, rr, http.StatusUnauthorized, "auth:bad")

	// No authorization
	rr = request("POST", "/auth", content, nil, nil)
	assertCode(t, rr, http.StatusUnauthorized, "auth:none")
}

func TestObject(t *testing.T) {
	p := "test payload"
	content := strings.NewReader(p)
	// don't make the unit tests last forever :)
	defaultExpires = 200 * time.Millisecond

	// Create
	rr := request("POST", "/test", content, &dt, nil)
	assertCode(t, rr, http.StatusNoContent, "obj:create")

	// Get
	rr = request("GET", "/test", nil, nil, nil)
	assertCode(t, rr, http.StatusOK, "obj:get")

	buf, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if c := string(buf); c != p {
		t.Errorf("handler returned invalid payload, got %v want %v",
			c, p)
	}

	// Check for deletion
	deleteTest := make(chan bool)
	go func() {
		time.Sleep(defaultExpires)

		rr := request("GET", "/test", nil, nil, nil)
		assertCode(t, rr, http.StatusNotFound, "obj:delete")

		deleteTest <- true
	}()

	// Custom expires after
	ea := "400" // 400ms
	content.Reset(p)
	rr = request("POST", "/test_expires", content, &dt, &ea)
	assertCode(t, rr, http.StatusNoContent, "obj:create_custom_expire")

	deleteCustomTest := make(chan bool)
	go func() {
		time.Sleep(425 * time.Millisecond)

		rr := request("GET", "/test_expires", nil, nil, nil)
		assertCode(t, rr, http.StatusNotFound, "obj:delete_custom_expire")

		deleteCustomTest <- true
	}()

	// Ensure POST checks for content
	rr = request("POST", "/test", nil, &dt, nil)
	assertCode(t, rr, http.StatusBadRequest, "obj:post_needs_content")

	// Ensure POST with Expires-After checks range
	testEa := func(val string, name string) {
		rr := request("POST", "/test_ea_range", content, &dt, &val)
		assertCode(t, rr, http.StatusBadRequest, name)
	}
	testEa(strconv.Itoa(minExpires-1), "obj:post_min_expire")
	testEa(strconv.Itoa(maxExpires+1), "obj:post_max_expire")

	// Ensure list endpoint returns a JSON
	listTest := make(chan bool)
	go func() {
		expected := "test_list"
		content.Reset(p)
		rr := request("POST", "/" + expected, content, &dt, nil)
		assertCode(t, rr, http.StatusNoContent, "obj:list_create")

		rr = request("GET", "", nil, &token, nil)
		assertCode(t, rr, http.StatusOK, "obj:list")

		result := map[string]string{}
		if err = json.NewDecoder(rr.Result().Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		found := false
		for _, val := range result {
			if val == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s in JSON response, does not exist.", expected)
		}
		listTest <- true
	}()

	// Clean up "parallel" tests
	<-deleteTest
	<-deleteCustomTest
	<-listTest
}

// save a few code-duplication-trees
func request(action string, key string, body io.Reader, auth *string, ea *string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/objects%s", key)
	req := httptest.NewRequest(action, url, body)
	if auth != nil {
		req.Header.Set("Authorization", *auth)
	}
	if ea != nil {
		req.Header.Set("X-Delete-After", *ea)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}
func assertCode(t *testing.T, rr *httptest.ResponseRecorder, want int, handler string) {
	if status := rr.Code; status != want {
		t.Errorf("%s returned wrong code, got %v want %v (%s)",
			handler, status, want, rr.Body.String())
	}
}
