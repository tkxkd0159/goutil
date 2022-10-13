package web

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckURLs(t *testing.T) {
	tcs := []string{
		"https://naver.com",
		"unavailable",
		"https://www.daum.net",
		"https://www.example.com",
		"https://www.google.com",
	}
	expected := map[string]Resurl{
		"https://naver.com":       {"https://naver.com", "200 OK", 200, nil},
		"https://www.daum.net":    {"https://www.daum.net", "200 OK", 200, nil},
		"https://www.example.com": {"https://www.example.com", "200 OK", 200, nil},
		"https://www.google.com":  {"https://www.google.com", "200 OK", 200, nil},
		"unavailable":             {"unavailable", "", 0, &url.Error{Op: "Get", URL: "unavailable", Err: errors.New("unsupported protocol scheme \"\"")}},
	}

	got := CheckURLs(tcs)
	for _, tc := range tcs {
		assert.Equal(t, expected[tc], got[tc])
	}
}
