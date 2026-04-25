package tools

import (
	"errors"
	"ivory/src/clients/database"
	"ivory/src/features"
	"ivory/src/features/tools/bloat"
	"ivory/src/features/vault"
	"ivory/src/storage/db"
	"ivory/src/storage/files"
)

type Service struct {
	bloat *bloat.Service
}

func NewService(
	vaultService *vault.Service,
) *Service {
	// DB
	st := db.NewStorage("ivory_tools.db")
	compactTableBucket := db.NewBucket[bloat.Bloat](st, "CompactTable")

	// FILES
	compactTableFiles := files.NewStorage("pgcompacttable", ".log")

	// REPO
	bloatRepo := bloat.NewRepository(compactTableBucket, compactTableFiles)

	return &Service{
		bloat: bloat.NewService(bloatRepo, vaultService),
	}
}

func (s *Service) SupportedFeatures(t database.Type) []features.Feature {
	switch t {
	case database.POSTGRES:
		return []features.Feature{
			features.ViewBloatList,
			features.ViewBloatItem,
			features.ViewBloatLogs,
			features.ManageBloatJob,
		}
	default:
		return []features.Feature{}
	}
}

func (s *Service) DeleteAll() error {
	errBloat := s.bloat.DeleteAll()
	return errors.Join(errBloat)
}
