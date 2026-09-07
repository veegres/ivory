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
	databaseRegistry *utils.Registry[database.PluginType, database.Adapter]
	vaultService     *vault.Service
	certService      *cert.Service

	appName  string
	chartMap map[database.PluginType]map[ChartType]Request
}

func NewService(
	repository *Repository,
	databaseRegistry *utils.Registry[database.PluginType, database.Adapter],
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

func (s *Service) SupportedFeatures(t database.PluginType) map[config.Feature]bool {
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
	s.chartMap = make(map[database.PluginType]map[ChartType]Request)
	for t, adapter := range s.databaseRegistry.All() {
		s.chartMap[t] = make(map[ChartType]Request)
		for name, query := range adapter.SystemCharts() {
			s.chartMap[t][name] = Request{Name: string(name), Query: query}
		}
	}
}

// initializeSystemQueries seeds the system queries a plugin ships but does not
// have stored yet, matched by name. Seeding used to skip a plugin that had any
// system query at all, which meant a query added to a plugin in a later release
// only ever reached installations created after it. Comparing per name still
// never duplicates or overwrites - a stored query keeps whatever edits were
// made to it, since only the ones absent by name are created.
func (s *Service) initializeSystemQueries() error {
	for plugin, adapter := range s.databaseRegistry.All() {
		existing, errExisting := s.repository.SystemQueryNamesForPlugin(plugin)
		if errExisting != nil {
			return errExisting
		}
		for _, req := range adapter.SystemRequests() {
			if existing[req.Name] {
				continue
			}
			_, _, err := s.Create(System, mapSystemRequest(plugin, req))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
