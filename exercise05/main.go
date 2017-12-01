package main

import (
	"fmt"
	"io"
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

	pages := get(urlForMapping)
	for _, page := range pages {
		fmt.Println(page)
	}
}

func get(urlForMapping string) []string {
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

	pages, err := hrefs(resp.Body, base)
	if err != nil {
		panic(err)
	}
	return filter(pages, withPrefix(base))
}

func hrefs(html io.Reader, baseURL string) ([]string, error) {
	links, err := link.Parse(html)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, baseURL+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret, nil
}

func filter(links []string, pred func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if pred(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
