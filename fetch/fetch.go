package fetch

import (
	"bytes"
	"net/http"
	"time"
)

type Api struct {
	BaseUrl string
	Headers map[string]string
	Client  *http.Client
	Timeout time.Duration //TODO: Apply Timeout
}

type ApiGetRequestOptions struct {
	Headers map[string]string
	Query   map[string]string
}

type ApiPostRequestOptions struct {
	Headers map[string]string
}

// TODO: Add missing request types - put, patch, delete
// TODO: Add Version of Request With Json Decoding BuiltIn - GetJson etc..
func (client Api) Get(url string, options ApiGetRequestOptions) (*http.Response, error) {
	req, err := http.NewRequest("GET", client.BaseUrl+url, nil)

	if err != nil {
		return &http.Response{}, err
	}

	for key, header := range client.Headers {
		req.Header.Add(key, header)
	}

	for key, header := range options.Headers {
		req.Header.Add(key, header)
	}

	queryString := req.URL.Query()

	for key, value := range options.Query {
		queryString.Add(key, value)
	}

	req.URL.RawQuery = queryString.Encode()

	return client.Client.Do(req)
}

func (client Api) Post(url string, rawBody []byte, options ApiPostRequestOptions) (*http.Response, error) {
	body := bytes.NewBuffer(rawBody)

	req, err := http.NewRequest("POST", client.BaseUrl+url, body)

	if err != nil {
		return &http.Response{}, err
	}

	for key, header := range client.Headers {
		req.Header.Add(key, header)
	}

	for key, header := range options.Headers {
		req.Header.Add(key, header)
	}

	return client.Client.Do(req)
}
