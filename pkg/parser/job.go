package parser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"semtest/pkg/cmap"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Job struct {
	Url         string
	Body        io.Reader
	Domain      string
	JobsCh      chan *Job
	QuitCh      chan struct{}
	CheckedUrls *cmap.CMap
	IsChecked   bool
	WG          *sync.WaitGroup
}

func Run(url string, RPS int) {
	re := regexp.MustCompile("^(https?:?)?//[^/]+")
	r := re.FindStringSubmatch(url)
	jobs := make(chan *Job, RPS)

	for i := 0; i < RPS; i++ {
		go GetPage(jobs)
	}

	job := &Job{
		Url:         url,
		Domain:      r[0],
		WG:          &sync.WaitGroup{},
		JobsCh:      jobs,
		CheckedUrls: cmap.NewCMap(),
	}

	job.WG.Add(1)
	jobs <- job

	job.WG.Wait()
}

func GetPage(jobs chan *Job) {
	for job := range jobs {
		if job.Check() {
			job.WG.Done()
			continue
		}

		fmt.Println(job.Url)

		res, err := http.Get(job.Url)
		if err != nil {
			log.Fatal(err)

			return
		}

		job.Body = res.Body
		go job.ParsePage()
		time.Sleep(1 * time.Second)
	}
}

func (job *Job) ParsePage() {
	defer job.WG.Done()
	tokens := html.NewTokenizer(job.Body)

	for {
		t := tokens.Next()

		switch {
		case t == html.ErrorToken:
			return
		case t == html.StartTagToken:
			t := tokens.Token()

			url := getURL(t)
			if url == "" {
				continue
			}

			if url[0] == '/' {
				url = job.Domain + url
			}

			if strings.HasPrefix(url, job.Domain) {
				job.WG.Add(1)
				job.JobsCh <- &Job{
					Url:         url,
					Domain:      job.Domain,
					JobsCh:      job.JobsCh,
					CheckedUrls: job.CheckedUrls,
					WG:          job.WG,
				}
			}
		}
	}

}

func (job *Job) Check() bool {
	if job.IsChecked {
		return true
	}

	return job.CheckedUrls.CheckAdd(job.Url)
}

func getURL(t html.Token) string {
	if t.Data != "a" {
		return ""
	}

	for _, a := range t.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}

	return ""
}
