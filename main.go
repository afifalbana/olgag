package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func httpClient() *http.Client {

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	return &client
}

func (f *File) buildFileName() {

	fileURL, err := url.Parse(f.url)
	checkError(err)
	path := fileURL.Path
	segments := strings.Split(path, "/")
	f.name = segments[len(segments)-1]
}

func (f *File) createFile() {

	file, err := os.Create(fmt.Sprintf("images/%s", f.name))
	checkError(err)
	f.nameOS = file
}

func (f File) putFile(client *http.Client) {

	resp, err := client.Get(f.url)
	checkError(err)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("Downloading file %s...\n", f.name)
		size, err := io.Copy(f.nameOS, resp.Body)
		checkError(err)
		defer f.nameOS.Close()
		fmt.Printf("Downloaded file %s with size %d\n", f.name, size)
	}
}

func main() {

	a := flag.Int("start", 1, "start page to download")
	b := flag.Int("end", 1, "end page to download")
	flag.Parse()

	timeStart := time.Now()
	defer fmt.Printf("\nDownload done in %v\n", time.Since(timeStart))

	urls := []string{}

	for i := *a; i <= *b; i++ {
		url := fmt.Sprintf("https://img-9gag-fun.9cache.com/photo/%d_700b.jpg", i)
		urls = append(urls, url)
	}

	c := make(chan bool, len(urls))

	for _, u := range urls {

		go func(u string) {
			url := File{url: u}
			url.buildFileName()
			url.createFile()
			url.putFile(httpClient())
			c <- true
		}(u)
	}

	for i := 0; i < len(urls); i++ {
		<-c
	}
}

type File struct {
	url    string
	name   string
	nameOS *os.File
}
