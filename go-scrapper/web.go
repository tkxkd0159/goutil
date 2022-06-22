package go_scrapper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tkxkd0159/goutil"
	"github.com/tkxkd0159/goutil/web"
	"log"
	"net/http"
)

type JobInfo struct {
	Title       string
	CompanyName string
	Location    string
	Salary      string
	JobType     string
	Summary     string
	URL         string
}

func GetPages(url string) (int, error) {
	pages := 0
	res, err := http.Get(url)
	if err != nil {
		return pages, err
	}
	web.CheckCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, err
	}
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages, nil
}

func GetPage(baseURL string, pagenum int, pageLimit int) []JobInfo {
	link := fmt.Sprintf("%s&start=%d&limit=%d", baseURL, pagenum*pageLimit, pageLimit)
	res, err := http.Get(link)
	goutil.CheckErr(err, "", 1)
	web.CheckCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	goutil.CheckErr(err, "", 1)

	var infos []JobInfo

	jobCards := doc.Find(".job_seen_beacon .resultContent")
	jobCards.Each(func(i int, s *goquery.Selection) {
		var jkurl string
		info := JobInfo{}

		title, ok := s.Find(".jobTitle span").Attr("title")
		if ok != true {
			log.Println("there is no title attribute on this jobTitle header")
		}
		jkid, ok := s.Find("a").Attr("data-jk")
		if ok != true {
			log.Println("there is no data-jk attribute on this tag")
		} else {
			jkurl = fmt.Sprintf("https://indeed.com/viewjob?jk=%s", jkid)
		}

		company := s.Find(".companyInfo")
		companyName := company.Find(".companyName > a").Text()
		companyLoc := company.Find("div.companyLocation").Text()

		metadata := s.Find(".metadataContainer")
		salary := metadata.Find(".estimated-salary > span").Text()
		jobType := metadata.Find("[aria-label='Job type']").Parent().Text()

		info.Title = title
		info.CompanyName = companyName
		info.Location = companyLoc
		info.Salary = salary
		info.JobType = jobType
		info.URL = jkurl

		infos = append(infos, info)
	})

	iter := 0
	doc.Find(".jobCardShelfContainer").Each(func(i int, s *goquery.Selection) {
		summary := s.Find(".underShelfFooter li").Text()
		infos[iter].Summary = summary
		iter++
	})

	return infos
}
