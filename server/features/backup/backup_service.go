package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ivory/features/cluster"
	"ivory/features/deployment"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/user"
	"log/slog"
	"mime/multipart"
	"strings"
)

var ErrInvalidQueryType = errors.New("invalid query type")
var ErrInvalidQueryVariety = errors.New("invalid query variety")
var ErrInvalidStatus = errors.New("invalid status")
var ErrInvalidFeature = errors.New("invalid feature")
var ErrUnsupportedBackupVersion = errors.New("unsupported backup version")

const latestBackupVersion = "v2"

// Service manages the lifecycle of system backups (export/import).
// It handles backward compatibility by dispatching imports to version-specific
// handlers based on the filename naming convention.
type Service struct {
	clusterService    *cluster.Service
	queryService      *query.Service
	permissionService *permission.Service
	deploymentService *deployment.Service
	userService       *user.Service
}

func NewService(
	clusterService *cluster.Service,
	queryService *query.Service,
	permissionService *permission.Service,
	deploymentService *deployment.Service,
	userService *user.Service,
) *Service {
	return &Service{
		clusterService:    clusterService,
		queryService:      queryService,
		permissionService: permissionService,
		deploymentService: deploymentService,
		userService:       userService,
	}
}

// GetFileName returns the recommended filename for a backup export.
// The filename includes the version (e.g., ivory.v1.bak) which is used
// for version detection during import.
func (s *Service) GetFileName() string {
	return fmt.Sprintf("ivory.%s.bak", latestBackupVersion)
}

// Export is the entry point for creating a backup file in the latest supported
// wire format.
//
// Callers should treat the returned bytes as an opaque payload and should not
// depend on any specific versioned backup struct. When the backup schema must
// change in a backward-incompatible way, this method should switch to exporting
// the new version while preserving older import handlers for legacy file.
func (s *Service) Export() ([]byte, error) {
	switch latestBackupVersion {
	case "v2":
		backupModel, err := s.exportV2()
		if err != nil {
			return nil, err
		}
		return json.Marshal(backupModel)
	default:
		return nil, ErrUnsupportedBackupVersion
	}
}

// Import is the entry point for restoring data from a backup file.
// It detects the backup version from the filename and delegates the
// actual unmarshaling and data restoration to the appropriate version handler.
func (s *Service) Import(file *multipart.FileHeader) error {
	f, errOpen := file.Open()
	if errOpen != nil {
		return errOpen
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close backup file", "error", err)
		}
	}()

	data, errRead := io.ReadAll(f)
	if errRead != nil {
		return errRead
	}

	if strings.Contains(file.Filename, ".v2.") {
		return s.importV2(data)
	}
	// NOTE: a name with no version at all is a v1 file - v1 predates the
	// version marker, so the absence of one identifies it
	if strings.Contains(file.Filename, ".v1.") || !strings.Contains(file.Filename, ".v") {
		return s.importV1(data)
	}

	return fmt.Errorf("unsupported backup version in filename: %s", file.Filename)
}
