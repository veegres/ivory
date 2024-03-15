package model

import "encoding/json"

// COMMON (WEB AND SERVER)

// SPECIFIC (SERVER)

type PatroniCluster struct {
	Members             []PatroniInstance           `json:"members"`
	ScheduledSwitchover *PatroniScheduledSwitchover `json:"scheduled_switchover"`
	Pause               bool                        `json:"pause"`
}

type PatroniInstance struct {
	Name           string `json:"name"`
	State          string `json:"state"`
	Role           string `json:"role"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	ApiUrl         string `json:"api_url"`
	PendingRestart bool   `json:"pending_restart"`
	// it can be int / nil / "unknown"
	Lag              json.RawMessage          `json:"lag"`
	Timeline         int                      `json:"timeline"`
	ScheduledRestart *PatroniScheduledRestart `json:"scheduled_restart"`
}

type PatroniScheduledRestart struct {
	RestartPending      bool   `json:"restart_pending"`
	Schedule            string `json:"schedule"`
	PostmasterStartTime string `json:"postmaster_start_time"`
}

type PatroniScheduledSwitchover struct {
	At   string `json:"at"`
	From string `json:"from"`
	To   string `json:"to"`
}
