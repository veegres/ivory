package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	. "ivory/src/model"
	"ivory/src/persistence"
	"os"
	"path"
	"strings"
)

type CertService struct {
	certRepository *persistence.CertRepository
	formats        []string
}

func NewCertService(certRepository *persistence.CertRepository) *CertService {
	return &CertService{
		certRepository: certRepository,

		formats: []string{".crt", ".key", ".chain"},
	}
}

func (s *CertService) Get(uuid uuid.UUID) (Cert, error) {
	return s.certRepository.Get(uuid)
}

// GetTLSConfigRootCA Setting Client CA
func (s *CertService) GetTLSConfigRootCA(certs Certs) (*x509.CertPool, error) {
	var rootCa *x509.CertPool

	if certs.ClientCAId != nil {
		clientCAInfo, errCert := s.Get(*certs.ClientCAId)
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
func (s *CertService) GetTLSConfigCertificates(certs Certs) ([]tls.Certificate, error) {
	var certificates []tls.Certificate
	if certs.ClientCertId != nil && certs.ClientKeyId != nil {
		clientCertInfo, errCert := s.Get(*certs.ClientCertId)
		if errCert != nil {
			return nil, errCert
		}
		clientKeyInfo, errKey := s.Get(*certs.ClientKeyId)
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
	if (certs.ClientCertId == nil || certs.ClientKeyId == nil) && certs.ClientCertId != certs.ClientKeyId {
		return nil, errors.New("to be able to establish mutual tls connection you need to provide both client cert and client private key")
	}

	return certificates, nil
}

func (s *CertService) EnrichTLSConfig(config *tls.Config, certs *Certs) error {
	if certs == nil {
		config = nil
		return nil
	}
	rootCA, errRoot := s.GetTLSConfigRootCA(*certs)
	if errRoot != nil {
		return errRoot
	}
	config.RootCAs = rootCA

	certsConfig, errCert := s.GetTLSConfigCertificates(*certs)
	if errCert != nil {
		return errCert
	}
	config.Certificates = certsConfig
	return nil
}

func (s *CertService) List() (CertMap, error) {
	return s.certRepository.List()
}

func (s *CertService) ListByType(certType CertType) (CertMap, error) {
	return s.certRepository.ListByType(certType)
}

func (s *CertService) Create(pathStr string, certType CertType, fileUsageType FileUsageType) (*Cert, error) {
	fileName := path.Base(pathStr)
	if fileName == "" {
		return nil, errors.New("file name cannot be empty")
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

func (s *CertService) Delete(certUuid uuid.UUID) error {
	return s.certRepository.Delete(certUuid)
}

func (s *CertService) DeleteAll() error {
	return s.certRepository.DeleteAll()
}
