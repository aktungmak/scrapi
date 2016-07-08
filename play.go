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
	handler := makeHandler(apiData)
	http.HandleFunc("/", handler)
	log.Printf("Starting HTTP server listening on %s", bindAddr)
	http.ListenAndServe(bindAddr, nil)
	return nil

}

func LoadFile(fname string) (map[string]string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(f)
	var ret map[string]string
	err = jsonParser.Decode(&ret)
	return ret, err
}

func makeHandler(endpoints map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "it works")
	}
}