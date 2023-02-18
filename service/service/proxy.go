package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"ivory/config"
	"ivory/model"
	"ivory/persistence"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type ProxyRequest[R any] struct {
	proxy *Proxy
}

func NewProxyRequest[R any](proxy *Proxy) *ProxyRequest[R] {
	return &ProxyRequest[R]{
		proxy: proxy,
	}
}

func (p *ProxyRequest[R]) Get(instance model.InstanceRequest, path string) (R, int, error) {
	response, status, err := p.proxy.send(http.MethodGet, instance, path, 1+time.Second)
	return response.(R), status, err
}

func (p *ProxyRequest[R]) Post(instance model.InstanceRequest, path string) (R, int, error) {
	response, status, err := p.proxy.send(http.MethodPost, instance, path, 0)
	return response.(R), status, err
}

func (p *ProxyRequest[R]) Put(instance model.InstanceRequest, path string) (R, int, error) {
	response, status, err := p.proxy.send(http.MethodPut, instance, path, 0)
	return response.(R), status, err
}

func (p *ProxyRequest[R]) Patch(instance model.InstanceRequest, path string) (R, int, error) {
	response, status, err := p.proxy.send(http.MethodPatch, instance, path, 0)
	return response.(R), status, err
}

func (p *ProxyRequest[R]) Delete(instance model.InstanceRequest, path string) (R, int, error) {
	response, status, err := p.proxy.send(http.MethodDelete, instance, path, 0)
	return response.(R), status, err
}

type Proxy struct {
	clusterRepository  *persistence.ClusterRepository
	passwordRepository *persistence.PasswordRepository
	certRepository     *persistence.CertRepository
	certFile           *config.FilePersistence
}

func NewProxy(
	clusterRepository *persistence.ClusterRepository,
	passwordRepository *persistence.PasswordRepository,
	certRepository *persistence.CertRepository,
	certFile *config.FilePersistence,
) *Proxy {
	return &Proxy{
		clusterRepository:  clusterRepository,
		passwordRepository: passwordRepository,
		certRepository:     certRepository,
		certFile:           certFile,
	}
}

func (p *Proxy) send(method string, instance model.InstanceRequest, path string, timeout time.Duration) (interface{}, int, error) {
	clusterInfo, err := p.clusterRepository.Get(instance.Cluster)
	client, protocol, err := p.getClient(clusterInfo.Certs, timeout)
	domain := instance.Host + ":" + strconv.Itoa(instance.Port)
	req, err := p.getRequest(clusterInfo.Credentials.PatroniId, method, protocol, domain, path, instance.Body)
	res, err := client.Do(req)

	var body interface{}
	var status int
	if res != nil && err == nil {
		contentType := res.Header.Get("Content-Type")
		status = res.StatusCode

		switch contentType {
		case "application/json":
			err = json.NewDecoder(res.Body).Decode(&body)
		case "text/html":
			text, errRead := io.ReadAll(res.Body)
			reflect.ValueOf(&body).Elem().SetString(string(text))
			err = errRead
		default:
			err = errors.New("doesn't support such Content-Type in response")
		}
	}

	return body, status, err
}

func (p *Proxy) getRequest(
	passwordId *uuid.UUID,
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

	if passwordId != nil {
		passInfo, errPass := p.passwordRepository.Get(*passwordId)
		err = errPass
		req.SetBasicAuth(passInfo.Username, passInfo.Password)
	}

	return req, err
}

func (p *Proxy) getClient(certs model.Certs, timeout time.Duration) (*http.Client, string, error) {
	// TODO error overloading doesn't work it should be changed to a lot of returns with errors
	var err error
	var rootCa *x509.CertPool
	protocol := "http"

	// Setting Client CA
	if certs.ClientCAId != nil {
		clientCAInfo, errCert := p.certRepository.Get(*certs.ClientCAId)
		clientCA, errCert := p.certFile.Read(clientCAInfo.Path)
		rootCa = x509.NewCertPool()
		rootCa.AppendCertsFromPEM(clientCA)
		protocol = "https"
		err = errCert
	}

	// Setting Client Cert with Client Private Key
	var certificates []tls.Certificate
	if certs.ClientCertId != nil && certs.ClientKeyId != nil {
		clientCertInfo, errCert := p.certRepository.Get(*certs.ClientCertId)
		clientKeyInfo, errCert := p.certRepository.Get(*certs.ClientKeyId)

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
