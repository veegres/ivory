package model

import "github.com/google/uuid"

type ClusterModel struct {
	Name           string     `json:"name"`
	CertsId        *uuid.UUID `json:"certsId"`
	PatroniCredId  *uuid.UUID `json:"patroniCredId"`
	PostgresCredId *uuid.UUID `json:"postgresCredId"`
	Nodes          []string   `json:"nodes"`
	Tags           []string   `json:"tags"`
}
