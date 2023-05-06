package model

// COMMON (WEB AND SERVER)

// SPECIFIC (SERVER)

type PatroniCluster struct {
	Members []PatroniInstance `json:"members"`
}

type PatroniInstance struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Role     string `json:"role"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	ApiUrl   string `json:"api_url"`
	Lag      *int   `json:"lag"`
	Timeline int    `json:"timeline"`
}

type PatroniInfo struct {
	DatabaseSystemIdentifier string               `json:"database_system_identifier"`
	PostmasterStartTime      string               `json:"postmaster_start_time"`
	Timeline                 int                  `json:"timeline"`
	ClusterUnlocked          bool                 `json:"cluster_unlocked"`
	State                    string               `json:"state"`
	Role                     string               `json:"role"`
	ServerVersion            int                  `json:"server_version"`
	Replication              []PatroniReplication `json:"replication"`
	Xlog                     Xlog                 `json:"xlog"`
	Patroni                  Patroni              `json:"patroni"`
}

type Patroni struct {
	Scope   string `json:"scope"`
	Version string `json:"version"`
}

type Xlog struct {
	Location          int    `json:"location"`
	ReceivedLocation  int    `json:"received_location"`
	ReplayedTimestamp string `json:"replayed_timestamp"`
	Paused            bool   `json:"paused"`
	ReplayedLocation  int    `json:"replayed_location"`
}

type PatroniReplication struct {
	SyncState       string `json:"sync_state"`
	SyncPriority    int    `json:"sync_priority"`
	ClientAddr      string `json:"client_addr"`
	State           string `json:"state"`
	ApplicationName string `json:"application_name"`
	Username        string `json:"username"`
}
