package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
)

func TestHandlePostTask_TableDriven(t *testing.T) {
	// Define a mock function or simulate the expected behavior of ConvertFormat.
	// If ConvertFormat is defined elsewhere, these test scenarios assume it returns 
	// specific error codes based on inputs.

	// Cleanup after test
	// os.Remove("../file-storage/completed/")
	
	tests := []struct {
		name           string
		requestPayload map[string]string
		expectedStatus int
		expectedSubstr string // substring to check in response body
	}{
		{
			name: "Success - Valid Task",
			requestPayload: map[string]string{
				"file_name":     "cat.jpeg",
				"output_format": "png",
			},
			expectedStatus: http.StatusOK,
			expectedSubstr: `"status":"success"`,
		},
		{
			name: "Failure - Invalid Output Format",
			requestPayload: map[string]string{
				"file_name":     "cat.jpeg",
				"output_format": "mp4",
			},
			expectedStatus: http.StatusBadRequest,
			expectedSubstr: `"status":"failed"`,
		},
		{
			name: "Failure - Non-existent File",
			requestPayload: map[string]string{
				"file_name":     "does_not_exist.txt",
				"output_format": "pdf",
			},
			expectedStatus: http.StatusBadRequest,
			expectedSubstr: `"status":"failed"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error

			if tc.requestPayload != nil {
				bodyBytes, err = json.Marshal(tc.requestPayload)
				if err != nil {
					t.Fatalf("failed to marshal test payload: %v", err)
				}
			} else {
				bodyBytes = []byte("{invalid-json") // simulate bad JSON
			}

			req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Call the worker handler directly
			handlePostTask(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			// Verify Status Code
			if res.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, res.StatusCode)
			}

			// Verify Response Content Type
			if ct := res.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("expected Content-Type application/json, got %q", ct)
			}
		})
	}
}