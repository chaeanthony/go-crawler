package main

import "testing"

func TestGetUrlsFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
					<a href="https://other.com/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "no urls in body",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<span>no urls</span>
				</body>
			</html>
			`,
			expected: []string{},
		},
		{
			name:     "invalid body urls",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="meh">
						<span>test</span>
					</a>
					<a href="invalid">
						<span>test</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{},
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urls, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - %s FAIL: expected urls: %v, actual: %v", i, tc.name, tc.expected, urls)
			}
		})
	}
}
