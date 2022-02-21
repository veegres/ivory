package model

type ClusterModel struct {
	Name  string   `json:"name"`
	Nodes []string `json:"nodes"`
}
