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
	client, protocol, err := getClient(clusterInfo.Certs, timeout)
	domain := instance.Host + ":" + strconv.Itoa(instance.Port)
	req, err := getRequest(clusterInfo.Credentials.PatroniId, method, protocol, domain, path, instance.Body)
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

func getClient(certs model.Certs, timeout time.Duration) (*http.Client, string, error) {
	// TODO error overturning doesn't work it should be changed to a lot of returns with errors
	var err error
	var rootCa *x509.CertPool
	protocol := "http"

	// Setting Client CA
	if certs.ClientCAId != nil {
		clientCAInfo, errCert := persistence.BoltDB.Cert.Get(*certs.ClientCAId)
		clientCA, errCert := persistence.File.Certs.Read(clientCAInfo.Path)
		rootCa = x509.NewCertPool()
		rootCa.AppendCertsFromPEM(clientCA)
		protocol = "https"
		err = errCert
	}

	// Setting Client Cert with Client Private Key
	var certificates []tls.Certificate
	if certs.ClientCertId != nil && certs.ClientKeyId != nil {
		clientCertInfo, errCert := persistence.BoltDB.Cert.Get(*certs.ClientCertId)
		clientKeyInfo, errCert := persistence.BoltDB.Cert.Get(*certs.ClientKeyId)

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
