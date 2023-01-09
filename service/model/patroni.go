package model

type PatroniCluster struct {
	Members []PatroniInstance `json:"members"`
}

type PatroniInstance struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Role     string `json:"role"`
	Host     string `json:"host"`
	Port     int8   `json:"port"`
	ApiUrl   string `json:"api_url"`
	Lag      int32  `json:"lag"`
	Timeline int32  `json:"timeline"`
}

type PatroniInfo struct {
	DatabaseSystemIdentifier string               `json:"database_system_identifier"`
	PostmasterStartTime      string               `json:"postmaster_start_time"`
	Timeline                 int32                `json:"timeline"`
	ClusterUnlocked          bool                 `json:"cluster_unlocked"`
	State                    string               `json:"state"`
	Role                     string               `json:"role"`
	ServerVersion            int32                `json:"server_version"`
	Replication              []PatroniReplication `json:"replication"`
	Xlog                     Xlog                 `json:"xlog"`
	Patroni                  Patroni              `json:"patroni"`
}

type Patroni struct {
	Scope   string `json:"scope"`
	Version string `json:"version"`
}

type Xlog struct {
	Location          int32 `json:"location"`
	ReceivedLocation  int32 `json:"received_location"`
	ReplayedTimestamp int32 `json:"replayed_timestamp"`
	Paused            int32 `json:"paused"`
	ReplayedLocation  int32 `json:"replayed_location"`
}

type PatroniReplication struct {
	SyncState       string `json:"sync_state"`
	SyncPriority    int8   `json:"sync_priority"`
	ClientAddr      string `json:"client_addr"`
	State           string `json:"state"`
	ApplicationName string `json:"application_name"`
	Username        string `json:"username"`
}
