package feed

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Fetcher struct{}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

type fetchResponse struct {
	URL  string
	Body []byte
	Err  error
}

func (f *Fetcher) FetchFromURLs(urls []string) []fetchResponse {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	results := make([]fetchResponse, 0, len(urls))
	for i, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("[%d/%d x %s -> %v\n]", i+1, len(urls), url, err)
			results = append(results, fetchResponse{URL: url, Err: err})
			continue
		}

		body, err := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			return nil
		}

		fmt.Printf("[%d/%d V %s -> %s \n]", i+1, len(urls), url, resp.Status)
		results = append(results, fetchResponse{URL: url, Body: body, Err: err})
	}

	return results
}
