package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	queryFmt = "https://api.github.com/users/%s/repos?page=%d;per_page=%d"
	MaxPage  = 100
)

func main() {
	/* FLAG */
	maxqueries := flag.Int("m", -1, "maximum number of queries")
	outfilename := flag.String("o", "repos.txt", "name of the outfile")
	pagination := flag.Int("p", 30, "pagination; results per page")
	timeout := flag.Float64("t", 3.0, "maxlatency on failed request ( default 3.0s )")
	username := flag.String("u", "", "username (required)")
	echo := flag.Bool("e", false, "echo (default false)")

	flag.Parse()
	checkFlag(*username != "", fmt.Errorf("username required"))
	checkFlag(0 < *pagination && *pagination < MaxPage, fmt.Errorf("pagination but be between (0,%d)", MaxPage))
	/* FLAG */

	/* REQUEST */
	var client = http.Client{
		Timeout: time.Second * time.Duration(*timeout),
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		log.Fatalf("creating new request: %v\n", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("Host", "api.github.com")
	req.Header.Set("User-Agent", "maxsei-gitlist")
	/* REQUEST */

	/* API */
	var repositories []GitHubResponse
	for i := 0; i < *maxqueries || *maxqueries < 0; i++ {
		// set the request url
		urlQuery := fmt.Sprintf(queryFmt, *username, i, *pagination)
		err := updateReqURL(req, urlQuery)
		if err != nil {
			log.Fatalf("updating request url: %v\n", err)
		}
		// echo the page grabbed
		if *echo {
			fmt.Printf("getting page %d", i)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("getting response: %v\n", err)
		}
		defer resp.Body.Close()
		// read the raw bytes of the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading body: %v\n", err)
		}
		if len(body) <= 2 {
			break
		}
		// decompress bytes from gzip
		greader, err := gzip.NewReader(bytes.NewBuffer(body))
		if err != nil {
			log.Fatalf("creating gzip reader: %v\n", err)
		}
		defer greader.Close()
		// decode into page and append to all repositories
		var page []GitHubResponse
		err = json.NewDecoder(greader).Decode(&page)
		if err != nil {
			log.Fatalf("decoding json: %v\n", err)
		}
		repositories = append(repositories, page...)
		if len(page) < *pagination {
			break
		}
	}
	/* API */

	/* IO */
	file, err := os.OpenFile(*outfilename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("opening file: %v\n", err)
	}
	defer file.Close()
	var repobytes []byte
	for _, repo := range repositories {
		if repo.Owner.Login != *username {
			continue
		}
		repobytes = append(repobytes, []byte(repo.HTMLURL+"\n")...)
	}
	io.Copy(file, bytes.NewBuffer(repobytes))
	/* IO */
}

// checkFlag will check to see if a boolean is true and will log a fatal error
// that is specified if the condition is not met
func checkFlag(boolean bool, errmsg error) {
	if !boolean {
		fmt.Fprintf(os.Stderr, "%v\n", errmsg)
		flag.Usage()
		os.Exit(1)
	}
}

func updateReqURL(req *http.Request, str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	req.URL = u
	return nil
}
