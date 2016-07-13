package scrapi

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ReqFunc func(url.URL) ([]byte, error)

func MakeClient(host, username, password string, insecure bool) ReqFunc {
	tr := &http.Transport{}
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
	return func(target url.URL) ([]byte, error) {
		if target.Host == "" {
			target.Host = host
		}
		if target.Scheme == "" {
			target.Scheme = "https"
		}
		req, err := http.NewRequest("GET", target.String(), nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(username, password)

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		log.Print(res.Status, target.String())

		return ioutil.ReadAll(res.Body)
	}
}
