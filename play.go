package scrapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Play(dumpFileName, bindAddr string) error {
	apiData, err := LoadFile(dumpFileName)
	if err != nil {
		return err
	}
	handler := makeHandler(apiData.State)
	http.HandleFunc("/", handler)
	log.Printf("Starting HTTP server listening on %s", bindAddr)
	return http.ListenAndServe(bindAddr, nil)

}

func LoadFile(fname string) (*Result, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(f)
	ret := &Result{}
	err = jsonParser.Decode(ret)
	return ret, err
}

func makeHandler(endpoints map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, DEFAULT_PAGE)
			return
		}
		body, ok := endpoints[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, body)
		}
	}
}
