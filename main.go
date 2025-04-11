package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"syscall"
)

const (
	EXPECTED_ARGS = 3
	ARGS_FORMAT   = "<website> <maxConcurreny> <maxPages>"
)

func main() {
	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)

	// Notify the channel when SIGINT or SIGTERM signals are received
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	args := os.Args
	args = args[1:] // ignore first arg = program path

	if len(args) < EXPECTED_ARGS {
		fmt.Printf("expected args: %s\n", ARGS_FORMAT)
		os.Exit(1)
	}
	if len(args) > EXPECTED_ARGS {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	arg_url := args[0]
	baseURL, err := url.Parse(arg_url)
	if err != nil {
		fmt.Printf("invalid base url: %s", arg_url)
	}
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("failed to convert maxConcurrency to integer. %v", err)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("failed to convert maxPages to integer. %v", err)
	}

	cfg := config{
		pages:         make(map[string]int),
		baseURL:       baseURL,
		visited:       make(map[string]bool),
		maxPages:      maxPages,
		maxConcurreny: maxConcurrency,
		mu:            &sync.Mutex{},
		wg:            &sync.WaitGroup{},
		skippedURLs:   []string{},
		overflowURLs:  []string{},
	}

	// Goroutine to handle the signal
	go func() {
		sig := <-signalChan
		fmt.Printf("\nReceived signal: %s. Cleaning up...\n", sig)
		cfg.outputDetails()
		os.Exit(1)
	}()

	fmt.Printf("Starting crawl of: %s\n", baseURL.String())
	err = cfg.crawlPage()
	if err != nil {
		fmt.Printf("failed to crawl %s. got: %v\n", baseURL.String(), err)
		os.Exit(1)
	}

	cfg.outputDetails()
}

func (cfg *config) outputDetails() {
	fmt.Printf("skipped different websites: %d\n", len(cfg.skippedURLs))
	fmt.Printf("overflow urls because queue buffer was too small: %d\n", len(cfg.overflowURLs))
	cfg.printReport()
}

func (cfg *config) printReport() {
	type kv struct {
		key string
		val int
	}
	var pairs []kv
	for k, v := range cfg.pages {
		pairs = append(pairs, kv{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].val > pairs[j].val // sort by desc order
	})

	fmt.Println("==========================")
	fmt.Printf("REPORT for %s\n", cfg.baseURL.String())
	fmt.Println("==========================")
	for _, pair := range pairs {
		fmt.Printf("Found %d internal links to %s\n", pair.val, pair.key)
	}
}
