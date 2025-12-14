package cert

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

var ErrMutualTLSRequiresBothCertAndKey = errors.New("to be able to establish mutual tls connection, you need to provide both client cert and client private key")
var ErrFileNameCannotBeEmpty = errors.New("file name cannot be empty")

type Service struct {
	certRepository *Repository
	formats        []string
}

func NewService(certRepository *Repository) *Service {
	return &Service{
		certRepository: certRepository,

		formats: []string{".crt", ".key", ".chain"},
	}
}

func (s *Service) Get(uuid uuid.UUID) (Cert, error) {
	return s.certRepository.Get(uuid)
}

// GetTLSConfigRootCA Setting Client CA
func (s *Service) GetTLSConfigRootCA(clientCAId *uuid.UUID) (*x509.CertPool, error) {
	var rootCa *x509.CertPool

	if clientCAId != nil {
		clientCAInfo, errCert := s.Get(*clientCAId)
		if errCert != nil {
			return nil, errCert
		}
		clientCA, errCa := os.ReadFile(clientCAInfo.Path)
		if errCa != nil {
			return nil, errCa
		}
		rootCa = x509.NewCertPool()
		rootCa.AppendCertsFromPEM(clientCA)
	}

	return rootCa, nil
}

// GetTLSConfigCertificates Setting Client Cert with Client Private Key
func (s *Service) GetTLSConfigCertificates(clientCertId *uuid.UUID, clientKeyId *uuid.UUID) ([]tls.Certificate, error) {
	var certificates []tls.Certificate
	if clientCertId != nil && clientKeyId != nil {
		clientCertInfo, errCert := s.Get(*clientCertId)
		if errCert != nil {
			return nil, errCert
		}
		clientKeyInfo, errKey := s.Get(*clientKeyId)
		if errKey != nil {
			return nil, errKey
		}

		cert, errX509 := tls.LoadX509KeyPair(clientCertInfo.Path, clientKeyInfo.Path)
		if errX509 != nil {
			return nil, errX509
		}
		certificates = append(certificates, cert)
	}

	// NOTE: xor operation for `nil` check
	if (clientCertId == nil || clientKeyId == nil) && clientCertId != clientKeyId {
		return nil, ErrMutualTLSRequiresBothCertAndKey
	}

	return certificates, nil
}

func (s *Service) EnrichTLSConfig(config **tls.Config, certs *Certs) error {
	if *config == nil {
		*config = &tls.Config{}
	}
	if certs == nil {
		*config = nil
		return nil
	}
	rootCA, errRoot := s.GetTLSConfigRootCA(certs.ClientCAId)
	if errRoot != nil {
		return errRoot
	}
	(*config).RootCAs = rootCA

	certsConfig, errCert := s.GetTLSConfigCertificates(certs.ClientCertId, certs.ClientKeyId)
	if errCert != nil {
		return errCert
	}
	(*config).Certificates = certsConfig
	return nil
}

func (s *Service) List() (CertMap, error) {
	return s.certRepository.List()
}

func (s *Service) ListByType(certType CertType) (CertMap, error) {
	return s.certRepository.ListByType(certType)
}

func (s *Service) Create(pathStr string, certType CertType, fileUsageType FileUsageType) (*Cert, error) {
	fileName := path.Base(pathStr)
	if fileName == "" {
		return nil, ErrFileNameCannotBeEmpty
	}

	fileFormat := path.Ext(pathStr)
	idx := slices.IndexFunc(s.formats, func(s string) bool { return s == fileFormat })
	if idx == -1 {
		return nil, errors.New("file format is not correct, allowed formats: " + strings.Join(s.formats, ", "))
	}

	cert := Cert{
		FileUsageType: fileUsageType,
		Type:          certType,
		FileName:      fileName,
	}

	return s.certRepository.Create(cert, pathStr)
}

func (s *Service) Delete(certUuid uuid.UUID) error {
	return s.certRepository.Delete(certUuid)
}

func (s *Service) DeleteAll() error {
	return s.certRepository.DeleteAll()
}
