package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/anwyho/gophercises/link"
)

func main() {
	urlFlag := flag.String("url", "http://localhost:3000/story/", "the URL that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the maximum number of links to follow")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println(page)
	}
	fmt.Println(len(pages))
}

type empty struct{}

func bfs(urlStr string, maxDepth int) (ret []string) {
	seen := make(map[string]empty)
	var q map[string]empty
	nq := map[string]empty{
		urlStr: empty{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]empty)
		for u := range q {
			fmt.Printf("Processing %s\n", u)
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = empty{}
			for _, link := range get(u) {
				nq[link] = empty{}
			}
		}
	}
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	return filter(baseURL, hrefs(resp.Body, baseURL), withPrefix(baseURL.String()), withSubstring("the"))
}

func hrefs(body io.Reader, base *url.URL) (ret []string) {
	links, _ := link.Parse(body)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		case strings.HasPrefix(l.Href, "#"):
		default:
			ref := &url.URL{
				Scheme: base.Scheme,
				Host:   base.Host,
				Path:   path.Join(base.Path, l.Href),
			}
			ret = append(ret, ref.String())
		}
	}
	return
}

func filter(base *url.URL, links []string, keepFns ...func(link string) bool) (ret []string) {
LinkParse:
	for _, link := range links {
		keep := true
		for _, fn := range keepFns {
			if fn(link) == false {
				keep = false
				continue LinkParse
			}
		}
		if keep {
			ret = append(ret, link)
		}
	}
	return
}

func withPrefix(pfx string) func(link string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func withSubstring(s string) func(link string) bool {
	return func(link string) bool {
		return strings.Contains(link, s)
	}

}
