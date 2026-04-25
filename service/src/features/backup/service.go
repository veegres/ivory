package backup

import (
	"errors"
	"fmt"
	"io"
	"ivory/src/features/cluster"
	"ivory/src/features/permission"
	"ivory/src/features/query"
	"mime/multipart"
	"strings"
)

var ErrInvalidQueryType = errors.New("invalid query type")
var ErrInvalidQueryVariety = errors.New("invalid query variety")
var ErrInvalidStatus = errors.New("invalid status")
var ErrInvalidFeature = errors.New("invalid feature")

// Service manages the lifecycle of system backups (export/import).
// It handles backward compatibility by dispatching imports to version-specific
// handlers based on the filename naming convention.
type Service struct {
	clusterService    *cluster.Service
	queryService      *query.Service
	permissionService *permission.Service
	currentVersion    string
}

func NewService(
	clusterService *cluster.Service,
	queryService *query.Service,
	permissionService *permission.Service,
) *Service {
	return &Service{
		clusterService:    clusterService,
		queryService:      queryService,
		permissionService: permissionService,
		currentVersion:    "v1",
	}
}

// GetFileName returns the recommended filename for a backup export.
// The filename includes the version (e.g., ivory.v1.bak) which is used
// for version detection during import.
func (s *Service) GetFileName() string {
	return fmt.Sprintf("ivory.%s.bak", s.currentVersion)
}

// Import is the entry point for restoring data from a backup file.
// It detects the backup version from the filename and delegates the
// actual unmarshaling and data restoration to the appropriate version handler.
func (s *Service) Import(file *multipart.FileHeader) error {
	f, errOpen := file.Open()
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()

	data, errRead := io.ReadAll(f)
	if errRead != nil {
		return errRead
	}

	if strings.Contains(file.Filename, ".v1.") || !strings.Contains(file.Filename, ".v") {
		return s.ImportV1(data)
	}

	return fmt.Errorf("unsupported backup version in filename: %s", file.Filename)
}
