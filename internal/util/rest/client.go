package rest_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SyamSolution/transaction-service/config"
	"go.uber.org/zap"
)

var (
	Client HttpClient
)

type RestTransport struct {
	Client  HttpClient
	Url     string
	Method  string
	Header  http.Header
	Payload interface{}
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewRestTransport(httpClient *http.Client) *RestTransport {
	return &RestTransport{
		Client: httpClient,
	}
}

func (r *RestTransport) Generate() (req *http.Request, err error) {
	body := new(bytes.Buffer)
	if r.Payload != nil {

		err := json.NewEncoder(body).Encode(r.Payload)
		if err != nil {
			return nil, fmt.Errorf("error encoding payload, %s", err)
		}
	}

	req, err = http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating http request, %s", err)
	}
	if r.Header != nil {
		req.Header = r.Header
	}

	return req, nil
}

func Request(req RestTransport, base *config.BaseDep) (res *http.Response, err error) {
	var httpReq *http.Request
	var body *bytes.Buffer
	if req.Payload != nil {
		b, err := json.Marshal(req.Payload)
		if err != nil {
			return res, err
		}
		body = bytes.NewBuffer(b)
	}

	if body == nil {
		httpReq, err = http.NewRequest(req.Method, req.Url, nil)
	} else {
		httpReq, err = http.NewRequest(req.Method, req.Url, body)
	}

	if err != nil {
		return res, err
	}

	if req.Header != nil {
		httpReq.Header = req.Header
	}

	res, err = req.Client.Do(httpReq)
	if err != nil {
		return res, err
	}

	return res, nil
}

func CloseResponseBody(r io.ReadCloser, base *config.BaseDep) {
	if err := r.Close(); err != nil {
		base.Logger.Error("failed closing", zap.Error(err))
	}
}
