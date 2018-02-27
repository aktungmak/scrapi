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
				log.Print("agg stopping")
				return
			case u := <-agg_url_input:
				log.Printf("agg got new url %s", u)
				_, ok := result.State[u.Path]
				if !ok {
					log.Printf("agg sending url %s to workers", u)
					worker_input <- u
				} else {
					log.Printf("agg ignoring url %s", u)
					pending.Add(-1)
				}
			case ep := <-agg_body_input:
				log.Printf("agg storing url %s", ep.Url)
				result.State[ep.Url.Path] = ep.Body
			}
		}
		log.Print("agg fell out of loop")
	}()

	// workers
	for i := 0; i < NUM_WORKERS; i++ {
		go func(id int) {
			for {
				select {
				case <-stop:
					log.Printf("worker%d stopping", id)
					return
				case u := <-worker_input:
					log.Printf("worker%d processing url %s", id, u)
					body, urls, err := Process(u, client)
					if err != nil {
						log.Printf("worker%d error processing %s: %s", id, u, err)
						pending.Add(-1)
						continue
					}

					// report result to aggregator
					log.Printf("worker%d sending body %s to aggregator", id, u)
					agg_body_input <- endpointData{Url: u, Body: string(body)}

					// send found urls to aggregator
					for _, u := range urls {
						log.Printf("worker%d sending new url %s to aggregator", id, u)
						pending.Add(1)
						agg_url_input <- u
					}

					// report we are done
					log.Printf("worker%d done with %s", id, u)
					pending.Add(-1)
				}
			}
			log.Printf("worker%d fell out of loop", id)
		}(i)
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
