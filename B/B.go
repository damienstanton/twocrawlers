package B

import (
	"sync"
)

// Fetcher ...
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type cache struct {
	depth   int
	fetcher Fetcher
	res     chan response
	req     chan request
	seen    map[string]struct{}
	wg      *sync.WaitGroup
}

type response struct {
	url  string
	body string
	err  error
}

type request struct {
	url   string
	depth int
}

func (c cache) worker(url string, depth int) {
	defer c.wg.Done()

	if depth <= 0 {
		return
	}

	body, urls, err := c.fetcher.Fetch(url)
	c.res <- response{url, body, err}

	if len(urls) == 0 {
		return
	}

	c.wg.Add(len(urls))
	for _, url := range urls {
		c.req <- request{url, depth - 1}
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) (res chan response) {
	res = make(chan response)
	seen := make(map[string]struct{})
	c := cache{
		req:     make(chan request),
		depth:   depth,
		fetcher: fetcher,
		res:     res,
		wg:      &sync.WaitGroup{},
	}

	c.wg.Add(1)
	go func() {
		for req := range c.req {
			if _, ok := seen[req.url]; ok {
				c.wg.Done()
				continue
			}
			seen[req.url] = struct{}{}
			go c.worker(req.url, req.depth)
		}
	}()

	// Wait for the wait group to finish, and then close the channel
	go func() {
		c.wg.Wait()
		close(res)
	}()

	// Send the first crawl request to the channel
	c.req <- request{url, depth}

	return
}
