package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	N_THREAD = 1000
	N_LOOP   = 1
)

var (
	URLS = []string{"http://104.45.179.41/", "http://52.186.102.90/"}
)

func main() {
	var wg sync.WaitGroup
	var totalCount int32
	var errCount int32

	for id := 0; id < N_THREAD; id++ {
		id := id
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			tr := &http.Transport{
				DisableCompression: true,
				MaxConnsPerHost:    1,
			}
			client := &http.Client{Transport: tr}

			for i := 0; i < N_LOOP; i++ {
				atomic.AddInt32(&totalCount, 1)

				url := URLS[r.Intn(len(URLS))]

				resp, err := client.Get(url)
				if err != nil {
					atomic.AddInt32(&errCount, 1)
					log.Printf("goroutine id: %4d, index: %2d, error: %s\n", id+1, i+1, err)
					return // terminate goroutine
				}
				defer resp.Body.Close()

				io.Copy(io.Discard, resp.Body) // necessary for connection to be reused
				log.Printf("goroutine id: %5d, index: %2d, status: %d, content length: %d\n", id+1, i+1, resp.StatusCode, resp.ContentLength)
			}
		}(id)
	}

	wg.Wait()

	log.Println("finished")
	log.Printf("total requests: %5d, failed requests: %5d", totalCount, errCount)
}
