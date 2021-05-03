package httpfetch

import (
	"bytes"
	"net/http"
	"sync"
)

type Request struct {
	Method string
	URL    string
	Body   []byte
}

type Result struct {
	StatusCode int
	Error      error
}

func FetchAll(c *http.Client, requests []Request) []Result {

	var wg sync.WaitGroup
	wg.Add(len(requests))
	results := make([]Result, len(requests))

	for i, v := range requests {
		go doRequest(c, v, &results[i], &wg)
	}

	wg.Wait()
	return results
}

func doRequest(c *http.Client, request Request, res *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	// Creating req object
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewReader(request.Body))

	if err != nil {
		// fmt.Printf("Error request: %v\n", request)
		// fmt.Printf("Error creating request: %s\n", err)
		return
	}

	resp, err := c.Do(req)
	if err != nil {
		// fmt.Printf("Error request: %v\n", request)
		// fmt.Printf("Error quering request: %s\n", err)
		res.Error = err
		return
	}

	defer resp.Body.Close()

	// fmt.Println(resp.Body)

	res.StatusCode = resp.StatusCode
}
