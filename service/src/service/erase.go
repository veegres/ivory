package service

import (
	"errors"
	"ivory/src/persistence"
)

type EraseService struct {
	// TODO change everything to service
	passwordService        *PasswordService
	clusterService         *ClusterService
	certRepository         *persistence.CertRepository
	tagRepository          *persistence.TagRepository
	compactTableRepository *persistence.BloatRepository
	queryService           *QueryService
	secretService          *SecretService
	encryption             *EncryptionService
}

func NewEraseService(
	passwordService *PasswordService,
	clusterService *ClusterService,
	certRepository *persistence.CertRepository,
	tagRepository *persistence.TagRepository,
	compactTableRepository *persistence.BloatRepository,
	queryService *QueryService,
	secretService *SecretService,
) *EraseService {
	return &EraseService{
		passwordService:        passwordService,
		compactTableRepository: compactTableRepository,
		clusterService:         clusterService,
		certRepository:         certRepository,
		tagRepository:          tagRepository,
		queryService:           queryService,
		secretService:          secretService,
	}
}

func (s *EraseService) Erase() error {
	errSecret := s.secretService.Clean()
	errPass := s.passwordService.DeleteAll()
	errCert := s.certRepository.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errComTable := s.compactTableRepository.DeleteAll()
	errTag := s.tagRepository.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery)
}
