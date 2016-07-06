package scrapi

import (
    "net/url"
)

func Process(target url.URL, rf ReqFunc) ([]string, string, error) {
	//make http request
    body, err := rf(target)
    if err != nil {
        return nil, "", err
    }
	//parse response for urls
	

	return urls, string(body), nil

}
