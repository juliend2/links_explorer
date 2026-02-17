package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleClientError(t *testing.T) {
	var testClientErrorPages ErroneousPages
	assert.Equal(t, testClientErrorPages.Length(), 0, "Should not have erroneous page")
	handleClientError(&testClientErrorPages, 404, "http://www.google.com/", "http://www.google.com/notfound1.html")
	assert.Equal(t, len(testClientErrorPages.Pages["http://www.google.com/"].ClientErrors), 1, "Should have 1 broken link")
	handleClientError(&testClientErrorPages, 404, "http://www.google.com/", "http://www.google.com/notfound1.html")
	assert.Equal(t, len(testClientErrorPages.Pages["http://www.google.com/"].ClientErrors), 1, "Should STILL have 1 broken link")
	handleClientError(&testClientErrorPages, 404, "http://www.google.com/", "http://www.google.com/notfound2.html")
	assert.Equal(t, len(testClientErrorPages.Pages["http://www.google.com/"].ClientErrors), 2, "Should have 2 broken links")
	handleClientError(&testClientErrorPages, 404, "http://www.google.com/", "http://www.google.com/notfound3.html")
	assert.Equal(t, len(testClientErrorPages.Pages["http://www.google.com/"].ClientErrors), 3, "Should have 3 broken links")
	handleClientError(&testClientErrorPages, 404, "http://www.google.com/", "http://www.google.com/notfound4.html")
	assert.Equal(t, testClientErrorPages.Length(), 1, "Should have only one erroneous page")
	// fmt.Println(testClientErrorPages.Pages["http://www.google.com/"].ClientErrors)
	assert.Equal(t, len(testClientErrorPages.Pages["http://www.google.com/"].ClientErrors), 4, "Should have 4 broken links")
}

func TestHandleServerError(t *testing.T) {
	var testServerErrors []string
	assert.Equal(t, len(testServerErrors), 0, "Should have 0 server error")
	handleServerError(&testServerErrors, 500, "http://www.google.com/broken")
	assert.Equal(t, len(testServerErrors), 1, "Should have 1 server error")
	handleServerError(&testServerErrors, 500, "http://www.google.com/broken")
	assert.Equal(t, len(testServerErrors), 1, "Should STILL have 1 server error")
	handleServerError(&testServerErrors, 500, "http://www.google.com/broken2")
	assert.Equal(t, len(testServerErrors), 2, "Should have 2 server error")
}

func TestPreMarshal(t *testing.T) {
	var testClientErrorPages ErroneousPages
	assert.Equal(t, testClientErrorPages.Length(), 0, "Should not have erroneous page")
	handleClientError(&testClientErrorPages, 404, "http://www.google.com/", "http://www.google.com/notfound1.html")
}

