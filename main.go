package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	URL      = "http://104.45.179.41/"
	N_THREAD = 100
	N_LOOP   = 10
)

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}

	var wg sync.WaitGroup

	for id := 0; id < N_THREAD; id++ {
		id := id
		wg.Add(1)
		go func(id int) {
			for i := 0; i < N_LOOP; i++ {
				resp, err := client.Get(URL)
				if err != nil {
					log.Fatalln(err)
				}
				defer resp.Body.Close()
				log.Printf("goroutine id: %4d, index: %4d, status: %d, content length: %d\n", id, i, resp.StatusCode, resp.ContentLength)
			}
			wg.Done()
		}(id)
	}

	wg.Wait()
}
