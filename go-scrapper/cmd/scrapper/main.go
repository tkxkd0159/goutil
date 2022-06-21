package main

import (
	"errors"
	"fmt"
)

func main() {
	baseURL := "https://indeed.com/jobs"
	lang := "python"
	pageNum := 0
	target := fmt.Sprintf("%s?q=%s&start=%d", baseURL, lang, pageNum)
	fmt.Println(target)
	myerr := errors.New("my error")

}
