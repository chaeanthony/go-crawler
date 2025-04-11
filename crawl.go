package main

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	baseURL       *url.URL
	visited       map[string]bool
	pages         map[string]int
	maxPages      int // max number of pages to visit. used to limit number of crawls
	maxConcurreny int
	mu            *sync.Mutex
	wg            *sync.WaitGroup
	skippedURLs   []string
	overflowURLs  []string
}

// crawlPage bfs crawls all internal links from baseURL. It retuns a map of urls to the number of times the url has been seen.
// bfs, bc why not.
func (cfg *config) crawlPage() error {
	queue := make(chan string, 1000)
	queue <- cfg.baseURL.String() // init queue
	cfg.wg.Add(1)

	for i := range cfg.maxConcurreny { // deploy workers
		workerId := i + 1
		fmt.Printf("Deployed worker %d\n", workerId)
		go func() {
			for rawUrl := range queue {
				cfg.mu.Lock()
				nPages := len(cfg.pages)
				cfg.mu.Unlock()
				if nPages >= cfg.maxPages {
					cfg.wg.Done()
					continue
				}
				fmt.Printf("Processing: %s\n", rawUrl)
				cfg.processURL(queue, rawUrl)
				cfg.wg.Done()
			}
		}()
	}

	cfg.wg.Wait() // Wait for all workers to finish
	close(queue)  // Close the queue only after all workers have finished

	return nil
}

// processURL is a helper function to process urls from a send-only queue channel. Adds 1 to wg whenever enqueue
func (cfg *config) processURL(queue chan<- string, rawURL string) {
	normURL, err := normalizeURL(rawURL) // used to add to visited and pages maps
	if err != nil {
		fmt.Printf("couldn't normalize %s, got: %v\n", rawURL, err)
		return
	}
	if cfg.isVisited(normURL) {
		fmt.Printf("already visited %s. skipping...\n", rawURL)
		return
	}
	cfg.addVisited(normURL)

	html, err := getHTML(rawURL) // get html body
	if err != nil {
		if errors.Is(err, ErrNotHTML) {
			fmt.Printf("error: %v. proceeding to next...\n", ErrNotHTML)
			return
		}
		fmt.Printf("couldn't get html. ignoring %s\n", rawURL)
		return
	}

	urls, err := getURLsFromHTML(html, rawURL)
	if err != nil {
		fmt.Printf("couldn't get urls html. ignoring %s\n", rawURL)
		return
	}
	cfg.addPage(normURL, len(urls))

	for _, u := range urls {
		// skip other websites
		different, err := cfg.isDifferentURLHost(u)
		if err != nil {
			fmt.Printf("got error checking url host: %v. proceeding to next...\n", err)
			continue
		} else if different {
			cfg.skippedURLs = append(cfg.skippedURLs, u)
			continue
		}

		nu, err := normalizeURL(u) // used to add to visited and pages maps
		if err != nil {
			fmt.Printf("couldn't normalize %s, got: %v\n", u, err)
			continue
		}
		if cfg.isVisited(nu) {
			continue
		}

		select {
		case queue <- u:
			cfg.wg.Add(1)
		default:
			cfg.overflowURLs = append(cfg.overflowURLs, u)
		}
	}
}

func (cfg *config) isDifferentURLHost(rawURL string) (bool, error) {
	// check if valid url
	curURL, err := url.Parse(rawURL)
	if err != nil {
		return false, fmt.Errorf("could not parse %s as url", rawURL)
	}

	if curURL.Hostname() != cfg.baseURL.Hostname() {
		return true, nil
	}

	return false, nil
}

// isVisited checks if normalized url is visited.
func (cfg *config) isVisited(normURL string) bool {
	isVisited := false
	cfg.mu.Lock()
	if cfg.visited[normURL] { // already visited
		isVisited = true
	}
	cfg.mu.Unlock()

	return isVisited
}

// addVisited adds normalized url to visited map
func (cfg *config) addVisited(normURL string) {
	cfg.mu.Lock()
	cfg.visited[normURL] = true
	cfg.mu.Unlock()
}

// addPage adds pages to pages map
func (cfg *config) addPage(normURL string, numPages int) {
	cfg.mu.Lock()
	cfg.pages[normURL] = numPages
	cfg.mu.Unlock()
}
