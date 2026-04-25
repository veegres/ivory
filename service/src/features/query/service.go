package query

import (
	"errors"
	"ivory/src/clients/database"
	"ivory/src/clients/database/postgres"
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

	chartMap map[database.QueryChartType]database.Query
}

func NewService(
	repository *Repository,
	databaseRegistry *database.Registry,
	vaultService *vault.Service,
	certService *cert.Service,
	secretService *secret.Service,
) *Service {
	queryService := &Service{
		repository:       repository,
		databaseRegistry: databaseRegistry,
		vaultService:     vaultService,
		certService:      certService,
		secretService:    secretService,
		chartMap:         postgres.CreateChartsMap(),
	}
	err := queryService.createDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}
