package scrapi

import (
	"encoding/json"
	"github.com/aktungmak/odata-client"
	"net/url"
)

func Process(target *url.URL, rf ReqFunc) ([]byte, []*url.URL, error) {
	//make http request
	body, err := rf(*target)
	if err != nil {
		return nil, nil, err
	}

	// sometimes the body is empty, so just return here
	if len(body) == 0 {
		return nil, nil, nil
	}

	// parse response as json
	var jdata map[string]interface{}
	err = json.Unmarshal(body, &jdata)
	if err != nil {
		return nil, nil, err
	}

	// parse response for urls
	lm := odata.ParseLinks(jdata, "")
	urls := make([]*url.URL, 0, len(lm))
	for _, u := range lm {
		urls = append(urls, u)
	}
	return body, urls, nil
}
