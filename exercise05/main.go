package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TheGhostHuCodes/gophercises/exercise04/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	maxDepth := flag.Int("depth", 3, "The maximum depth of links to traverse.")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Recieved %d arguments (%v), expected 1 positional argument, with flags first.\n", flag.NArg(), flag.Args())
		os.Exit(1)
	}
	urlForMapping := flag.Arg(0)

	pages := bfs(urlForMapping, *maxDepth)
	if err := xmlSitemapWriter(os.Stdout, pages); err != nil {
		panic(err)
	}
}

func xmlSitemapWriter(w io.Writer, pages []string) error {
	toXML := urlset{
		Urls:  make([]loc, 0, len(pages)),
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}
	fmt.Fprint(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(toXML)
}

func bfs(urlForMapping string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlForMapping: struct{}{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
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
