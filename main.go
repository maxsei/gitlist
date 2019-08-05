package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	queryFmt = "https://api.github.com/users/%s/repos?page=%d;per_page=%d"
)

func main() {
	/*
		[maximillian@ThinkPadT440 gitlist]$ ./gitlist
		username required
		Usage of ./gitlist:
		  -m int
		        maximum number of queries (default -1) (default -1)
		  -o string
		        outfilename (default repos.txt) (default "repos.txt")
		  -p string
		        password (required)(not stored)
		  -pp int
		        pagination; results per page (default 10) (default 10)
		  -t int
		        maxlatency on failed request ( default 3s ) (default 3)
		  -u string
		        username (required)
	*/

	maxqueries := flag.Int("m", -1, "maximum number of queries (default -1)")
	outfilename := flag.String("o", "repos.txt", "outfile (default \"repos.txt\"")
	password := flag.String("p", "", "password (required)(not stored)")
	pagination := flag.Int("pp", 30, "pagination; results per page (default 10) (default 10)")
	timeout := flag.Float64("t", 3.0, "maxlatency on failed request ( default 3.0s )")
	username := flag.String("u", "", "username (required)")
	echo := flag.Bool("e", false, "echo (default false)")

	flag.Parse()
	checkFlag(*password != "", fmt.Errorf("password required"))
	checkFlag(*username != "", fmt.Errorf("username required"))

}

// checkFlag will check to see if a boolean is true and will log a fatal error
// that is specified if the condition is not met
func checkFlag(boolean bool, errmsg error) {
	if !boolean {
		log.Fatalln(fmt.Errorf("%v\n%s", errmsg, flag.Usage))
	}
}
