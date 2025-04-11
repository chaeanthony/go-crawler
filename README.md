# Go Crawler

Website crawler written in Golang. Recursively searches internal links starting from base url provided. Returns report of number of links for each page.

### Installation

1. Fork/Clone the repository:

2. Install dependencies:

- [Go](https://golang.org/doc/install)

```bash
go mod download
```

### Usage

```bash
go run . <website> <maxConcurrency> <maxPages>
```

maxConcurrency = number of workers to concurrently search
maxPages = set number of pages desired to process
