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
				return nil, status, errors.Join(errMar, errors.New(string(bytesBody)))
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

func (p *SidecarClient) getClient(certs Certs, timeout time.Duration) (*http.Client, string, error) {
	var rootCa *x509.CertPool
	protocol := "http"

	// Setting Client CA
	if certs.ClientCAId != nil {
		protocol = "https"
		clientCA, errCa := p.certService.GetFile(*certs.ClientCAId)
		if errCa != nil {
			return nil, protocol, errCa
		}
		rootCa = x509.NewCertPool()
		rootCa.AppendCertsFromPEM(clientCA)
	}

	// Setting Client Cert with Client Private Key
	var certificates []tls.Certificate
	if certs.ClientCertId != nil && certs.ClientKeyId != nil {
		protocol = "https"
		clientCertInfo, errCert := p.certService.Get(*certs.ClientCertId)
		if errCert != nil {
			return nil, protocol, errCert
		}
		clientKeyInfo, errKey := p.certService.Get(*certs.ClientKeyId)
		if errKey != nil {
			return nil, protocol, errKey
		}

		cert, errX509 := tls.LoadX509KeyPair(clientCertInfo.Path, clientKeyInfo.Path)
		if errX509 != nil {
			return nil, protocol, errX509
		}
		certificates = append(certificates, cert)
	}

	// xor operation for `nil` check
	if (certs.ClientCertId == nil || certs.ClientKeyId == nil) && certs.ClientCertId != certs.ClientKeyId {
		return nil, protocol, errors.New("to be able to establish mutual tls connection you need to provide both client cert and client private key")
	}

	tlsConfig := &tls.Config{RootCAs: rootCa, Certificates: certificates}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Timeout: timeout}

	return client, protocol, nil
}
