package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	. "ivory/src/model"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type SidecarRequest[R any] struct {
	sidecar *SidecarClient
}

func NewSidecarRequest[R any](sidecar *SidecarClient) *SidecarRequest[R] {
	return &SidecarRequest[R]{
		sidecar: sidecar,
	}
}

func (p *SidecarRequest[R]) Get(instance InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodGet, instance, path, 1+time.Second)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Post(instance InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPost, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Put(instance InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPut, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Patch(instance InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodPatch, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) Delete(instance InstanceRequest, path string) (*R, int, error) {
	response, errRes := p.sidecar.send(http.MethodDelete, instance, path, 0)
	body, status, errParse := p.parseResponse(response)
	return body, status, errors.Join(errRes, errParse)
}

func (p *SidecarRequest[R]) parseResponse(res *http.Response) (*R, int, error) {
	var body R
	status := http.StatusBadGateway
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
				return nil, status, errMar
			}
		}
	}
	return &body, status, nil
}

type SidecarClient struct {
	passwordService *PasswordService
	certService     *CertService
}

func NewSidecarClient(
	passwordService *PasswordService,
	certService *CertService,
) *SidecarClient {
	return &SidecarClient{
		passwordService: passwordService,
		certService:     certService,
	}
}

func (p *SidecarClient) send(method string, instance InstanceRequest, path string, timeout time.Duration) (*http.Response, error) {
	if instance.Host == "" {
		return nil, errors.New("host cannot be empty")
	}
	client, protocol, errClient := p.getClient(instance.Certs, timeout)
	if errClient != nil {
		return nil, errClient
	}
	domain := instance.Host + ":" + strconv.Itoa(instance.Port)
	req, errRequest := p.getRequest(instance.CredentialId, method, protocol, domain, path, instance.Body)
	if errRequest != nil {
		return nil, errRequest
	}
	res, errDo := client.Do(req)
	return res, errDo
}

func (p *SidecarClient) getRequest(
	credentilId *uuid.UUID,
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

	if credentilId != nil {
		passInfo, errPass := p.passwordService.GetDecrypted(*credentilId)
		err = errPass
		req.SetBasicAuth(passInfo.Username, passInfo.Password)
	}

	return req, err
}

func (p *SidecarClient) getClient(certs Certs, timeout time.Duration) (*http.Client, string, error) {
	// TODO error overloading doesn't work it should be changed to a lot of returns with errors
	var err error
	var rootCa *x509.CertPool
	protocol := "http"

	// Setting Client CA
	if certs.ClientCAId != nil {
		clientCA, errCert := p.certService.GetFile(*certs.ClientCAId)
		rootCa = x509.NewCertPool()
		rootCa.AppendCertsFromPEM(clientCA)
		protocol = "https"
		err = errCert
	}

	// Setting Client Cert with Client Private Key
	var certificates []tls.Certificate
	if certs.ClientCertId != nil && certs.ClientKeyId != nil {
		clientCertInfo, errCert := p.certService.Get(*certs.ClientCertId)
		clientKeyInfo, errCert := p.certService.Get(*certs.ClientKeyId)

		cert, errCert := tls.LoadX509KeyPair(clientCertInfo.Path, clientKeyInfo.Path)
		certificates = append(certificates, cert)
		err = errCert
	}

	tlsConfig := &tls.Config{RootCAs: rootCa, Certificates: certificates}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Timeout: timeout}

	if certs.ClientCertId != nil || certs.ClientKeyId != nil {
		return client, protocol, errors.New("to be able to establish mutual tls connection you need to provide both client cert and client private key")
	}

	return client, protocol, err
}
