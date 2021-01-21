package connect

import (
	"bytes"
	"io"
	"net/http"
)

type Requests struct {
	Name    string
	Request *http.Request
}

func (b *Spider) NewRequests(method, u string, body []byte) (*Requests, error) {
	var reqBody io.Reader
	if body == nil {
		reqBody = nil
	} else {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(b.CtxRoot, method, u, reqBody)
	if err != nil {
		return nil, err
	}

	requests := &Requests{}
	requests.Name = b.Name
	requests.Request = req
	req.Header = b.DefaultHeader

	return requests, nil
}

func (b *Spider) GET(u string) (*Requests, error) {
	return b.NewRequests("GET", u, nil)
}

func (b *Spider) POST(u string, body []byte) (*Requests, error) {
	return b.NewRequests("POST", u, body)
}

func (b *Spider) HEAD(c *Spider, u string) (*Requests, error) {
	return b.NewRequests("HEAD", u, nil)
}

func (b *Spider) OPTIONS(c *Spider, u string) (*Requests, error) {
	return b.NewRequests("OPTIONS", u, nil)
}
