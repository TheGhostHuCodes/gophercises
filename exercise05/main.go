package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TheGhostHuCodes/gophercises/exercise04/link"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Printf("Recieved %d arguments (%v), expected 1 argument\n", len(args), args)
		os.Exit(1)
	}

	urlForMapping := args[0]

	resp, err := http.Get(urlForMapping)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	requestURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: requestURL.Scheme,
		Host:   requestURL.Host,
	}
	base := baseURL.String()

	links, err := link.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	for _, href := range hrefs {
		fmt.Println(href)
	}
}
