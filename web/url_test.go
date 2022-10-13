package web

import (
	"errors"
	"net"
	"net/url"
	"os"
	"syscall"
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
		"http://localhost",
	}
	expected := map[string]Resurl{
		"https://naver.com":       {"https://naver.com", "200 OK", 200, nil},
		"https://www.daum.net":    {"https://www.daum.net", "200 OK", 200, nil},
		"https://www.example.com": {"https://www.example.com", "200 OK", 200, nil},
		"https://www.google.com":  {"https://www.google.com", "200 OK", 200, nil},
		"unavailable":             {"unavailable", "", 0, &url.Error{Op: "Get", URL: "unavailable", Err: errors.New("unsupported protocol scheme \"\"")}},
		"http://localhost": {
			"http://localhost",
			"",
			0,
			&url.Error{
				Op: "Get", URL: "http://localhost",
				Err: &net.OpError{
					Op:     "dial",
					Net:    "tcp",
					Source: nil,
					Addr: &net.TCPAddr{
						IP:   net.IP{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7f, 0x00, 0x00, 0x01},
						Port: 80, Zone: "",
					},
					Err: &os.SyscallError{Syscall: "connect", Err: syscall.Errno(0x3d)},
				},
			},
		},
	}

	got := CheckURLs(tcs)
	for _, tc := range tcs {
		assert.Equal(t, expected[tc], got[tc])
	}
}
