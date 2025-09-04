package sidecar

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"ivory/src/features/cert"
	"ivory/src/features/instance"
	"ivory/src/features/password"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SidecarRequest[R any] struct {
	sidecar *SidecarGateway
}

func NewSidecarRequest[R any](sidecar *SidecarGateway) *SidecarRequest[R] {
	return &SidecarRequest[R]{
		sidecar: sidecar,
	}
}

func (p *SidecarRequest[R]) Get(instance instance.InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodGet, instance, path, 1+time.Second)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Post(instance instance.InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPost, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Put(instance instance.InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPut, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Patch(instance instance.InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPatch, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Delete(instance instance.InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodDelete, instance, path, 0)
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

type SidecarGateway struct {
	passwordService *password.PasswordService
	certService     *cert.CertService
}

func NewSidecarGateway(
	passwordService *password.PasswordService,
	certService *cert.CertService,
) *SidecarGateway {
	return &SidecarGateway{
		passwordService: passwordService,
		certService:     certService,
	}
}

func (p *SidecarGateway) send(method string, instance instance.InstanceRequest, path string, timeout time.Duration) (*http.Response, error) {
	if instance.Sidecar.Host == "" {
		return nil, errors.New("host cannot be empty")
	}
	client, protocol, errClient := p.getClient(instance.Certs, timeout)
	if errClient != nil {
		return nil, errClient
	}
	domain := instance.Sidecar.Host + ":" + strconv.Itoa(instance.Sidecar.Port)
	req, errRequest := p.getRequest(instance.CredentialId, method, protocol, domain, path, instance.Body)
	if errRequest != nil {
		return nil, errRequest
	}
	res, errDo := client.Do(req)
	return res, errDo
}

func (p *SidecarGateway) getRequest(
	credentialId *uuid.UUID,
	method string,
	protocol string,
	domain string,
	path string,
	body any,
) (*http.Request, error) {
	var err error
	var bodyBytes []byte

	if body != nil {
		bodyBytes, err = json.Marshal(body)
	}

	req, err := http.NewRequest(method, protocol+"://"+domain+path, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	if credentialId != nil {
		passInfo, errPass := p.passwordService.GetDecrypted(*credentialId)
		err = errPass
		req.SetBasicAuth(passInfo.Username, passInfo.Password)
	}

	return req, err
}

func (p *SidecarGateway) getClient(certs *cert.Certs, timeout time.Duration) (*http.Client, string, error) {
	protocol := "http"
	tlsConfig := &tls.Config{}

	if certs != nil {
		protocol = "https"
	}

	errTls := p.certService.EnrichTLSConfig(&tlsConfig, certs)
	if errTls != nil {
		return nil, protocol, errTls
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Timeout: timeout}
	return client, protocol, nil
}
