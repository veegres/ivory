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
	authorised, authError := s.authInfo(authHeader, appConfig.Auth)

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

func (s *GeneralService) GetAppConfig() (*AppConfig, error) {
	if s.appConfig != nil {
		return s.appConfig, nil
	}

	read, err := s.configFiles.Read("data/config/" + s.appConfigFileName + ".json")
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig
	errUnmarshal := json.Unmarshal(read, &appConfig)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	s.appConfig = &appConfig
	return s.appConfig, nil
}

func (s *GeneralService) SetConfig(config AppConfig) error {
	if config.Company == "" {
		return errors.New("company name cannot be empty")
	}

	if s.configFiles.Exist(s.appConfigFileName) {
		return errors.New("config is already set up")
	}

	s.configFiles.Create(s.appConfigFileName)
	file, errOpen := s.configFiles.Open("data/config/" + s.appConfigFileName + ".json")
	if errOpen != nil {
		return errOpen
	}

	jsonAuth, errMarshall := json.Marshal(config)
	if errMarshall != nil {
		return errMarshall
	}

	_, errWrite := file.Write(jsonAuth)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

// TODO this should be moved to AuthService
func (s *GeneralService) authInfo(header string, authConfig AuthConfig) (bool, string) {
	switch authConfig.Type {
	case NONE:
		if header != "" {
			return false, "authentication header should be empty"
		}
		return true, ""
	case BASIC:
		token, errHeader := s.authService.GetTokenFromHeader(header)
		if errHeader != nil {
			return false, errHeader.Error()
		} else {
			errValidate := s.authService.ValidateToken(token, authConfig.Body["username"])
			if errValidate != nil {
				return false, errValidate.Error()
			}
		}
		return true, ""
	default:
		return false, "authentication type is not specified or incorrect"
	}
}
