package main

import (
	"fmt"
	"log"
	"os"
	"path"

	goscrapper "github.com/tkxkd0159/goutil/go-scrapper"
	"github.com/tkxkd0159/goutil/io"
)

func main() {

	for {
		lang := ""
		pageLimit := 10
		filename := "jobs.csv"

		fmt.Println(" > Enter programming language which you are looking for")
		_, err := fmt.Scanln(&lang)
		io.CheckErr(err, "", 1)
		fmt.Println(" > Enter the number per page")
		_, err = fmt.Scanln(&pageLimit)
		io.CheckErr(err, "", 1)
		fmt.Println(" > Enter the file name to save")
		_, err = fmt.Scanln(&filename)
		io.CheckErr(err, "", 1)

		baseURL := "https://indeed.com"
		target := fmt.Sprintf("%s/jobs?q=%s", baseURL, lang)

		var pgnums []int
		log.Println("==> aggregating page numbers...")
		for i := 0; ; i++ {
			url := goscrapper.GetJobURL(target, i, pageLimit)
			if !goscrapper.CheckNext(url) {
				pgnums = append(pgnums, i)
				break
			} else {
				pgnums = append(pgnums, i)
			}
		}

		log.Println("==> start to get job infos...")
		var allJobs []goscrapper.JobInfo
		tmpSave := make(chan []goscrapper.JobInfo, len(pgnums))
		for _, pn := range pgnums {
			go goscrapper.GetJobInfos(tmpSave, target, pn, pageLimit)
		}

		for _, _ = range pgnums {
			allJobs = append(allJobs, <-tmpSave...)
		}

		goscrapper.WriteJobInfos(allJobs, path.Join(os.Getenv("HOME"), filename))

		for {
			var endSig string
			fmt.Println(" > Do you wanna search another job positions (y/n)")
			_, _ = fmt.Scanln(&endSig)
			if endSig == "y" {
				break
			} else if endSig == "n" {
				os.Exit(0)
			} else {
				fmt.Println("this is wrong input")
			}
		}

	}

}
