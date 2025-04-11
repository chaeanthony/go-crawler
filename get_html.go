package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrNotHTML = errors.New("content-type not text/html")
)

func getHTML(rawURL string) (string, error) {
	baseUrl, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(baseUrl.String())
	if err != nil {
		return "", fmt.Errorf("couldn't get %s. got: %v", baseUrl, err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("couldn't read body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to get html. got: code: %v", resp.Status)
	} else if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", ErrNotHTML
	}

	return string(body), nil
}
