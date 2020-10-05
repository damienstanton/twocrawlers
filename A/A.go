// Package A implements a tour webcrawler that uses serialized access to a shared map
package A

import (
	"fmt"
	"sync"
)

// A Fetcher retrieves data by crawling a given URL
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type cache struct {
	wg sync.WaitGroup
	mu sync.Mutex

	m map[string]struct{}
}

func (h *cache) seen(url string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.m[url]; ok {
		return true
	}
	h.m[url] = struct{}{}
	return false
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, c *cache) {
	defer c.wg.Done()
	if depth <= 0 || c.seen(url) {
		return
	}

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		c.wg.Add(1)
		go Crawl(u, depth-1, fetcher, c)
	}

	return
}
