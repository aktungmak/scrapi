package scrapi

type Result struct {
	sync.RWMutex
	m map[string]string
}

func Rec(rootUri string) map[string]string {
	num_workers := 16
	queue := make(chan string, num_workers*4)
	done := make(chan struct{})
	idle := &sync.WaitGroup{}
	result := &Result{m: make(map[string]string)}

	queue <- rootUri
	for i := 0; i < num_workers; i++ {
		go Worker(queue, done, idle, result)
	}
	idle.Wait()
	for i := 0; i < num_workers; i++ {
		done <- struct{}{}
	}
	println(result)
}

func Worker(queue chan string, done chan struct{}, wg waitGroup, result *Result) {
	for {
		wg.Add(1)
		select {
		case <-done:
			return
		case url <- queue:
			wg.Add(-1)
			body, urls := Process(url)
			result.Lock()
			resulti.m[url] = body
			result.Unlock()
			for u := range urls {
				result.RLock()
				_, ok := result.m[u]
				result.RUnlock()
				if !ok {
					queue <- u
				}
			}
		}
	}
}
