package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"ivory/model"
	"ivory/persistence"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func Get[R any](instance model.InstanceRequest, path string) (R, int, error) {
	return send[R](http.MethodGet, instance, path, 1+time.Second)
}

func Post[R any](instance model.InstanceRequest, path string) (R, int, error) {
	return send[R](http.MethodPost, instance, path, 0)
}

func Put[R any](instance model.InstanceRequest, path string) (R, int, error) {
	return send[R](http.MethodPut, instance, path, 0)
}

func Patch[R any](instance model.InstanceRequest, path string) (R, int, error) {
	return send[R](http.MethodPatch, instance, path, 0)
}

func Delete[R any](instance model.InstanceRequest, path string) (R, int, error) {
	return send[R](http.MethodDelete, instance, path, 0)
}

func send[R any](method string, instance model.InstanceRequest, path string, timeout time.Duration) (R, int, error) {
	clusterInfo, err := persistence.BoltDB.Cluster.Get(instance.Cluster)
	client, protocol, err := getClient(clusterInfo.CertId, timeout)
	domain := instance.Host + ":" + strconv.Itoa(instance.Port)
	req, err := getRequest(clusterInfo.PatroniCredId, method, protocol, domain, path, instance.Body)
	res, err := client.Do(req)

	var body R
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

func getRequest(
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
		passInfo, errPass := persistence.BoltDB.Credential.Get(*passwordId)
		err = errPass
		req.SetBasicAuth(passInfo.Username, passInfo.Password)
	}

	return req, err
}

func getClient(certId *uuid.UUID, timeout time.Duration) (*http.Client, string, error) {
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
	client := &http.Client{Transport: transport, Timeout: timeout}

	return client, protocol, err
}
