package web

import (
	"log"
	"net/http"
)

type Resurl struct {
	URL        string
	Status     string
	StatusCode int
	Err        error
}

func CheckURLs(urls []string) map[string]Resurl {
	total := len(urls)

	ans := make(map[string]Resurl)
	recvCh := make(chan Resurl, total)

	for _, url := range urls {
		muURL := url

		go func() {
			res, err := http.Get(muURL)

			if err != nil {
				recvCh <- Resurl{URL: muURL, Err: err, Status: "", StatusCode: 0}
				return
			}

			if res.StatusCode > InfoSuccessUpperBound {
				recvCh <- Resurl{URL: muURL, Status: res.Status, StatusCode: res.StatusCode, Err: nil}
				return
			}
			recvCh <- Resurl{URL: muURL, Status: res.Status, StatusCode: res.StatusCode, Err: nil}
		}()
	}

	for i := 0; i < total; i++ {
		got := <-recvCh
		ans[got.URL] = got
	}

	return ans
}

func CheckCode(res *http.Response) {
	if res.StatusCode != OK {
		log.Printf("ERROR : %d %s\n", res.StatusCode, res.Status)
	}
}
