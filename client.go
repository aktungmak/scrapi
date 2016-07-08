package scrapi

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ReqFunc func(url.URL) ([]byte, error)

func MakeClient(host, username, password string, insecure bool) ReqFunc {
	return func(target url.URL) ([]byte, error) {
		tr := &http.Transport{}
		if insecure {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client := &http.Client{Transport: tr}
		if target.Host == "" {
			target.Host = host
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

		return ioutil.ReadAll(res.Body)
	}
}
