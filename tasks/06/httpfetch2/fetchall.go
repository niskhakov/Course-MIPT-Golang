package httpfetch2

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net"
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

var dbg *log.Logger

func init() {
	dbg = log.New(ioutil.Discard, "DBG: ", log.Lmicroseconds)
}

func FetchAll(ctx context.Context, c *http.Client, requests <-chan Request) <-chan Result {
	var wg sync.WaitGroup

	ch := make(chan Result)

	go func() {
		counter := 0
		doThing := true
		for doThing {
			dbg.Println("Main GRTN -> Selecting option")
			select {
				case <- ctx.Done():
					dbg.Println("Main GRTN EVENT -> Got cancelled context")
					doThing = false
				case v, ok := <- requests:
					doThing = ok
					wg.Add(1)
					go doRequest(ctx, c, v, ch, &wg, counter) 
					counter++
			}
		}

		dbg.Printf("Main GRTN EVENT -> Waiting WaitGroup of %d goroutines\n", counter)
		wg.Wait()
		dbg.Println("Main GRTN EVENT -> All Goroutines were executed")
		close(ch)
	}()
	dbg.Println("Main GRTN EVENT -> Exiting function")
    
	return ch
}

func doRequest(ctx context.Context, c *http.Client, request Request, resCh chan<- Result, wg *sync.WaitGroup, id int) {
	defer func() {
		dbg.Printf("doRequest -> Worker %d Done\n", id)
		wg.Done()
	}()
	
	dbg.Printf("doRequest -> Enter worker %d\n", id)

	var res Result

	// Creating req object
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewReader(request.Body))

	if err != nil {
		//dbg.Printf("doRequest -> Error request: %v\n", request)
		//dbg.Printf("doRequest -> Error creating request: %s\n", err)
		return
	}

	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		// dbg.Printf("doRequest -> Error request: %v\n", request)
		// dbg.Printf("doRequest -> Error quering request: %T\n", errors.Unwrap(err))
		if isTimeoutError(err) {
			return
		}
		res.Error = err
		resCh <- res
		return 
	}

	defer resp.Body.Close()

	res.StatusCode = resp.StatusCode
	resCh <- res
}

func isTimeoutError(err error) bool {
	 e, ok := err.(net.Error)
	 return ok && e.Timeout()
}

