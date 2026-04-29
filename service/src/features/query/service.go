package query

import (
	"errors"
	"fmt"
	"ivory/src/clients/database"
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/secret"
	"ivory/src/features/vault"
)

var ErrQueryEmpty = errors.New("query is empty")
var ErrVaultProblems = errors.New("vault problems, check if it exists")
var ErrAllFieldsRequired = errors.New("all fields have to be filled")
var ErrNameChangeNotAllowed = errors.New("name change is not allowed for system queries")
var ErrTypeChangeNotAllowed = errors.New("type change is not allowed for system queries")
var ErrDescriptionChangeNotAllowed = errors.New("description change is not allowed for system queries")
var ErrDeletionOfSystemQueriesRestricted = errors.New("deletion of system queries is restricted")

type Service struct {
	repository       *Repository
	databaseRegistry *database.Registry
	vaultService     *vault.Service
	certService      *cert.Service
	secretService    *secret.Service

	appName  string
	chartMap map[database.Type]map[ChartType]Request
}

func NewService(
	repository *Repository,
	databaseRegistry *database.Registry,
	vaultService *vault.Service,
	certService *cert.Service,
	secretService *secret.Service,
	appName string,
) *Service {
	queryService := &Service{
		repository:       repository,
		databaseRegistry: databaseRegistry,
		vaultService:     vaultService,
		certService:      certService,
		secretService:    secretService,
		appName:          appName,
	}
	queryService.initializeSystemCharts()
	err := queryService.initializeSystemQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *Service) GetApplicationName(session string) string {
	return s.appName + " [" + fmt.Sprintf("%.7s", session) + "]"
}

func (s *Service) SupportedFeatures(t database.Type) []features.Feature {
	c, e := s.databaseRegistry.Get(t)
	if e != nil {
		return []features.Feature{}
	}
	return c.SupportedFeatures()
}

func (s *Service) initializeSystemCharts() {
	s.chartMap = make(map[database.Type]map[ChartType]Request)
	for t, adapter := range s.databaseRegistry.All() {
		s.chartMap[t] = make(map[ChartType]Request)
		for name, query := range adapter.SystemCharts() {
			s.chartMap[t][name] = Request{Name: string(name), Query: query}
		}
	}
}

func (s *Service) initializeSystemQueries() error {
	if !s.secretService.IsRefEmpty() {
		return nil
	}

	for _, adapter := range s.databaseRegistry.All() {
		for _, req := range adapter.SystemRequests() {
			_, _, err := s.Create(System, s.mapSystemRequest(req))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) mapSystemRequest(req database.SystemRequest) Request {
	t := req.Type
	v := make([]VarietyType, len(req.Varieties))
	for i, val := range req.Varieties {
		v[i] = val
	}
	return Request{
		Name:        req.Name,
		Type:        &t,
		Description: &req.Description,
		Query:       req.Query,
		Varieties:   v,
		Params:      req.Params,
	}
}
