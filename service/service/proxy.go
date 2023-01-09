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

func Get[R any](instance model.InstanceRequest, path string) (R, error) {
	return send[R](http.MethodGet, instance, path)
}

func Post[R any](instance model.InstanceRequest, path string) (R, error) {
	return send[R](http.MethodPost, instance, path)
}

func Put[R any](instance model.InstanceRequest, path string) (R, error) {
	return send[R](http.MethodPut, instance, path)
}

func Patch[R any](instance model.InstanceRequest, path string) (R, error) {
	return send[R](http.MethodPatch, instance, path)
}

func Delete[R any](instance model.InstanceRequest, path string) (R, error) {
	return send[R](http.MethodDelete, instance, path)
}

func send[R any](method string, instance model.InstanceRequest, path string) (R, error) {
	clusterInfo, err := persistence.BoltDB.Cluster.Get(instance.Cluster)
	client, protocol, err := getClient(clusterInfo.CertId)
	domain := instance.Host + ":" + strconv.Itoa(int(instance.Port))
	req, err := getRequest(clusterInfo.PatroniCredId, method, protocol, domain, path, instance.Body)
	res, err := client.Do(req)
	var body R
	if err == nil {
		err = json.NewDecoder(res.Body).Decode(&body)
	}
	return body, err
}

func getRequest(
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
		passInfo, errPass := persistence.BoltDB.Credential.Get(*passwordId)
		err = errPass
		req.SetBasicAuth(passInfo.Username, passInfo.Password)
	}

	return req, err
}

func getClient(certId *uuid.UUID) (*http.Client, string, error) {
	var err error
	var caCertPool *x509.CertPool
	protocol := "http"

	if certId != nil {
		certInfo, errCert := persistence.BoltDB.Cert.Get(*certId)
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
