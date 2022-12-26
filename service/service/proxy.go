package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/google/uuid"
	"ivory/model"
	"ivory/persistence"
	"net/http"
	"strconv"
)

var proxy = &prx{}

type prx struct{}

func (p *prx) Get(instance model.InstanceRequest, path string) (interface{}, error) {
	return p.send(http.MethodGet, instance, path)
}

func (p *prx) Post(instance model.InstanceRequest, path string) (interface{}, error) {
	return p.send(http.MethodPost, instance, path)
}

func (p *prx) Put(instance model.InstanceRequest, path string) (interface{}, error) {
	return p.send(http.MethodPut, instance, path)
}

func (p *prx) Patch(instance model.InstanceRequest, path string) (interface{}, error) {
	return p.send(http.MethodPatch, instance, path)
}

func (p *prx) Delete(instance model.InstanceRequest, path string) (interface{}, error) {
	return p.send(http.MethodDelete, instance, path)
}

func (p *prx) send(method string, instance model.InstanceRequest, path string) (interface{}, error) {
	clusterInfo, err := persistence.Database.Cluster.Get(instance.Cluster)
	client, protocol, err := p.getClient(clusterInfo.CertId)
	domain := instance.Host + ":" + strconv.Itoa(int(instance.Port))
	req, err := p.getRequest(clusterInfo.PatroniCredId, method, protocol, domain, path, instance.Body)
	res, err := client.Do(req)
	var body interface{}
	return json.NewDecoder(res.Body).Decode(&body), err
}

func (p *prx) getRequest(
	passwordId *uuid.UUID,
	method string,
	protocol string,
	domain string,
	path string,
	body string,
) (*http.Request, error) {
	var err error

	req, err := http.NewRequest(method, protocol+"://"+domain+path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	if passwordId != nil {
		passInfo, errPass := persistence.Database.Credential.Get(*passwordId)
		err = errPass
		req.SetBasicAuth(passInfo.Username, passInfo.Password)
	}

	return req, err
}

func (p *prx) getClient(certId *uuid.UUID) (*http.Client, string, error) {
	var err error
	var caCertPool *x509.CertPool
	protocol := "http"

	if certId != nil {
		certInfo, errCert := persistence.Database.Cert.Get(*certId)
		caCert, errCert := persistence.File.Certs.Read(certInfo.Path)
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		protocol = "https"
		err = errCert
	}

	tlsConfig := &tls.Config{RootCAs: caCertPool}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	return client, protocol, err
}
