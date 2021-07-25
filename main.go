package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	URL      = "http://104.45.179.41/"
	N_THREAD = 1000
	N_LOOP   = 10
)

func main() {
	var wg sync.WaitGroup

	for id := 0; id < N_THREAD; id++ {
		id := id
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			tr := &http.Transport{
				DisableCompression: true,
				MaxConnsPerHost:    1,
			}
			client := &http.Client{Transport: tr}

			for i := 0; i < N_LOOP; i++ {
				resp, err := client.Get(URL)
				if err != nil {
					log.Println(err)
					return // terminate goroutine
				}
				defer resp.Body.Close()
				io.Copy(io.Discard, resp.Body) // necessary for connection to be reused
				log.Printf("goroutine id: %4d, index: %2d, status: %d, content length: %d\n", id+1, i+1, resp.StatusCode, resp.ContentLength)
			}
		}(id)
	}

	wg.Wait()
}
