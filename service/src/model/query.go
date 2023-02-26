package model

type QueryType int8

const (
	BLOAT QueryType = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

type QueryCreation string

const (
	Manual QueryCreation = "manual"
	System               = "system"
)

type Query struct {
	Name        string        `json:"name"`
	Type        QueryType     `json:"type"`
	Creation    QueryCreation `json:"creation"`
	Description string        `json:"description"`
	Default     string        `json:"default"`
	Edited      *string       `json:"edited"`
}
