package model

type Cluster struct {
	Name  string   `json:"name"`
	Nodes []string `json:"nodes"`
}
