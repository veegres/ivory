package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	. "ivory/src/model"
	"ivory/src/persistence"
	"path"
	"strings"
)

type CertService struct {
	certRepository *persistence.CertRepository
}

func NewCertService(certRepository *persistence.CertRepository) *CertService {
	return &CertService{certRepository: certRepository}
}

func (s *CertService) Get(uuid uuid.UUID) (Cert, error) {
	return s.certRepository.Get(uuid)
}

func (s *CertService) GetFile(uuid uuid.UUID) ([]byte, error) {
	return s.certRepository.GetFile(uuid)
}

func (s *CertService) List() (CertMap, error) {
	return s.certRepository.List()
}

func (s *CertService) ListByType(certType CertType) (CertMap, error) {
	return s.certRepository.ListByType(certType)
}

func (s *CertService) Create(pathStr string, certType CertType, fileUsageType FileUsageType) (*Cert, error) {
	formats := []string{".crt", ".key", ".chain"}

	fileName := path.Base(pathStr)
	if fileName == "" {
		return nil, errors.New("file name cannot be empty")
	}

	fileFormat := path.Ext(pathStr)
	idx := slices.IndexFunc(formats, func(s string) bool { return s == fileFormat })
	if idx == -1 {
		return nil, errors.New("file format is not correct, allowed formats: " + strings.Join(formats, ", "))
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
