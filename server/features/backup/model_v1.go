package backup

// BackupV1 represents the legacy backup format used in Ivory v1.
//
// SACRED RULE: This structure and all its subtypes MUST NOT be changed.
// If the internal system models change, the mapping logic in service_import_v1.go
// must be updated to translate these fixed structures into the new internal models.
// If a new backup format is required, create a new BackupV2 instead.
type BackupV1 struct {
	Clusters    []backupClusterV1     `json:"clusters"`
	Queries     []backupQueryV1       `json:"queries"`
	Permissions []backupPermissionsV1 `json:"permissions"`
}

type backupClusterV1 struct {
	Name     string            `json:"name"`
	Tags     []string          `json:"tags"`
	Sidecars []backupSidecarV1 `json:"sidecars"`
}

type backupSidecarV1 struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type backupQueryV1 struct {
	Name        string                 `json:"name"`
	Type        backupQueryTypeV1      `json:"type"`
	Varieties   []backupQueryVarietyV1 `json:"varieties"`
	Params      []string               `json:"params"`
	Description *string                `json:"description"`
	Default     string                 `json:"default"`
	Custom      string                 `json:"custom"`
}

type backupQueryTypeV1 int8

const (
	BLOAT_V1 backupQueryTypeV1 = iota
	ACTIVITY_V1
	REPLICATION_V1
	STATISTIC_V1
	OTHER_V1
)

type backupQueryVarietyV1 int8

const (
	DatabaseSensitiveV1 backupQueryVarietyV1 = iota
	MasterOnlyV1
	ReplicaRecommendedV1
)

type backupPermissionsV1 struct {
	Username    string                            `json:"username"`
	Permissions map[string]backupPermissionTypeV1 `json:"permissions"`
}

type backupPermissionTypeV1 int

const (
	NOT_PERMITTED_V1 backupPermissionTypeV1 = iota
	PENDING_V1
	GRANTED_V1
)
