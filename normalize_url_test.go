package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorContains string
	}{
		{
			name:     "https url",
			inputURL: "https://test.dev/path",
			expected: "test.dev/path",
		},
		{
			name:     "http url",
			inputURL: "http://test.dev/path",
			expected: "test.dev/path",
		},
		{
			name:     "ending slash",
			inputURL: "https://test.dev/path/",
			expected: "test.dev/path",
		},
		{
			name:     "no change needed",
			inputURL: "test.dev/path",
			expected: "test.dev/path",
		},
		{
			name:     "base url",
			inputURL: "https://test.dev",
			expected: "test.dev",
		},
		{
			name:     "remove scheme and capitals",
			inputURL: "https://TEST.dev/pATH",
			expected: "test.dev/path",
		},
		{
			name:          "invalid",
			inputURL:      ":\\invalid.dev",
			expected:      "",
			errorContains: "couldn't parse URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil && (!strings.Contains(err.Error(), tc.errorContains) || tc.errorContains == "") {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
