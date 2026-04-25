package query

import (
	"ivory/src/clients/database"

	"github.com/google/uuid"
)

func (s *Service) GetLog(queryUuid uuid.UUID) ([]database.QueryFields, error) {
	return s.repository.getLog(queryUuid)
}

func (s *Service) AddLog(uuid uuid.UUID, element any) error {
	return s.repository.addLog(uuid, element)
}

func (s *Service) DeleteLog(queryUuid uuid.UUID) error {
	return s.repository.deleteLog(queryUuid)
}

func (s *Service) HasLog(queryUuid uuid.UUID) bool {
	return s.repository.hasLog(queryUuid)
}

func (s *Service) DeleteAllLogs() error {
	return s.repository.deleteAllLogs()
}
