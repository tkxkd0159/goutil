package web

import (
	"context"
	"io"
	"net/http"
)

func Get(url string, body io.Reader) (*http.Response, error) {
	return GetWithContext(context.Background(), body, url)
}

func GetWithContext(ctx context.Context, body io.Reader, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
