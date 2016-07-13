package scrapi

import (
	"encoding/json"
	"fmt"
	"html/template"
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

func makeHandler(apiData *Result) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			t, err := template.New("page").Parse(DEFAULT_PAGE)
			if err != nil {
				// this should never happen unless we have messed up the template
				// panic so it is never missed in test
				panic(err)
			}
			tdata := struct {
				R *Result
				D string
			}{
				R: apiData,
				D: BUILD_TIME,
			}
			err = t.Execute(w, tdata)
			//fmt.Fprintf(w, DEFAULT_PAGE)
			return
		}
		body, ok := apiData.State[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, body)
		}
	}
}
