package connect

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	Name string
	Resp *http.Response
	Body []byte
}

func NewResponse(name string, resp *http.Response) (*Response, error) {
	var err error
	response := &Response{}
	response.Name = name
	response.Resp = resp
	response.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, err
}
