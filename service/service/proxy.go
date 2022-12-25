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

func (p *prx) Get(cluster string, instance model.Instance, path string) (interface{}, error) {
    return p.send(http.MethodGet, cluster, instance, path)
}

func (p *prx) Post(cluster string, instance model.Instance, path string) (interface{}, error) {
    return p.send(http.MethodPost, cluster, instance, path)
}

func (p *prx) Put(cluster string, instance model.Instance, path string) (interface{}, error) {
    return p.send(http.MethodPut, cluster, instance, path)
}

func (p *prx) Patch(cluster string, instance model.Instance, path string) (interface{}, error) {
    return p.send(http.MethodPatch, cluster, instance, path)
}

func (p *prx) Delete(cluster string, instance model.Instance, path string) (interface{}, error) {
    return p.send(http.MethodDelete, cluster, instance, path)
}

func (p *prx) send(method string, cluster string, instance model.Instance, path string) (interface{}, error) {
	clusterInfo, err := persistence.Database.Cluster.Get(cluster)
	client, protocol, err := p.getClient(clusterInfo.CertId)
    domain := instance.Host+":"+strconv.Itoa(int(instance.Port))
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
