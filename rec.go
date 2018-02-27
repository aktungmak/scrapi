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
	NUM_WORKERS = 32
	QUEUE_SIZE  = NUM_WORKERS * 200
)

type Result struct {
	sync.RWMutex
	State map[string]string
	Date  time.Time
	Host  string
	Root  string
	Note  string
}

type endpointData struct {
	Url  *url.URL
	Body string
}

func Rec(servRoot *url.URL, fname string, client ReqFunc, note string) error {
	worker_input := make(chan *url.URL, QUEUE_SIZE)
	agg_url_input := make(chan *url.URL, QUEUE_SIZE)
	agg_body_input := make(chan endpointData, QUEUE_SIZE)
	stop := make(chan struct{})
	var pending sync.WaitGroup

	result := &Result{
		State: make(map[string]string),
		Date:  time.Now(),
		Root:  servRoot.Path,
		Host:  servRoot.Host,
		Note:  note,
	}

	// send the first url
	pending.Add(1)
	agg_url_input <- servRoot

	// aggregator
	go func() {
		for {
			select {
			case <-stop:
				return
			case u := <-agg_url_input:
				_, ok := result.State[u.Path]
				if !ok { // only queue the URL if it has not been seen
					worker_input <- u
				} else {
					pending.Add(-1)
				}
			case ep := <-agg_body_input:
				result.State[ep.Url.Path] = ep.Body
			}
		}
	}()

	// workers
	for i := 0; i < NUM_WORKERS; i++ {
		go func() {
			for {
				select {
				case <-stop:
					return
				case u := <-worker_input:
					body, urls, err := Process(u, client)
					if err != nil {
						log.Printf("worker err: %s", err)
						pending.Add(-1)
						continue
					}

					// report result to aggregator
					agg_body_input <- endpointData{Url: u, Body: string(body)}

					// send found urls to aggregator
					for _, u := range urls {
						pending.Add(1)
						agg_url_input <- u
					}

					// report we are done
					pending.Add(-1)
				}
			}
		}()
	}

	// wait until there are zero pending urls
	pending.Wait()

	// tell the aggregator and workers to stop
	close(stop)

	log.Printf("captured %d endpoints", len(result.State))
	return DumpToFile(result, fname)
}

func DumpToFile(data *Result, fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	return enc.Encode(data)
}
