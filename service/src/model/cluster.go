package model

import "github.com/google/uuid"

type ClusterModel struct {
	Name        string      `json:"name"`
	Certs       Certs       `json:"certs"`
	Credentials Credentials `json:"credentials"`
	Instances   []Sidecar   `json:"instances"`
	Tags        []string    `json:"tags"`
}

type ClusterAutoModel struct {
	Name        string      `json:"name"`
	Certs       Certs       `json:"certs"`
	Credentials Credentials `json:"credentials"`
	Instance    Sidecar     `json:"instance"`
	Tags        []string    `json:"tags"`
}

type Certs struct {
	ClientCAId   *uuid.UUID `json:"clientCAId"`
	ClientKeyId  *uuid.UUID `json:"clientKeyId"`
	ClientCertId *uuid.UUID `json:"clientCertId"`
}

type Credentials struct {
	PatroniId  *uuid.UUID `json:"patroniId"`
	PostgresId *uuid.UUID `json:"postgresId"`
}
