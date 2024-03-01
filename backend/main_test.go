package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "os"
)

func TestHandleOptions(t *testing.T) {
    req, err := http.NewRequest("OPTIONS", "/api/generate", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handleOptions)

    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status %d, but got %d", http.StatusOK, rr.Code)
    }

    expectedHeaders := map[string]string{
        "Access-Control-Allow-Origin":  "*",
        "Content-Type":                 "text/html; charset=utf-8",
        "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type, X-Requested-With, Authorization",
    }

    for header, expectedValue := range expectedHeaders {
        actualValue := rr.Header().Get(header)
        if actualValue != expectedValue {
            t.Errorf("Expected header value for %s to be %s, but got %s", header, expectedValue, actualValue)
        }
    }
}

func TestGenerateIAC(t *testing.T) {
    reqBody := `{"i_want": "some_iac_prompt"}`
    req, err := http.NewRequest("POST", "/api/generate", strings.NewReader(reqBody))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(generateIAC)

    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status %d, but got %d", http.StatusOK, rr.Code)
    }

    expectedContentType := "text/plain"
    actualContentType := rr.Header().Get("Content-Type")
    if actualContentType != expectedContentType {
        t.Errorf("Expected Content-Type to be %s, but got %s", expectedContentType, actualContentType)
    }

    expectedResponse := "\"What steps can you take to ensure that your infrastructure as code is secure?\""
    actualResponse := rr.Body.String()
    if actualResponse != expectedResponse {
        t.Errorf("Expected response body to be %s, but got %s", expectedResponse, actualResponse)
    }
}

func TestMain(m *testing.M) {
    // Set up necessary initialization code before running tests
    // (e.g., initializing environment variables)

    // Prepare environment variables for testing
    os.Setenv("PORT", ":8086")

    // Run the tests
    exitCode := m.Run()

    // Clean up any resources after running tests (if needed)

    // Exit with the appropriate exit code
    os.Exit(exitCode)
}


