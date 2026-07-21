package patroni

import "encoding/json"

// COMMON (WEB AND SERVER)

// SPECIFIC (SERVER)

type cluster struct {
	Members             []instance           `json:"members"`
	ScheduledSwitchover *scheduledSwitchover `json:"scheduled_switchover"`
	Pause               bool                 `json:"pause"`
}

type instance struct {
	Name           string `json:"name"`
	State          string `json:"state"`
	Role           string `json:"role"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	ApiUrl         string `json:"api_url"`
	PendingRestart bool   `json:"pending_restart"`
	// it can be int / nil / "unknown"
	Lag              json.RawMessage   `json:"lag"`
	Timeline         int               `json:"timeline"`
	ScheduledRestart *scheduledRestart `json:"scheduled_restart"`
	Tags             *map[string]any   `json:"tags"`
}

type scheduledRestart struct {
	RestartPending      bool   `json:"restart_pending"`
	Schedule            string `json:"schedule"`
	PostmasterStartTime string `json:"postmaster_start_time"`
}

type scheduledSwitchover struct {
	At   string `json:"at"`
	From string `json:"from"`
	To   string `json:"to"`
}

type configPause struct {
	Pause bool `json:"pause"`
}
