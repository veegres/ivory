package tools

import (
	"errors"
	"ivory/clients/shell"
	"ivory/features"
	"ivory/features/job"
	"ivory/features/tools/pg_compacttable"
	"ivory/features/vault"
	"ivory/plugins/database"
	"ivory/storage/db"
)

type Service struct {
	bloat *bloat.Service
}

func NewService(
	vaultService *vault.Service,
	shellClient *shell.Client,
	jobManager *job.Manager,
) *Service {
	// DB
	st := db.NewStorage("ivory_tools.db")
	compactTableBucket := db.NewBucket[bloat.Response](st, "CompactTable")

	// REPO
	bloatRepo := bloat.NewRepository(compactTableBucket)

	return &Service{
		bloat: bloat.NewService(bloatRepo, shellClient, vaultService, jobManager),
	}
}

func (s *Service) SupportedFeatures(t database.Plugin) []features.Feature {
	switch t {
	case database.POSTGRES:
		return []features.Feature{
			features.ViewToolBloatList,
			features.ViewToolBloatItem,
			features.ViewToolBloatLogs,
			features.ManageToolBloatJob,
		}
	default:
		return []features.Feature{}
	}
}

func (s *Service) DeleteAll() error {
	errBloat := s.bloat.DeleteAll()
	return errors.Join(errBloat)
}
