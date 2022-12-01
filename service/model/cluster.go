package model

import "github.com/google/uuid"

type ClusterModel struct {
	Name           string     `json:"name"`
	CertId         *uuid.UUID `json:"certId"`
	PatroniCredId  *uuid.UUID `json:"patroniCredId"`
	PostgresCredId *uuid.UUID `json:"postgresCredId"`
	Nodes          []string   `json:"nodes"`
	Tags           []string   `json:"tags"`
}
