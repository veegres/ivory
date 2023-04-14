package service

import (
	"errors"
)

type EraseService struct {
	passwordService *PasswordService
	clusterService  *ClusterService
	certService     *CertService
	tagService      *TagService
	bloatService    *BloatService
	queryService    *QueryService
	secretService   *SecretService
	encryption      *EncryptionService
}

func NewEraseService(
	passwordService *PasswordService,
	clusterService *ClusterService,
	certService *CertService,
	tagService *TagService,
	bloatService *BloatService,
	queryService *QueryService,
	secretService *SecretService,
) *EraseService {
	return &EraseService{
		passwordService: passwordService,
		bloatService:    bloatService,
		clusterService:  clusterService,
		certService:     certService,
		tagService:      tagService,
		queryService:    queryService,
		secretService:   secretService,
	}
}

func (s *EraseService) Erase() error {
	errSecret := s.secretService.Clean()
	errPass := s.passwordService.DeleteAll()
	errCert := s.certService.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errComTable := s.bloatService.DeleteAll()
	errTag := s.tagService.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery)
}
