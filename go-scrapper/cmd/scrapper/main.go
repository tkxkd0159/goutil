package main

import (
	"fmt"
	"github.com/tkxkd0159/goutil"
	go_scrapper "github.com/tkxkd0159/goutil/go-scrapper"
)

func main() {

	for {
		lang := ""
		pageLimit := 10

		fmt.Println(" * Enter programming language which you are looking for")
		_, err := fmt.Scanln(&lang)
		goutil.CheckErr(err, "", 1)
		fmt.Println(" * Enter the number per page")
		_, err = fmt.Scanln(&pageLimit)
		goutil.CheckErr(err, "", 1)

		baseURL := "https://indeed.com"
		target := fmt.Sprintf("%s/jobs?q=%s", baseURL, lang)
		totalPages, _ := go_scrapper.GetPages(target)

		for i := 0; i < totalPages; i++ {
			go_scrapper.GetPage(target, i, pageLimit)
		}

	}

	//testurl := "https://indeed.com/jobs?q=python&start=200&limit=50"
	//res, _ := http.Get(testurl)
	//defer res.Body.Close()
	//addMore, _ := goquery.NewDocumentFromReader(res.Body)
	//addMore.Find(".pagination > .pagination-list").Each(func(i int, s *goquery.Selection) {
	//	fmt.Println(s.Find("[aria-label=\"Previous\"]").Length())
	//	fmt.Println(s.Find("[aria-label=\"Next\"]").Length())
	//})

}
