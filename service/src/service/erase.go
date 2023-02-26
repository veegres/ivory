package service

import (
	"errors"
	"ivory/src/persistence"
)

type EraseService struct {
	// TODO change everything to service
	passwordRepository     *persistence.PasswordRepository
	clusterRepository      *persistence.ClusterRepository
	certRepository         *persistence.CertRepository
	tagRepository          *persistence.TagRepository
	compactTableRepository *persistence.CompactTableRepository
	queryService           *QueryService
	secretService          *SecretService
	encryption             *Encryption
}

func NewEraseService(
	passwordRepository *persistence.PasswordRepository,
	clusterRepository *persistence.ClusterRepository,
	certRepository *persistence.CertRepository,
	tagRepository *persistence.TagRepository,
	compactTableRepository *persistence.CompactTableRepository,
	queryService *QueryService,
	secretService *SecretService,
) *EraseService {
	return &EraseService{
		passwordRepository:     passwordRepository,
		compactTableRepository: compactTableRepository,
		clusterRepository:      clusterRepository,
		certRepository:         certRepository,
		tagRepository:          tagRepository,
		queryService:           queryService,
		secretService:          secretService,
	}
}

func (s *EraseService) Erase() error {
	errSecret := s.secretService.Clean()
	errPass := s.passwordRepository.DeleteAll()
	errCert := s.certRepository.DeleteAll()
	errCluster := s.clusterRepository.DeleteAll()
	errComTable := s.compactTableRepository.DeleteAll()
	errTag := s.tagRepository.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery)
}
