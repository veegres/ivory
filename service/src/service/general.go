package service

import (
	"errors"
)

type GeneralService struct {
	passwordService *PasswordService
	clusterService  *ClusterService
	certService     *CertService
	tagService      *TagService
	bloatService    *BloatService
	queryService    *QueryService
	secretService   *SecretService
	encryption      *EncryptionService
}

func NewGeneralService(
	passwordService *PasswordService,
	clusterService *ClusterService,
	certService *CertService,
	tagService *TagService,
	bloatService *BloatService,
	queryService *QueryService,
	secretService *SecretService,
) *GeneralService {
	return &GeneralService{
		passwordService: passwordService,
		bloatService:    bloatService,
		clusterService:  clusterService,
		certService:     certService,
		tagService:      tagService,
		queryService:    queryService,
		secretService:   secretService,
	}
}

func (s *GeneralService) Erase() error {
	errSecret := s.secretService.Clean()
	errPass := s.passwordService.DeleteAll()
	errCert := s.certService.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errComTable := s.bloatService.DeleteAll()
	errTag := s.tagService.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery)
}

func (s *GeneralService) ChangeSecret(previousKey string, newKey string) error {
	prevSha, newSha, err := s.secretService.Update(previousKey, newKey)
	if err != nil {
		return err
	}
	errEnc := s.passwordService.ReEncryptPasswords(prevSha, newSha)
	if errEnc != nil {
		return errEnc
	}
	return nil
}
