package service

import (
	"encoding/json"
	"errors"
	"ivory/src/config"
	. "ivory/src/model"
)

type GeneralService struct {
	env             *config.Env
	configFiles     *config.FileGateway
	authService     *AuthService
	passwordService *PasswordService
	clusterService  *ClusterService
	certService     *CertService
	tagService      *TagService
	bloatService    *BloatService
	queryService    *QueryService
	secretService   *SecretService
	encryption      *EncryptionService

	appConfigFileName string
	appConfig         *AppConfig
}

func NewGeneralService(
	env *config.Env,
	configFiles *config.FileGateway,
	authService *AuthService,
	passwordService *PasswordService,
	clusterService *ClusterService,
	certService *CertService,
	tagService *TagService,
	bloatService *BloatService,
	queryService *QueryService,
	secretService *SecretService,
) *GeneralService {
	generalService := &GeneralService{
		env:             env,
		configFiles:     configFiles,
		authService:     authService,
		passwordService: passwordService,
		bloatService:    bloatService,
		clusterService:  clusterService,
		certService:     certService,
		tagService:      tagService,
		queryService:    queryService,
		secretService:   secretService,

		appConfigFileName: "application",
		appConfig:         nil,
	}

	appConfig, _ := generalService.GetAppConfig()
	generalService.appConfig = appConfig

	return generalService
}

func (s *GeneralService) Erase() error {
	errSecret := s.secretService.Clean()
	errPass := s.passwordService.DeleteAll()
	errCert := s.certService.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errComTable := s.bloatService.DeleteAll()
	errTag := s.tagService.DeleteAll()
	errQuery := s.queryService.DeleteAll()

	errConfig := s.configFiles.DeleteAll()
	s.appConfig = nil

	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery, errConfig)
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

	appConfig, errAppConfig := s.GetAppConfig()
	if errAppConfig != nil {
		return errAppConfig
	}
	reEncryptAuthConfig, errAuthConfig := s.authService.ReEncryptAuthConfig(appConfig.Auth, prevSha, newSha)
	if errAuthConfig != nil {
		return errAuthConfig
	}
	updatedAppConfig := AppConfig{
		Company:      appConfig.Company,
		Availability: appConfig.Availability,
		Auth:         reEncryptAuthConfig,
	}
	errUpdateAppConfig := s.UpdateAppConfig(updatedAppConfig)
	if errUpdateAppConfig != nil {
		return errUpdateAppConfig
	}

	return nil
}

func (s *GeneralService) GetAppInfo(authHeader string) *AppInfo {
	appConfig, err := s.GetAppConfig()
	if err != nil {
		return &AppInfo{
			Company:      "",
			Configured:   false,
			Secret:       s.secretService.Status(),
			Version:      s.env.Version,
			Availability: Availability{ManualQuery: false},
			Auth: AuthInfo{
				Type:       NONE,
				Authorised: false,
				Error:      "",
			},
		}
	}
	authorised, authError := s.authService.ValidateHeader(authHeader, appConfig.Auth)

	return &AppInfo{
		Company:      appConfig.Company,
		Configured:   true,
		Secret:       s.secretService.Status(),
		Version:      s.env.Version,
		Availability: appConfig.Availability,
		Auth: AuthInfo{
			Type:       appConfig.Auth.Type,
			Authorised: authorised,
			Error:      authError,
		},
	}
}

func (s *GeneralService) UpdateAppConfig(config AppConfig) error {
	file, errOpen := s.configFiles.Open(s.appConfigFileName)
	if errOpen != nil {
		return errOpen
	}

	errTruncate := file.Truncate(0)
	if errTruncate != nil {
		return errTruncate
	}

	jsonAuth, errMarshall := json.MarshalIndent(config, "", "  ")
	if errMarshall != nil {
		return errMarshall
	}

	_, errWrite := file.Write(jsonAuth)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func (s *GeneralService) GetAppConfig() (*AppConfig, error) {
	if s.appConfig != nil {
		return s.appConfig, nil
	}

	read, err := s.configFiles.Read(s.appConfigFileName)
	if err != nil {
		return nil, errors.New("config is not specified or it cannot be read from file")
	}

	var appConfig AppConfig
	errUnmarshal := json.Unmarshal(read, &appConfig)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	s.appConfig = &appConfig
	return s.appConfig, nil
}

func (s *GeneralService) SetAppConfig(newAppConfig AppConfig) error {
	if newAppConfig.Company == "" {
		return errors.New("company name cannot be empty")
	}

	// NOTE: we skip error checking here, because we need to know only if config is not set up yet
	currentAppConfig, _ := s.GetAppConfig()
	if currentAppConfig != nil {
		return errors.New("config is already set up")
	}

	encryptedAuthConfig, errAuthConfig := s.authService.EncryptAuthConfig(newAppConfig.Auth)
	if errAuthConfig != nil {
		return errAuthConfig
	}

	updatedAppConfig := AppConfig{
		Company:      newAppConfig.Company,
		Availability: newAppConfig.Availability,
		Auth:         encryptedAuthConfig,
	}

	_, errCreate := s.configFiles.Create(s.appConfigFileName)
	if errCreate != nil {
		return errCreate
	}

	return s.UpdateAppConfig(updatedAppConfig)
}
