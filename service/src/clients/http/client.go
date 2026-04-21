package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	nethttp "net/http"
	"reflect"
	"strconv"
	"time"
)

var ErrHostEmpty = errors.New("host cannot be empty")

type Credentials struct {
	Username string
	Password string
}

type Request struct {
	Method      string
	Host        string
	Port        int
	Path        string
	Body        any
	Credentials *Credentials
	TLSConfig   *tls.Config
	Timeout     time.Duration
}

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Send(req Request) (*nethttp.Response, error) {
	if req.Host == "" {
		return nil, ErrHostEmpty
	}

	protocol := "http"
	if req.TLSConfig != nil {
		protocol = "https"
	}

	transport := &nethttp.Transport{TLSClientConfig: req.TLSConfig}
	httpClient := &nethttp.Client{Transport: transport, Timeout: req.Timeout}

	domain := req.Host + ":" + strconv.Itoa(req.Port)
	url := protocol + "://" + domain + req.Path

	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq, err := nethttp.NewRequest(req.Method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if req.Credentials != nil {
		httpReq.SetBasicAuth(req.Credentials.Username, req.Credentials.Password)
	}

	return httpClient.Do(httpReq)
}

type JSONRequest[R any] struct {
	client *Client
}

func NewJSONRequest[R any](client *Client) *JSONRequest[R] {
	return &JSONRequest[R]{
		client: client,
	}
}

func (p *JSONRequest[R]) Get(request Request) (*R, int, error) {
	request.Method = nethttp.MethodGet
	if request.Timeout == 0 {
		request.Timeout = 1 * time.Second
	}
	return p.Do(request)
}

func (p *JSONRequest[R]) Post(request Request) (*R, int, error) {
	request.Method = nethttp.MethodPost
	return p.Do(request)
}

func (p *JSONRequest[R]) Put(request Request) (*R, int, error) {
	request.Method = nethttp.MethodPut
	return p.Do(request)
}

func (p *JSONRequest[R]) Patch(request Request) (*R, int, error) {
	request.Method = nethttp.MethodPatch
	return p.Do(request)
}

func (p *JSONRequest[R]) Delete(request Request) (*R, int, error) {
	request.Method = nethttp.MethodDelete
	return p.Do(request)
}

func (p *JSONRequest[R]) Do(req Request) (*R, int, error) {
	response, errRes := p.client.Send(req)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *JSONRequest[R]) parseResponse(res *nethttp.Response) (*R, int, error) {
	var body R
	status := nethttp.StatusBadRequest
	if res != nil {
		defer res.Body.Close()
		status = res.StatusCode
		bytesBody, errRead := io.ReadAll(res.Body)
		if errRead != nil {
			return nil, status, errRead
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			errCode := errors.New(res.Status)
			errParse := errors.New(string(bytesBody))
			return nil, status, errors.Join(errCode, errParse)
		}

		errMar := json.Unmarshal(bytesBody, &body)
		if errMar != nil {
			if reflect.TypeOf(body) == reflect.TypeOf("") {
				reflect.ValueOf(&body).Elem().SetString(string(bytesBody))
			} else {
				return nil, status, errors.Join(errMar, errors.New(string(bytesBody)))
			}
		}
	}
	return &body, status, nil
}
