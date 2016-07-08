package scrapi

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"sync"
)

type Result struct {
	sync.RWMutex
	m map[url.URL]string
}

func Rec(servRoot *url.URL, fname string, client ReqFunc) error {
	num_workers := 16
	queue := make(chan *url.URL, num_workers*4)
	done := make(chan struct{})
	idle := &sync.WaitGroup{}
	result := &Result{m: make(map[url.URL]string)}

	idle.Add(num_workers)
	queue <- servRoot
	for i := 0; i < num_workers; i++ {
		go Worker(queue, done, idle, client, result)
	}
	idle.Wait()
	for i := 0; i < num_workers; i++ {
		done <- struct{}{}
	}
	log.Printf("result: %v", result)
	return DumpToFile(result, fname)
}

func Worker(queue chan *url.URL, done chan struct{}, wg *sync.WaitGroup, client ReqFunc, result *Result) {
	log.Print("starting worker")
	for {
		wg.Add(-1)
		select {
		case <-done:
			return
		case nextUrl := <-queue:
			wg.Add(1)
			body, urls, err := Process(nextUrl, client)
			if err != nil {
				log.Print(err)
				continue
			}
			result.Lock()
			result.m[*nextUrl] = string(body)
			result.Unlock()
			for _, u := range urls {
				result.RLock()
				_, ok := result.m[*u]
				result.RUnlock()
				if !ok {
					queue <- u
				}
			}
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
