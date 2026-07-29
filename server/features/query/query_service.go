package query

import (
	"errors"
	"fmt"
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/plugins/database"
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
	databaseRegistry *utils.Registry[database.Plugin, database.Adapter]
	vaultService     *vault.Service
	certService      *cert.Service

	appName  string
	chartMap map[database.Plugin]map[ChartType]Request
}

func NewService(
	repository *Repository,
	databaseRegistry *utils.Registry[database.Plugin, database.Adapter],
	vaultService *vault.Service,
	certService *cert.Service,
	appName string,
) *Service {
	queryService := &Service{
		repository:       repository,
		databaseRegistry: databaseRegistry,
		vaultService:     vaultService,
		certService:      certService,
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

func (s *Service) SupportedFeatures(t database.Plugin) map[config.Feature]bool {
	c, e := s.databaseRegistry.Get(t)
	if e != nil {
		return map[config.Feature]bool{}
	}
	return c.SupportedFeatures()
}

func (s *Service) getDatabaseAdapter(queryCtx Context) (database.Adapter, database.Context, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, database.Context{}, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Plugin)
	return client, ctx, err
}

func (s *Service) initializeSystemCharts() {
	s.chartMap = make(map[database.Plugin]map[ChartType]Request)
	for t, adapter := range s.databaseRegistry.All() {
		s.chartMap[t] = make(map[ChartType]Request)
		for name, query := range adapter.SystemCharts() {
			s.chartMap[t][name] = Request{Name: string(name), Query: query}
		}
	}
}

// initializeSystemQueries seeds each plugin's system queries only when that
// plugin has none yet, so new plugins get their templates on upgrade without
// duplicating existing ones on restart.
func (s *Service) initializeSystemQueries() error {
	for plugin, adapter := range s.databaseRegistry.All() {
		exists, errExists := s.repository.HasSystemQueriesForPlugin(plugin)
		if errExists != nil {
			return errExists
		}
		if exists {
			continue
		}
		for _, req := range adapter.SystemRequests() {
			_, _, err := s.Create(System, mapSystemRequest(plugin, req))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
