package scrapi

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	NUM_WORKERS = 4
	QUEUE_SIZE  = 30 + NUM_WORKERS*10
)

type Result struct {
	sync.RWMutex
	State map[string]string
	Date  time.Time
	Host  string
	Root  string
	Notes string
}

func Rec(servRoot *url.URL, fname string, client ReqFunc) error {
	queue := make(chan *url.URL, QUEUE_SIZE)
	stop := make(chan struct{})
	pending := &sync.WaitGroup{}
	result := &Result{
		State: make(map[string]string),
		Date:  time.Now(),
		Root:  servRoot.Path,
		Host:  servRoot.Host,
	}

	pending.Add(1)
	queue <- servRoot

	for i := 0; i < NUM_WORKERS; i++ {
		go Worker(queue, stop, pending, client, result)
	}
	pending.Wait()

	for i := 0; i < NUM_WORKERS; i++ {
		stop <- struct{}{}
	}

	return DumpToFile(result, fname)
}

func Worker(queue chan *url.URL, stop chan struct{}, wg *sync.WaitGroup, client ReqFunc, result *Result) {
	log.Print("starting worker")
	defer log.Print("worker closing")
	for {
		select {
		case <-stop:
			return
		case nextUrl := <-queue:
			body, urls, err := Process(nextUrl, client)
			if err != nil {
				log.Print(err)
				wg.Add(-1)
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
