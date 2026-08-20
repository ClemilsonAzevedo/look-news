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
		res := f.fetch(client, url)
		if res.Err != nil {
			res = f.fetch(client, url)
		}

		if res.Err != nil {
			fmt.Printf("[%d/%d x %s -> %v\n]", i+1, len(urls), url, res.Err)
		} else {
			fmt.Printf("[%d/%d V %s -> success \n]", i+1, len(urls), url)
		}

		results = append(results, res)
	}
	return results
}

func (f *Fetcher) fetch(client http.Client, url string) fetchResponse {
	resp, err := client.Get(url)
	if err != nil {
		return fetchResponse{
			URL: url,
			Err: err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fetchResponse{
			URL: url,
			Err: fmt.Errorf("invalid status: %s", resp.Status),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fetchResponse{
			URL: url,
			Err: err,
		}
	}

	return fetchResponse{
		URL:  url,
		Body: body,
	}
}
