package web

import (
	"net/http"
)

type resurl struct {
	URL        string
	Status     string
	StatusCode int
	Err        error
}

func CheckURLs(urls []string) map[string]resurl {

	total := len(urls)

	ans := make(map[string]resurl)
	recvCh := make(chan resurl, total)

	for _, url := range urls {
		muURL := url
		go func() {
			res, err := http.Get(muURL)
			if err != nil {
				recvCh <- resurl{URL: muURL, Err: err}
				return
			}
			if res.StatusCode > 299 {
				recvCh <- resurl{URL: muURL, Status: res.Status, StatusCode: res.StatusCode, Err: nil}
				return

			}
			recvCh <- resurl{URL: muURL, Status: res.Status, StatusCode: res.StatusCode, Err: nil}

		}()
	}

	for i := 0; i < total; i++ {
		select {
		case got := <-recvCh:
			ans[got.URL] = got
		}
	}

	return ans
}
