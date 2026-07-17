package feed

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Minute}

type Result struct {
	URL  string
	Body []byte
	Err  error
}

func FetchFromURLs(urls []string) []Result {
	results := make([]Result, 0, len(urls))

	for i, url := range urls {
		resp, err := httpClient.Get(url)
		if err != nil {
			fmt.Printf("[%d/%d x %s -> %v\n]", i+1, len(urls), url, err)
			results = append(results, Result{URL: url, Err: err})
			continue
		}

		body, err := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			return nil
		}
		fmt.Printf("[%d/%d V %s -> %s \n]", i+1, len(urls), url, resp.Status)
		results = append(results, Result{URL: url, Body: body, Err: err})
	}

	return results
}
