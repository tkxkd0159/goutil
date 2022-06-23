package go_scrapper

import "testing"

func TestCheckNext(t *testing.T) {
	testurl := "https://indeed.com/jobs?q=python&start=200&limit=50"
	if CheckNext(testurl) != true {
		t.Error("There is no next scope")
	}
}
