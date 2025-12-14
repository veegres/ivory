package sidecar

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var ErrHostEmpty = errors.New("host cannot be empty")

type SidecarRequest[R any] struct {
	sidecar *Gateway
}

func NewSidecarRequest[R any](sidecar *Gateway) *SidecarRequest[R] {
	return &SidecarRequest[R]{
		sidecar: sidecar,
	}
}

func (p *SidecarRequest[R]) Get(request Request, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodGet, request, path, 1+time.Second)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Post(request Request, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPost, request, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Put(request Request, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPut, request, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Patch(request Request, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPatch, request, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Delete(request Request, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodDelete, request, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) parseResponse(res *http.Response) (*R, int, error) {
	var body R
	status := http.StatusBadRequest
	if res != nil {
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

type Gateway struct{}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (p *Gateway) send(method string, request Request, path string, timeout time.Duration) (*http.Response, error) {
	if request.Sidecar.Host == "" {
		return nil, ErrHostEmpty
	}
	client, protocol, errClient := p.getClient(request.TlsConfig, timeout)
	if errClient != nil {
		return nil, errClient
	}
	domain := request.Sidecar.Host + ":" + strconv.Itoa(request.Sidecar.Port)
	req, errRequest := p.getRequest(request.Credentials, method, protocol, domain, path, request.Body)
	if errRequest != nil {
		return nil, errRequest
	}
	res, errDo := client.Do(req)
	return res, errDo
}

func (p *Gateway) getRequest(
	credentials *Credentials,
	method string,
	protocol string,
	domain string,
	path string,
	body any,
) (*http.Request, error) {
	var bodyBytes []byte

	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	url := protocol + "://" + domain + path
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	if credentials != nil {
		req.SetBasicAuth(credentials.Username, credentials.Password)
	}

	return req, nil
}

func (p *Gateway) getClient(tlsConfig *tls.Config, timeout time.Duration) (*http.Client, string, error) {
	protocol := "http"

	if tlsConfig != nil {
		protocol = "https"
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Timeout: timeout}
	return client, protocol, nil
}
