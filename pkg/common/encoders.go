package common

import (
	"fmt"
	"net/url"
	"strings"
)

func EncodeFilenameDots(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	segments := strings.Split(parsedURL.Path, "/")
	if len(segments) == 0 {
		return "", fmt.Errorf("invalid path in URL")
	}

	filename := segments[len(segments)-1]
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex > 0 {
		name := strings.ReplaceAll(filename[:dotIndex], ".", "%2E")
		ext := filename[dotIndex:]
		filename = name + ext
	} else {
		filename = strings.ReplaceAll(filename, ".", "%2E")
	}

	segments[len(segments)-1] = filename
	path := strings.Join(segments, "/")
	result := parsedURL.Scheme + "://" + parsedURL.Host + path
	if parsedURL.RawQuery != "" {
		result += "?" + parsedURL.RawQuery
	}
	if parsedURL.Fragment != "" {
		result += "#" + parsedURL.Fragment
	}

	fmt.Println("Function Output:", result) // Debug
	return result, nil
}
