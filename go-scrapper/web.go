package go_scrapper

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tkxkd0159/goutil"
	"github.com/tkxkd0159/goutil/web"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {

		}
	}(res.Body)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, err
	}
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages, nil
}

func GetJobURL(baseURL string, pagenum int, pageLimit int) string {
	return fmt.Sprintf("%s&start=%d&limit=%d", baseURL, pagenum*pageLimit, pageLimit)
}

func GetJobInfos(ch chan<- []JobInfo, baseURL string, pagenum int, pageLimit int) {
	link := GetJobURL(baseURL, pagenum, pageLimit)
	res, err := http.Get(link)
	goutil.CheckErr(err, "", 1)
	web.CheckCode(res)

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {

		}
	}(res.Body)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	goutil.CheckErr(err, "", 1)

	var infos []JobInfo

	jobCards := doc.Find(".job_seen_beacon")
	jobCards.Each(func(i int, s *goquery.Selection) {
		var jkurl string
		info := JobInfo{}

		main := s.Find(".resultContent")

		title, ok := main.Find("a > span").Attr("title")
		if ok != true {
			log.Println("there is no title attribute on this job card")
		}
		jkid, ok := main.Find("a").Attr("data-jk")
		if ok != true {
			log.Println("there is no data-jk attribute on this tag")
		} else {
			jkurl = fmt.Sprintf("https://indeed.com/viewjob?jk=%s", jkid)
		}

		company := main.Find(".companyInfo")
		companyName := company.Find(".companyName > a").Text()
		companyLoc := company.Find("div.companyLocation").Text()

		metadata := main.Find(".metadataContainer")
		salary := metadata.Find(".estimated-salary > span").Text()
		jobType := metadata.Find("[aria-label='Job type']").Parent().Text()

		info.Title = goutil.CleanString(title)
		info.CompanyName = goutil.CleanString(companyName)
		info.Location = goutil.CleanString(companyLoc)
		info.Salary = goutil.CleanString(salary)
		info.JobType = goutil.CleanString(jobType)
		info.URL = goutil.CleanString(jkurl)

		shelf := s.Find(".jobCardShelfContainer")
		summary := shelf.Find(".underShelfFooter li").Text()
		info.Summary = goutil.CleanString(summary)

		infos = append(infos, info)
	})

	ch <- infos
}

func CheckNext(url string) bool {
	var isNext bool
	res, err := http.Get(url)
	goutil.CheckErr(err, "", 1)
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {

		}
	}(res.Body)
	addMore, _ := goquery.NewDocumentFromReader(res.Body)
	addMore.Find(".pagination > .pagination-list").Each(func(i int, s *goquery.Selection) {
		if s.Find("[aria-label=\"Next\"]").Length() > 0 {
			isNext = true
		} else {
			isNext = false
		}
	})
	return isNext
}

func WriteJobInfos(jobs []JobInfo, savepath string) {
	err := os.MkdirAll(filepath.Dir(savepath), 0740)
	goutil.CheckErr(err, "", 0)

	fp, err := os.Create(savepath)
	goutil.CheckErr(err, "", 0)
	w := csv.NewWriter(fp)
	defer w.Flush()

	headers := []string{"Title", "Company name", "Location", "Salary", "Job type", "Sumamry", "URL"}
	err = w.Write(headers)
	goutil.CheckErr(err, "", 0)

	for _, job := range jobs {
		s := []string{job.Title, job.CompanyName, job.Location, job.Salary, job.JobType, job.Summary, job.URL}
		w.Write(s)
	}

}
