package tools

import (
	"errors"
	"ivory/features"
	"ivory/features/tools/pg_compacttable"
	"ivory/plugins/database"
)

type Tool interface {
	SupportedFeatures(t database.Plugin) []features.Feature
	DeleteAll() error
}

// NOTE: validate that is matches interface in compile-time
var _ Tool = (*pg_compacttable.Service)(nil)

type Service struct {
	tools          []Tool
	pgCompactTable *pg_compacttable.Service
}

func NewService(
	pgCompactTable *pg_compacttable.Service,
) *Service {
	return &Service{
		tools: []Tool{
			pgCompactTable,
		},
		pgCompactTable: pgCompactTable,
	}
}

func (s *Service) SupportedFeatures(t database.Plugin) []features.Feature {
	allFeatures := make([]features.Feature, 0)
	for _, tool := range s.tools {
		allFeatures = append(allFeatures, tool.SupportedFeatures(t)...)
	}
	return allFeatures
}

func (s *Service) DeleteAll() error {
	errs := make([]error, 0)
	for _, tool := range s.tools {
		errs = append(errs, tool.DeleteAll())
	}
	return errors.Join(errs...)
}
