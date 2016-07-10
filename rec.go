package scrapi

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"sync"
	"time"
)

const NUM_WORKERS = 16

type Result struct {
	sync.RWMutex
	State map[url.URL]string
	Date  time.Time
	Notes string
}

func Rec(servRoot *url.URL, fname string, client ReqFunc) error {
	queue := make(chan *url.URL, NUM_WORKERS*4)
	done := make(chan struct{})
	pending := &sync.WaitGroup{}
	result := &Result{
		State: make(map[url.URL]string),
		Date:  time.Now(),
	}

	pending.Add(1)
	queue <- servRoot

	for i := 0; i < num_workers; i++ {
		go Worker(queue, done, pending, client, result)
	}
	pending.Wait()

	for i := 0; i < num_workers; i++ {
		done <- struct{}{}
	}

	return DumpToFile(result, fname)
}

func Worker(queue chan *url.URL, done chan struct{}, wg *sync.WaitGroup, client ReqFunc, result *Result) {
	log.Print("starting worker")
	for {
		select {
		case <-done:
			return
		case nextUrl := <-queue:
			body, urls, err := Process(nextUrl, client)
			if err != nil {
				log.Print(err)
				continue
			}
			result.Lock()
			result.State[nextUrl.Path] = string(body)
			result.Unlock()
			for _, u := range urls {
				result.RLock()
				_, ok := result.State[u.Path]
				result.RUnlock()
				if !ok {
					wg.Add(1)
					queue <- u
				}
			}
			wg.Add(-1)
		}
	}
}

func DumpToFile(data *Result, fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	return enc.Encode(data)
}
