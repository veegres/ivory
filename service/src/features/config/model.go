package config

type AuthType int8

const (
	NONE AuthType = iota
	BASIC
)

type AuthConfig struct {
	Type AuthType          `json:"type"`
	Body map[string]string `json:"body"`
}

type AppConfig struct {
	Company      string       `json:"company"`
	Availability Availability `json:"availability"`
	Auth         AuthConfig   `json:"auth"`
}

type Availability struct {
	ManualQuery bool `json:"manualQuery"`
}
